// Functional API Wrapper with Explicit Context Pattern
// This implements the new functional SDK design on top of generated clients

import { createOpenCodeClient } from "./generated/client.js"
import { z } from "zod"

// ============================================================================
// CONTEXT INTERFACE - Explicit dependency injection
// ============================================================================

export interface OpenCodeContext {
  baseUrl: string
  headers?: Record<string, string>
  fetch?: typeof fetch
}

// ============================================================================
// ZOD SCHEMAS - Enhanced validation layer on generated types
// ============================================================================

// Session schemas with enhanced validation
const SessionSchema = z.object({
  id: z.string(),
  parentID: z.string().optional(),
  share: z
    .object({
      url: z.string(),
    })
    .optional(),
  title: z.string().min(1),
  version: z.string(),
  time: z.object({
    created: z.number(),
    updated: z.number(),
  }),
  revert: z
    .object({
      messageID: z.string(),
      partID: z.string().optional(),
      snapshot: z.string().optional(),
      diff: z.string().optional(),
    })
    .optional(),
})

const SessionListResponseSchema = z.object({
  sessions: z.array(SessionSchema),
  total: z.number().int().min(0),
})

// Create session has no request body based on OpenAPI spec
const CreateSessionRequestSchema = z.object({}).optional()

const ChatMessageRequestSchema = z.object({
  messageID: z.string().optional(),
  providerID: z.string(),
  modelID: z.string(),
  agent: z.string().optional(),
  system: z.string().optional(),
  tools: z.record(z.boolean()).optional(),
  parts: z.array(z.any()), // TextPartInput | FilePartInput | AgentPartInput
})

// App schemas
const AppInfoSchema = z.object({
  name: z.string().min(1),
  version: z.string().min(1),
})

// Config schemas
const AppConfigSchema = z.object({
  providers: z.array(
    z.object({
      name: z.string(),
      enabled: z.boolean(),
      config: z.record(z.unknown()),
    }),
  ),
  models: z.array(
    z.object({
      id: z.string(),
      name: z.string(),
      provider: z.string(),
    }),
  ),
})

// ============================================================================
// RESPONSE HANDLER - Consistent error handling and validation
// ============================================================================

const handleResponse = async <T>(
  response: { data?: unknown; error?: unknown },
  schema?: z.ZodSchema<T>,
): Promise<T> => {
  if (response.error) {
    throw new Error(`API Error: ${response.error}`)
  }

  if (!response.data) {
    throw new Error("No data received from API")
  }

  return schema ? schema.parse(response.data) : (response.data as T)
}

// ============================================================================
// FUNCTIONAL API CLIENT - Explicit context pattern
// ============================================================================

export const createOpenCodeAPI = (ctx: OpenCodeContext) => {
  // Create underlying generated client
  const client = createOpenCodeClient(ctx)

  return {
    // Session operations - functional style with explicit context
    session: {
      // List all sessions
      list: async () => {
        const response = await client.GET("/session")
        return handleResponse(response, SessionListResponseSchema)
      },

      // Create a new session (no request body based on OpenAPI spec)
      create: async () => {
        const response = await client.POST("/session")
        return handleResponse(response, SessionSchema)
      },

      // Get a specific session
      get: async (id: string) => {
        const response = await client.GET("/session/{id}", {
          params: { path: { id } },
        })
        return handleResponse(response, SessionSchema)
      },

      // Delete a session
      delete: async (id: string) => {
        const response = await client.DELETE("/session/{id}", {
          params: { path: { id } },
        })
        return handleResponse(response)
      },

      // Update session properties
      update: async (id: string, data: any) => {
        const response = await client.PATCH("/session/{id}", {
          params: { path: { id } },
          body: data,
        })
        return handleResponse(response, SessionSchema)
      },

      // Initialize session analysis
      init: async (id: string) => {
        const response = await client.POST("/session/{id}/init", {
          params: { path: { id } },
        })
        return handleResponse(response)
      },

      // Abort a session
      abort: async (id: string) => {
        const response = await client.POST("/session/{id}/abort", {
          params: { path: { id } },
        })
        return handleResponse(response)
      },

      // Share a session
      share: async (id: string, data: any) => {
        const response = await client.POST("/session/{id}/share", {
          params: { path: { id } },
          body: data,
        })
        return handleResponse(response)
      },

      // Unshare a session
      unshare: async (id: string) => {
        const response = await client.DELETE("/session/{id}/share", {
          params: { path: { id } },
        })
        return handleResponse(response)
      },

      // Summarize session
      summarize: async (id: string) => {
        const response = await client.POST("/session/{id}/summarize", {
          params: { path: { id } },
        })
        return handleResponse(response)
      },

      // Get session messages
      messages: async (id: string) => {
        const response = await client.GET("/session/{id}/message", {
          params: { path: { id } },
        })
        return handleResponse(response)
      },

      // Send chat message
      chat: async (id: string, params: z.infer<typeof ChatMessageRequestSchema>) => {
        const validatedParams = ChatMessageRequestSchema.parse(params)
        const response = await client.POST("/session/{id}/message", {
          params: { path: { id } },
          body: validatedParams,
        })
        return handleResponse(response)
      },

      // Run shell command
      shell: async (id: string, data: any) => {
        const response = await client.POST("/session/{id}/shell", {
          params: { path: { id } },
          body: data,
        })
        return handleResponse(response)
      },

      // Revert message
      revert: async (id: string, data: any) => {
        const response = await client.POST("/session/{id}/revert", {
          params: { path: { id } },
          body: data,
        })
        return handleResponse(response)
      },

      // Restore reverted messages
      unrevert: async (id: string) => {
        const response = await client.POST("/session/{id}/unrevert", {
          params: { path: { id } },
        })
        return handleResponse(response)
      },
    },

    // App operations
    app: {
      // Get app info
      get: async () => {
        const response = await client.GET("/app")
        return handleResponse(response, AppInfoSchema)
      },

      // Initialize app
      init: async () => {
        const response = await client.POST("/app/init")
        return handleResponse(response)
      },

      // Write log entry
      log: async (data: any) => {
        const response = await client.POST("/log", { body: data })
        return handleResponse(response)
      },
    },

    // Config operations
    config: {
      // Get configuration
      get: async () => {
        const response = await client.GET("/config")
        return handleResponse(response, AppConfigSchema)
      },

      // Get providers
      providers: async () => {
        const response = await client.GET("/config/providers")
        return handleResponse(response)
      },
    },

    // Find operations
    find: {
      // Find text in files
      text: async (params?: any) => {
        const response = await client.GET("/find", {
          params: { query: params },
        })
        return handleResponse(response)
      },

      // Find files
      files: async (params?: any) => {
        const response = await client.GET("/find/file", {
          params: { query: params },
        })
        return handleResponse(response)
      },

      // Find symbols
      symbols: async (params?: any) => {
        const response = await client.GET("/find/symbol", {
          params: { query: params },
        })
        return handleResponse(response)
      },
    },

    // File operations
    file: {
      // Read file
      read: async (params?: any) => {
        const response = await client.GET("/file", {
          params: { query: params },
        })
        return handleResponse(response)
      },

      // Get file status
      status: async (params?: any) => {
        const response = await client.GET("/file/status", {
          params: { query: params },
        })
        return handleResponse(response)
      },
    },
  }
}

// ============================================================================
// CLOUDFLARE WORKERS OPTIMIZATION
// ============================================================================

// Cloudflare-specific context extension
export interface CloudflareContext extends OpenCodeContext {
  env: {
    API_TOKEN: string
    CF_ACCESS_CLIENT_ID?: string
    CF_ACCESS_CLIENT_SECRET?: string
  }
}

// Create Cloudflare-optimized API client
export const createCloudflareAPI = (ctx: CloudflareContext) => {
  const cloudflareHeaders: Record<string, string> = {
    Authorization: `Bearer ${ctx.env.API_TOKEN}`,
    "User-Agent": "opencode-sdk-cloudflare/1.0.0",
  }

  // Add Cloudflare Access headers if available
  if (ctx.env.CF_ACCESS_CLIENT_ID && ctx.env.CF_ACCESS_CLIENT_SECRET) {
    cloudflareHeaders["CF-Access-Client-Id"] = ctx.env.CF_ACCESS_CLIENT_ID
    cloudflareHeaders["CF-Access-Client-Secret"] = ctx.env.CF_ACCESS_CLIENT_SECRET
  }

  return createOpenCodeAPI({
    ...ctx,
    headers: {
      ...cloudflareHeaders,
      ...ctx.headers,
    },
    fetch: ctx.fetch || fetch.bind(globalThis),
  })
}

// ============================================================================
// TYPE EXPORTS
// ============================================================================

// Export inferred types from Zod schemas for external use
export type Session = z.infer<typeof SessionSchema>
export type CreateSessionRequest = z.infer<typeof CreateSessionRequestSchema>
export type ChatMessageRequest = z.infer<typeof ChatMessageRequestSchema>
export type SessionListResponse = z.infer<typeof SessionListResponseSchema>
export type AppInfo = z.infer<typeof AppInfoSchema>
export type AppConfig = z.infer<typeof AppConfigSchema>

// Export functional API type
export type OpenCodeAPI = ReturnType<typeof createOpenCodeAPI>
