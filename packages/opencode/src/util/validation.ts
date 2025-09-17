import { z } from "zod"

/**
 * Parses request context using a Zod schema and returns typed data.
 * @param schema The Zod schema to validate against
 * @param ctx The context to validate
 * @returns The parsed and validated data
 */
export function validator<T extends z.ZodTypeAny>(schema: T, ctx: any): z.infer<T> {
  // This is a simplified implementation - in reality, this would extract
  // the relevant data from the context (like request body, query params, etc.)
  // and validate it against the schema
  const data = ctx.req ? ctx.req.body || ctx.req.query || {} : ctx
  return schema.parse(data)
}
