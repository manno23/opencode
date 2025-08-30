import createClient from "openapi-fetch"
import type { paths } from "./types.js"

export type OpencodeClient = ReturnType<typeof createClient<paths>>

export function createOpencodeClient(
  options: {
    baseUrl?: string
    headers?: Record<string, string>
    fetch?: typeof fetch
  } = {},
): OpencodeClient {
  const { baseUrl = "http://localhost:4096", headers = {}, fetch: customFetch } = options

  return createClient<paths>({
    baseUrl,
    headers: {
      "User-Agent": "opencode-sdk-typescript/1.0.0",
      ...headers,
    },
    ...(customFetch && { fetch: customFetch }),
  })
}

export * from "./types.js"
