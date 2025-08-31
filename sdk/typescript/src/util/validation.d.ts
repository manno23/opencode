import type { z } from "zod"

declare module "../util/validation" {
  /**
• Parses request context using a Zod schema and returns typed data.
• Exact runtime is owned by packages/opencode; this is a type shim.
*/
  export function validator(schema: TSchema, ctx: unknown): z.infer
}
