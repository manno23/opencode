export function validateUtf8(str: string): boolean {
  try {
    new TextDecoder("utf-8", { fatal: true }).decode(new TextEncoder().encode(str))
    return true
  } catch {
    return false
  }
}

export function strictDecode(bytes: Uint8Array): string {
  return new TextDecoder("utf-8", { fatal: true }).decode(bytes)
}

export function asciiOnly(str: string): boolean {
  for (let i = 0; i < str.length; i++) {
    if (str.charCodeAt(i) > 127) return false
  }
  return true
}

export function enforceAscii(str: string): string {
  return str.replace(/[^\x00-\x7F]/g, "?")
}

export function normalizeNFC(str: string): string {
  return str.normalize("NFC")
}
