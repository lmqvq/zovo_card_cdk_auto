/**
 * Read a customer-provided Session file in the browser.
 *
 * The file never leaves the page: callers receive the extracted text and can
 * put it into the existing Session input before the normal preflight flow.
 */

import { extractSessionJsonFromText, rtfToText } from './rtf'

export const MAX_SESSION_FILE_BYTES = 1024 * 1024
export const SESSION_FILE_ACCEPT =
  '.txt,.docx,.rtf,text/plain,application/rtf,text/rtf,text/richtext,application/vnd.openxmlformats-officedocument.wordprocessingml.document'

export type SessionFileErrorCode =
  | 'too_large'
  | 'unsupported'
  | 'empty'
  | 'invalid_docx'
  | 'invalid_rtf'
  | 'read_failed'

export class SessionFileError extends Error {
  readonly code: SessionFileErrorCode

  constructor(code: SessionFileErrorCode) {
    super(code)
    this.name = 'SessionFileError'
    this.code = code
  }
}

function fileKind(file: File): 'txt' | 'docx' | 'rtf' | null {
  const name = String(file.name || '').toLowerCase()
  const type = String(file.type || '').toLowerCase()
  if (name.endsWith('.rtf') || type === 'application/rtf' || type === 'text/rtf' || type === 'text/richtext') {
    return 'rtf'
  }
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

/** RTF 是 7-bit 文本 + `\'hh` 字节转义，按单字节读取即可（避免 UTF-8 解码把 `\'e4` 之类拆坏）。 */
function decodeRtfBytes(bytes: Uint8Array): string {
  let start = 0
  if (bytes.length >= 3 && bytes[0] === 0xef && bytes[1] === 0xbb && bytes[2] === 0xbf) start = 3
  let out = ''
  const chunkSize = 0x8000
  for (let i = start; i < bytes.length; i += chunkSize) {
    out += String.fromCharCode(...bytes.subarray(i, i + chunkSize))
  }
  return out
}

/**
 * Read an RTF document (Word / TextEdit / WPS 另存的 .rtf) and return the
 * Session JSON inside it. RTF 里的 `\{` `\}`、`\u8203` 零宽空格、`\par`
 * 拆行都由 rtf.ts 还原，调用方拿到的就是可直接提交的 JSON。
 */
async function readRtf(file: File): Promise<string> {
  try {
    const raw = decodeRtfBytes(new Uint8Array(await file.arrayBuffer()))
    const text = rtfToText(raw)
    const json = extractSessionJsonFromText(text)
    if (!json) throw new SessionFileError('invalid_rtf')
    return json
  } catch (error) {
    if (error instanceof SessionFileError) throw error
    throw new SessionFileError('invalid_rtf')
  }
}

/** Read a .txt / .docx / .rtf Session file locally and return its trimmed text. */
export async function readSessionFile(file: File): Promise<string> {
  if (file.size > MAX_SESSION_FILE_BYTES) throw new SessionFileError('too_large')

  const kind = fileKind(file)
  if (!kind) throw new SessionFileError('unsupported')

  try {
    const text =
      kind === 'txt'
        ? await file.text()
        : kind === 'rtf'
          ? await readRtf(file)
          : await readDocx(file)
    const trimmed = text.trim()
    if (!trimmed) throw new SessionFileError('empty')
    return trimmed
  } catch (error) {
    if (error instanceof SessionFileError) throw error
    throw new SessionFileError('read_failed')
  }
}
