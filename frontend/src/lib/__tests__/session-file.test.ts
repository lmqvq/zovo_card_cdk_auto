import { describe, expect, it } from 'vitest'
import {
  MAX_SESSION_FILE_BYTES,
  readSessionFile,
} from '../session-file'

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
