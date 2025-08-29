// Modern TypeScript SDK for OpenCode API
// Optimized for Cloudflare Workers and edge runtimes

export { createOpenCodeClient } from "./generated/client.js"
export type { OpenCodeClient } from "./generated/client.js"
export type * from "./generated/types.js"

// Re-export commonly used types for convenience
export type { paths, components, operations, webhooks } from "./generated/types.js"
