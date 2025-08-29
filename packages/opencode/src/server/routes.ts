import { z } from "zod";
import { validator } from "../util/validation";

export const withValidation = (schema: z.ZodObject, handler: any) => async (ctx) => {
  try {
    const parsed = validator(schema, ctx);
    return handler({ ...ctx, parsed });
  } catch (error) {
    ctx.throw(400, error.errors || "Validation failed");
  }
};
export const withValidation = (schema: z.ZodObject, handler: any) => async (ctx) => { try { const parsed = validator(schema, ctx); return handler({ ...
ctx, parsed }); } catch (error) { ctx.throw(400, error.errors || "Validation failed"); } };
