import { z } from "zod"
import { Tool } from "./tool"
import { Provider } from "../provider/provider"
import DESCRIPTION from "./optimize-prompt.txt"

export const OptimizePromptTool = Tool.define("optimize_prompt", {
  description: DESCRIPTION,
  parameters: z.object({
    prompt: z.string().describe("The original prompt to optimize"),
    context: z.string().optional().describe("Additional context about how the prompt will be used"),
    targetModel: z.string().optional().describe("The target model the optimized prompt will be sent to"),
    optimizationGoal: z
      .enum(["clarity", "specificity", "efficiency", "comprehensiveness"])
      .optional()
      .describe("The type of optimization to apply"),
  }),
  async execute(params, ctx) {
    const { prompt, context, targetModel, optimizationGoal = "clarity" } = params

    // Get the optimizing model (default to GPT-4 if available, otherwise use default)
    let optimizingModel
    try {
      optimizingModel = await Provider.getModel("openai", "gpt-4")
    } catch {
      // Fallback to default model if GPT-4 not available
      const defaultModel = await Provider.defaultModel()
      optimizingModel = await Provider.getModel(defaultModel.providerID, defaultModel.modelID)
    }

    // Create optimization prompt
    const optimizationPrompt = `You are an expert prompt engineer. Your task is to optimize the following prompt to make it more effective.

Original Prompt:
${prompt}

${context ? `Context: ${context}` : ""}

Optimization Goal: ${optimizationGoal}
- clarity: Make the prompt clearer and more understandable
- specificity: Make the prompt more specific and detailed
- efficiency: Make the prompt more concise while maintaining effectiveness
- comprehensiveness: Make the prompt more thorough and complete

${targetModel ? `Target Model: ${targetModel}` : ""}

Please provide an optimized version of this prompt that achieves the specified goal. Return only the optimized prompt without any additional explanation or formatting.`

    try {
      const result = await optimizingModel.language.generateText({
        prompt: optimizationPrompt,
        temperature: 0.3,
        maxTokens: 2000,
      })

      const optimizedPrompt = result.text.trim()

      ctx.metadata({
        title: `Optimized prompt (${optimizationGoal})`,
        metadata: {
          originalLength: prompt.length,
          optimizedLength: optimizedPrompt.length,
          optimizationGoal,
          targetModel: targetModel || "default",
        },
      })

      return {
        title: `Prompt optimized for ${optimizationGoal}`,
        metadata: {
          originalLength: prompt.length,
          optimizedLength: optimizedPrompt.length,
          optimizationGoal,
          targetModel: targetModel || "default",
        },
        output: optimizedPrompt,
      }
    } catch (error) {
      throw new Error(`Failed to optimize prompt: ${error instanceof Error ? error.message : "Unknown error"}`)
    }
  },
})
