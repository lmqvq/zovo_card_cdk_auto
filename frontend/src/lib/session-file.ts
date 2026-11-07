/**
 * Read a customer-provided Session file in the browser.
 *
 * The file never leaves the page: callers receive the extracted text and can
 * put it into the existing Session input before the normal preflight flow.
 */

export const MAX_SESSION_FILE_BYTES = 1024 * 1024
export const SESSION_FILE_ACCEPT =
  '.txt,.docx,text/plain,application/vnd.openxmlformats-officedocument.wordprocessingml.document'

export type SessionFileErrorCode =
  | 'too_large'
  | 'unsupported'
  | 'empty'
  | 'invalid_docx'
  | 'read_failed'

export class SessionFileError extends Error {
  readonly code: SessionFileErrorCode

  constructor(code: SessionFileErrorCode) {
    super(code)
    this.name = 'SessionFileError'
    this.code = code
  }
}

function fileKind(file: File): 'txt' | 'docx' | null {
  const name = String(file.name || '').toLowerCase()
  const type = String(file.type || '').toLowerCase()
  if (name.endsWith('.txt') || type === 'text/plain') return 'txt'
  if (
    name.endsWith('.docx') ||
    type === 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
  ) {
    return 'docx'
  }
  return null
}

function readWordText(node: Node): string {
  if (node.nodeType !== 1) return ''
  const element = node as Element
  const localName = element.localName || element.nodeName.split(':').pop() || ''
  if (localName === 't') return element.textContent || ''
  if (localName === 'tab') return '\t'
  if (localName === 'br' || localName === 'cr') return '\n'

  let text = ''
  for (const child of Array.from(node.childNodes)) text += readWordText(child)
  if (localName === 'p') text += '\n'
  return text
}

function isZipFile(bytes: Uint8Array): boolean {
  return (
    bytes.length >= 4 &&
    bytes[0] === 0x50 &&
    bytes[1] === 0x4b &&
    (bytes[2] === 0x03 || bytes[2] === 0x05 || bytes[2] === 0x07) &&
    (bytes[3] === 0x04 || bytes[3] === 0x06 || bytes[3] === 0x08)
  )
}

function looksLikeSessionText(text: string): boolean {
  const trimmed = text.trim()
  if (!trimmed) return false
  if (trimmed.startsWith('{')) return true
  return trimmed.startsWith('eyJ') && trimmed.split('.').length >= 5
}

async function readDocx(file: File): Promise<string> {
  try {
    const bytes = new Uint8Array(await file.arrayBuffer())

    // Some providers export the Session JSON as a text file and only give it
    // a .docx name. Accept that form when it is clearly a Session payload.
    if (!isZipFile(bytes)) {
      const text = new TextDecoder().decode(bytes).trim()
      if (looksLikeSessionText(text)) return text
      throw new SessionFileError('invalid_docx')
    }

    const { unzipSync, strFromU8 } = await import('fflate')
    const files = unzipSync(bytes)
    const documentXml = files['word/document.xml']
    if (!documentXml) throw new SessionFileError('invalid_docx')

    const xml = new DOMParser().parseFromString(strFromU8(documentXml), 'application/xml')
    if (xml.getElementsByTagName('parsererror').length > 0) {
      throw new SessionFileError('invalid_docx')
    }
    return readWordText(xml.documentElement).trim()
  } catch (error) {
    if (error instanceof SessionFileError) throw error
    throw new SessionFileError('read_failed')
  }
}

/** Read a .txt or .docx Session file locally and return its trimmed text. */
export async function readSessionFile(file: File): Promise<string> {
  if (file.size > MAX_SESSION_FILE_BYTES) throw new SessionFileError('too_large')

  const kind = fileKind(file)
  if (!kind) throw new SessionFileError('unsupported')

  try {
    const text = kind === 'txt' ? await file.text() : await readDocx(file)
    const trimmed = text.trim()
    if (!trimmed) throw new SessionFileError('empty')
    return trimmed
  } catch (error) {
    if (error instanceof SessionFileError) throw error
    throw new SessionFileError('read_failed')
  }
}
