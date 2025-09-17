import { z } from "zod"
import { validator } from "../util/validation"

export const withValidation = (schema: z.ZodObject<any, any, any>, handler: any) => async (ctx: any) => {
  try {
    const parsed = validator(schema, ctx)
    return handler({ ...ctx, parsed })
  } catch (error: any) {
    ctx.throw(400, error.errors || "Validation failed")
  }
}
