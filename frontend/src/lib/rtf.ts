/**
 * 从 RTF 文档里取出 ChatGPT Session JSON。
 *
 * 客户经常把 Session 粘进 Word / TextEdit 再存成 .rtf 发过来。RTF 不是纯文本：
 *   - `{` `}` 会被转义成 `\{` `\}`，`\` 会变成 `\\`；
 *   - 长 token 会被编辑器插入 `\u8203`（U+200B 零宽空格）当断行点；
 *   - 还可能被 `\par` 拆行、被智能引号（U+201C/201D）替换掉 `"`；
 *   - 字体表 / 颜色表 / 生成器信息等组不是正文，必须跳过。
 *
 * 这里只做"RTF → 可见文本 → JSON 候选"的提取，不依赖 DOM，浏览器和 Node 都能跑。
 */

const SKIP_DESTINATIONS = new Set([
  'fonttbl',
  'colortbl',
  'stylesheet',
  'info',
  'pict',
  'object',
  'header',
  'headerl',
  'headerr',
  'footer',
  'footerl',
  'footerr',
  'footnote',
  'field',
  'fldinst',
  'bkmkstart',
  'bkmkend',
  'listtable',
  'listoverridetable',
  'listpicture',
  'rsidtbl',
  'generator',
  'xmlnstbl',
  'themedata',
  'colorschememapping',
  'datastore',
  'latentstyles',
  'filetbl',
  'revtbl',
  'pgdsctbl',
  'mmathPr',
])

const ZERO_WIDTH_RE = /[\u200b\u200c\u200d\u2060\ufeff]/g
const SMART_QUOTES_RE = /[\u201c\u201d\u201e\u201f]/g

function decoderFor(codepage: string): TextDecoder {
  const candidates = [codepage, 'utf-8', 'windows-1252']
  for (const label of candidates) {
    try {
      return new TextDecoder(label)
    } catch {
      /* 试下一个 */
    }
  }
  return new TextDecoder()
}

/**
 * RTF → 可见文本（保留换行与制表，丢掉控制字与样式表）。
 * 非法 RTF 直接抛错，由调用方转成用户可读提示。
 */
export function rtfToText(input: string): string {
  const src = String(input || '').replace(/^\uFEFF/, '')
  if (!/^\s*\{\\rtf/i.test(src)) throw new Error('not-rtf')

  const cpMatch = src.match(/\\ansicpg(\d+)/)
  const codepage =
    cpMatch && cpMatch[1] === '936'
      ? 'gbk'
      : cpMatch && cpMatch[1] === '950'
        ? 'big5'
        : cpMatch && cpMatch[1] === '65001'
          ? 'utf-8'
          : 'windows-1252'
  const decoder = decoderFor(codepage)

  let out = ''
  const pendingBytes: number[] = []
  const flushBytes = () => {
    if (!pendingBytes.length) return
    out += decoder.decode(new Uint8Array(pendingBytes))
    pendingBytes.length = 0
  }
  const emit = (text: string) => {
    flushBytes()
    out += text
  }

  const skipStack: boolean[] = [false]
  let depth = 0
  let ucSkip = 1
  let i = 0

  while (i < src.length) {
    const ch = src[i]

    if (ch === '{') {
      depth += 1
      skipStack[depth] = skipStack[depth - 1]
      i += 1
      continue
    }
    if (ch === '}') {
      if (depth > 0) {
        skipStack.pop()
        depth -= 1
      }
      i += 1
      continue
    }
    if (ch === '\r' || ch === '\n') {
      // RTF 源码里的换行只是排版，不属于正文
      i += 1
      continue
    }
    if (ch !== '\\') {
      if (!skipStack[depth]) {
        flushBytes()
        out += ch
      }
      i += 1
      continue
    }

    // ---- 反斜杠开始的转义 / 控制字 ----
    const next = src[i + 1]
    if (!next) break

    if (next === '\\' || next === '{' || next === '}') {
      if (!skipStack[depth]) emit(next)
      i += 2
      continue
    }
    if (next === "'") {
      const hex = src.slice(i + 2, i + 4)
      if (/^[0-9a-fA-F]{2}$/.test(hex)) {
        if (!skipStack[depth]) pendingBytes.push(parseInt(hex, 16))
        i += 4
        continue
      }
      i += 2
      continue
    }
    if (next === '*') {
      skipStack[depth] = true
      i += 2
      continue
    }
    if (next === '~') {
      if (!skipStack[depth]) emit('\u00a0')
      i += 2
      continue
    }
    if (next === '-') {
      i += 2
      continue
    }
    if (next === '_') {
      if (!skipStack[depth]) emit('-')
      i += 2
      continue
    }
    if (next === '\r' || next === '\n') {
      if (!skipStack[depth]) emit('\n')
      i += 2
      continue
    }
    if (!/[a-zA-Z]/.test(next)) {
      // 未定义的控制符号：按字面输出（如 \"）
      if (!skipStack[depth]) emit(next)
      i += 2
      continue
    }

    // 控制字：\word[-]number[ ]
    let j = i + 1
    while (j < src.length && /[a-zA-Z]/.test(src[j])) j += 1
    const word = src.slice(i + 1, j)
    let sign = 1
    if (src[j] === '-') {
      sign = -1
      j += 1
    }
    let digits = ''
    while (j < src.length && /[0-9]/.test(src[j])) {
      digits += src[j]
      j += 1
    }
    const num = digits ? sign * parseInt(digits, 10) : null
    let hadDelimiter = false
    if (src[j] === ' ') {
      hadDelimiter = true
      j += 1
    }

    if (word === 'bin' && num !== null && num > 0) {
      j += num
      i = j
      continue
    }
    if (word === 'uc' && num !== null) {
      ucSkip = num
      i = j
      continue
    }
    if (word === 'u' && num !== null) {
      if (!skipStack[depth]) {
        const code = num < 0 ? num + 65536 : num
        emit(String.fromCharCode(code))
      }
      // \uN 后面的替换字符：Cocoa 用空格作分隔（hadDelimiter=true），Word 用 `?`。
      if (!hadDelimiter) {
        let left = ucSkip
        while (left > 0 && j < src.length) {
          if (src[j] === '\\' && src[j + 1] === "'") {
            j += 4
            left -= 1
            continue
          }
          if (src[j] === '\\' || src[j] === '{' || src[j] === '}') break
          j += 1
          left -= 1
        }
      } else if (src[j] === '?') {
        // Word 偶尔写成 "\uN ?"：空格是分隔符，问号才是替换字符
        j += 1
      }
      i = j
      continue
    }
    if (SKIP_DESTINATIONS.has(word)) {
      skipStack[depth] = true
      i = j
      continue
    }
    if (!skipStack[depth]) {
      if (word === 'par' || word === 'line' || word === 'sect' || word === 'page') emit('\n')
      else if (word === 'tab') emit('\t')
      else if (word === 'cell' || word === 'row') emit('\n')
      else if (word === 'emdash') emit('\u2014')
      else if (word === 'endash') emit('\u2013')
      else if (word === 'lquote') emit('\u2018')
      else if (word === 'rquote') emit('\u2019')
      else if (word === 'ldblquote') emit('\u201c')
      else if (word === 'rdblquote') emit('\u201d')
      else if (word === 'bullet') emit('\u2022')
    }
    i = j
  }

  flushBytes()
  return out
}

/** 找到从 start 的 `{` 开始、字符串感知的配对 `}`，返回下标（找不到返回 -1）。 */
function matchBrace(text: string, start: number): number {
  let depth = 0
  let inStr = false
  let esc = false
  for (let i = start; i < text.length; i += 1) {
    const ch = text[i]
    if (inStr) {
      if (esc) esc = false
      else if (ch === '\\') esc = true
      else if (ch === '"') inStr = false
      continue
    }
    if (ch === '"') {
      inStr = true
      continue
    }
    if (ch === '{') depth += 1
    else if (ch === '}') {
      depth -= 1
      if (depth === 0) return i
    }
  }
  return -1
}

/** 字符串外空白全删、字符串内只删换行（token 被编辑器拆行时拼回去）。 */
function squashJsonWhitespace(src: string): string {
  let out = ''
  let inStr = false
  let esc = false
  for (let i = 0; i < src.length; i += 1) {
    const ch = src[i]
    if (inStr) {
      if (ch === '\r' || ch === '\n') continue
      out += ch
      if (esc) esc = false
      else if (ch === '\\') esc = true
      else if (ch === '"') inStr = false
      continue
    }
    if (ch === '"') {
      inStr = true
      out += ch
      continue
    }
    if (/\s/.test(ch)) continue
    out += ch
  }
  return out
}

function stripTrailingCommas(src: string): string {
  let out = ''
  let inStr = false
  let esc = false
  for (let i = 0; i < src.length; i += 1) {
    const ch = src[i]
    if (inStr) {
      out += ch
      if (esc) esc = false
      else if (ch === '\\') esc = true
      else if (ch === '"') inStr = false
      continue
    }
    if (ch === '"') {
      inStr = true
      out += ch
      continue
    }
    if (ch === ',') {
      let j = i + 1
      while (j < src.length && /\s/.test(src[j])) j += 1
      if (src[j] === '}' || src[j] === ']') continue
    }
    out += ch
  }
  return out
}

function looksLikeSessionJson(value: unknown): value is Record<string, unknown> {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return false
  const obj = value as Record<string, unknown>
  return Boolean(obj.sessionToken || obj.session_token || obj.accessToken || obj.access_token)
}

/** 把一段文本整理成可解析的 Session JSON；做不到返回空串。 */
function relaxCandidate(raw: string): string {
  const normalized = raw.replace(SMART_QUOTES_RE, '"').trim()
  if (!normalized) return ''

  const attempts = [
    normalized,
    squashJsonWhitespace(normalized),
    stripTrailingCommas(squashJsonWhitespace(normalized)),
  ]
  for (const attempt of attempts) {
    try {
      const parsed = JSON.parse(attempt)
      if (looksLikeSessionJson(parsed)) return attempt
    } catch {
      /* 试下一种 */
    }
  }
  return ''
}

/**
 * 从可见文本中提取 Session JSON（含 sessionToken / accessToken 的对象）。
 * 支持正文里先出现示例、后面才是真实 JSON 的情况；失败返回空串。
 */
export function extractSessionJsonFromText(input: string): string {
  const text = String(input || '').replace(ZERO_WIDTH_RE, '')
  if (!text.trim()) return ''

  const direct = relaxCandidate(text)
  if (direct) return direct

  for (let i = 0; i < text.length; i += 1) {
    if (text[i] !== '{') continue
    const end = matchBrace(text, i)
    if (end <= i) continue
    const hit = relaxCandidate(text.slice(i, end + 1))
    if (hit) return hit
  }

  // 兜底：整体 JSON 被 Word / 翻译工具改坏（引号错位、key 被翻译）时，
  // 只要 token 本身还在，就按字段把它捞出来，组成最小可用 Session。
  const recovered = recoverSessionTokens(text)
  if (recovered) return recovered
  return ''
}

function decodeBase64UrlJson(part: string): Record<string, unknown> | null {
  try {
    const b64 = part.replace(/-/g, '+').replace(/_/g, '/')
    const padded = b64.padEnd(Math.ceil(b64.length / 4) * 4, '=')
    const binary = atob(padded)
    const bytes = Uint8Array.from(binary, (ch) => ch.charCodeAt(0))
    const value = JSON.parse(new TextDecoder().decode(bytes)) as unknown
    return value && typeof value === 'object' && !Array.isArray(value)
      ? (value as Record<string, unknown>)
      : null
  } catch {
    return null
  }
}

function looksLikeJwe(token: string): boolean {
  const parts = token.split('.')
  if (parts.length !== 5) return false
  const header = decodeBase64UrlJson(parts[0])
  return Boolean(header && typeof header.alg === 'string' && typeof header.enc === 'string')
}

function looksLikeJwt(token: string): boolean {
  const parts = token.split('.')
  if (parts.length !== 3) return false
  const payload = decodeBase64UrlJson(parts[1])
  if (!payload) return false
  return (
    typeof payload.exp === 'number' ||
    Boolean(payload['https://api.openai.com/auth']) ||
    typeof payload.email === 'string'
  )
}

/**
 * 从任意文本里按字段捞 token：
 *   - sessionToken 必须是五段 JWE（header 可解码）；
 *   - accessToken 可选，必须是三段 JWT 且 payload 可解码。
 * 只做"取值 + 结构校验"，不改写 token 内容。
 */
function recoverSessionTokens(text: string): string {
  const flat = text.replace(/\s+/g, '')
  const pick = (key: string, valid: (token: string) => boolean): string => {
    const re = new RegExp(`${key}[^A-Za-z0-9]{0,24}(eyJ[A-Za-z0-9_\\-.]{40,})`, 'i')
    const hit = flat.match(re)
    if (!hit) return ''
    return valid(hit[1]) ? hit[1] : ''
  }
  const sessionToken = pick('sessionToken', looksLikeJwe) || pick('session_token', looksLikeJwe)
  if (!sessionToken) return ''
  const accessToken = pick('accessToken', looksLikeJwt) || pick('access_token', looksLikeJwt)
  const out: Record<string, string> = { sessionToken }
  if (accessToken) out.accessToken = accessToken
  return JSON.stringify(out)
}
