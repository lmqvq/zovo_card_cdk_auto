import { describe, expect, it } from 'vitest'
import {
  MAX_SESSION_FILE_BYTES,
  readSessionFile,
} from '../session-file'

// 合法的 JWE 头 + 占位段；JWT 的 payload 可解码（{"exp":…,"email":"a@b.com"}）
const JWE =
  'eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..aXYtMTIzNDU2Nzg5MGFi.Y2lwaGVydGV4dC0xMjM0.dGFnLTEyMzQ1Ng'
const JWT = 'eyJhbGciOiJSUzI1NiJ9.eyJleHAiOjE3OTQzMjIwNzcsImVtYWlsIjoiYUBiLmNvbSJ9.c2lnbmF0dXJl'

function wrapRtf(body: string): string {
  return `{\\rtf1\\ansi\\ansicpg1252\\deff0{\\fonttbl{\\f0 Calibri;}}{\\*\\generator Test;}\\uc1\\pard\\f0\\fs22 ${body}\\par}`
}

function rtfFile(body: string): File {
  return new File([wrapRtf(body)], 'session.rtf', { type: 'application/rtf' })
}

describe('readSessionFile', () => {
  it('reads and trims a TXT Session file', async () => {
    const file = new File(['\n  {"sessionToken":"fixture"}  \n'], 'session.txt', {
      type: 'text/plain',
    })

    await expect(readSessionFile(file)).resolves.toBe('{"sessionToken":"fixture"}')
  })

  it('accepts Session JSON saved with a .docx filename', async () => {
    const file = new File(['\n{"sessionToken":"fixture"}\n'], 'session.docx', {
      type: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    })

    await expect(readSessionFile(file)).resolves.toBe('{"sessionToken":"fixture"}')
  })

  it('reads Session JSON from an RTF document (escaped braces, zero-width break, split line)', async () => {
    // TextEdit 会把长 token 拆行并插入 \u8203 零宽空格；Word 会把 { } 转义成 \{ \}
    const body =
      `\\{"accessToken":"${JWT}","sessionToken":"${JWE.slice(0, 12)}\\u8203 ${JWE.slice(12, 40)}` +
      `\\par ${JWE.slice(40)}"\\}`
    const out = await readSessionFile(rtfFile(body))
    const parsed = JSON.parse(out)
    expect(parsed.sessionToken).toBe(JWE)
    expect(parsed.accessToken).toBe(JWT)
    expect(out).not.toMatch(/[\u200b\u200c\u200d\u2060\ufeff\r\n]/)
  })

  it('recovers tokens when the RTF JSON was mangled by a translator/editor', async () => {
    // key 被翻译、引号错位，整体 JSON 不可解析；token 本身仍完整
    const body = `\\{"令牌"\\ldblquote sessionToken\\rdblquote :\\ldblquote ${JWE}\\rdblquote , broken\\}`
    const out = await readSessionFile(rtfFile(body))
    expect(JSON.parse(out).sessionToken).toBe(JWE)
  })

  it('rejects an RTF without any Session token', async () => {
    const file = rtfFile('\\{ hello world \\}')
    await expect(readSessionFile(file)).rejects.toMatchObject({ code: 'invalid_rtf' })
  })

  it('rejects unsupported files and files above the local limit', async () => {
    const unsupported = new File(['content'], 'session.pdf', { type: 'application/pdf' })
    await expect(readSessionFile(unsupported)).rejects.toMatchObject({
      code: 'unsupported',
    })

    const oversized = new File(['x'.repeat(MAX_SESSION_FILE_BYTES + 1)], 'session.txt', {
      type: 'text/plain',
    })
    await expect(readSessionFile(oversized)).rejects.toMatchObject({
      code: 'too_large',
    })
  })
})
