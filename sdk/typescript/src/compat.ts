import { createOpenCodeClient } from "./generated/client.js"

/**
 * Compatibility layer that adapts the TypeScript SDK to match the JS SDK interface
 * This maps the existing object-style API to the new generated client
 */
export interface OpencodeClientCompat {
  session: {
    list(): Promise<any>
    create(data?: any): Promise<any>
    get(id: string): Promise<any>
    delete(id: string): Promise<any>
    update(id: string, data: any): Promise<any>
    init(id: string): Promise<any>
    abort(id: string): Promise<any>
    unshare(id: string): Promise<any>
    share(id: string, data: any): Promise<any>
    summarize(id: string): Promise<any>
    messages(id: string): Promise<any>
    chat(id: string, data: any): Promise<any>
    shell(id: string, data: any): Promise<any>
    revert(id: string, data: any): Promise<any>
    unrevert(id: string): Promise<any>
  }
  app: {
    get(): Promise<any>
    init(): Promise<any>
    log(data: any): Promise<any>
  }
  config: {
    get(): Promise<any>
    providers(): Promise<any>
  }
  find: {
    text(data?: any): Promise<any>
    files(data?: any): Promise<any>
    symbols(data?: any): Promise<any>
  }
  file: {
    get(data?: any): Promise<any>
    status(data?: any): Promise<any>
  }
}

/**
 * Create a compatibility client that matches the JS SDK interface
 * Maps to the actual OpenAPI paths from the generated types
 */
export function createOpencodeClient(
  options: {
    baseUrl?: string
    headers?: Record<string, string>
    fetch?: typeof fetch
  } = {},
): OpencodeClientCompat {
  const client = createOpenCodeClient(options)

  return {
    session: {
      list: () => client.GET("/session"),
      create: (data) => client.POST("/session", { body: data }),
      get: (id) => client.GET("/session/{id}", { params: { path: { id } } }),
      delete: (id) => client.DELETE("/session/{id}", { params: { path: { id } } }),
      update: (id, data) => client.PATCH("/session/{id}", { params: { path: { id } }, body: data }),
      init: (id) => client.POST("/session/{id}/init", { params: { path: { id } } }),
      abort: (id) => client.POST("/session/{id}/abort", { params: { path: { id } } }),
      unshare: (id) => client.DELETE("/session/{id}/share", { params: { path: { id } } }),
      share: (id, data) => client.POST("/session/{id}/share", { params: { path: { id } }, body: data }),
      summarize: (id) => client.POST("/session/{id}/summarize", { params: { path: { id } } }),
      messages: (id) => client.GET("/session/{id}/message", { params: { path: { id } } }),
      chat: (id, data) => client.POST("/session/{id}/message", { params: { path: { id } }, body: data }),
      shell: (id, data) => client.POST("/session/{id}/shell", { params: { path: { id } }, body: data }),
      revert: (id, data) => client.POST("/session/{id}/revert", { params: { path: { id } }, body: data }),
      unrevert: (id) => client.POST("/session/{id}/unrevert", { params: { path: { id } } }),
    },
    app: {
      get: () => client.GET("/app"),
      init: () => client.POST("/app/init"),
      log: (data) => client.POST("/log", { body: data }),
    },
    config: {
      get: () => client.GET("/config"),
      providers: () => client.GET("/config/providers"),
    },
    find: {
      text: (data) => client.GET("/find", { params: { query: data } }),
      files: (data) => client.GET("/find/file", { params: { query: data } }),
      symbols: (data) => client.GET("/find/symbol", { params: { query: data } }),
    },
    file: {
      get: (data) => client.GET("/file", { params: { query: data } }),
      status: (data) => client.GET("/file/status", { params: { query: data } }),
    },
  }
}
