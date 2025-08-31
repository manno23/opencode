export interface ClientOptions {
  baseUrl?: string
  headers?: Record<string, string>
  fetch?: typeof fetch
}

export class OpenCodeClient {
  private baseUrl: string
  private headers: Record<string, string>
  private fetch: typeof fetch

  constructor(options: ClientOptions = {}) {
    this.baseUrl = options.baseUrl || "http://localhost:3000"
    this.headers = options.headers || {}
    this.fetch = options.fetch || fetch
  }

  async request<T>(method: string, path: string, body?: any): Promise<T> {
    const url = `${this.baseUrl}${path}`
    const headers = {
      "Content-Type": "application/json",
      "User-Agent": "opencode-sdk-typescript/1.0.0",
      ...this.headers,
    }
    const resp = await this.fetch(url, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    })
    if (!resp.ok) throw new Error(`HTTP ${resp.status}`)
    return resp.json()
  }

  async getApp(): Promise<any> {
    return this.request("GET", "/app")
  }
  async initApp(body: any): Promise<any> {
    return this.request("POST", "/app/init", body)
  }
  async subscribeToEvents(): Promise<any> {
    return this.request("GET", "/event")
  }
}

export function createOpenCodeClient(options: ClientOptions = {}): OpenCodeClient {
  return new OpenCodeClient(options)
}
