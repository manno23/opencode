import { validateUtf8, asciiOnly, normalizeNFC } from "../../../../opencode/src/util/encoding"

describe("TUI message rendering", () => {
  it("should be valid UTF-8, ASCII-only, and NFC-normalized", () => {
    const sampleRendering = "Hello, world! This is a sample TUI message."

    expect(validateUtf8(sampleRendering)).toBe(true)
    expect(asciiOnly(sampleRendering)).toBe(true)
    expect(normalizeNFC(sampleRendering)).toBe(sampleRendering)
  })
})
