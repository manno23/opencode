import { Provider } from "../provider/provider"
import type { Message } from "../session/protocol/types"

// Model configurations for the two-agent loop
export const MODEL_CONFIGS = {
  executor_small: {
    provider: "opencode", // or local provider
    model: "grok-code", // or qwen2.5-coder-7b-instruct
    temperature: 0.2,
    top_p: 0.9,
    max_output_tokens: 1200,
  },
  supervisor_gpt5: {
    provider: "openai",
    model: "gpt-5-thinking",
    temperature: 0.1,
    top_p: 0.8,
    max_output_tokens: 2000,
  },
} as const

export type ModelKey = keyof typeof MODEL_CONFIGS

/**
 * Call a model with the given prompt and return the response
 */
export async function callModel(
  modelKey: ModelKey,
  options: {
    system?: string
    user?: string
    tool?: Message
  },
): Promise<Message> {
  const config = MODEL_CONFIGS[modelKey]
  const { language } = await Provider.getModel(config.provider, config.model)

  // Build the prompt
  let prompt = ""
  if (options.system) {
    prompt += `System: ${options.system}\n\n`
  }
  if (options.user) {
    prompt += `User: ${options.user}\n\n`
  }
  if (options.tool) {
    prompt += `Tool: ${JSON.stringify(options.tool)}\n\n`
  }

  // Generate response
  const result = await language.generateText({
    prompt,
    temperature: config.temperature,
    topP: config.top_p,
    maxTokens: config.max_output_tokens,
  })

  // Parse JSON response
  try {
    return JSON.parse(result.text)
  } catch (error) {
    throw new Error(`Failed to parse model response as JSON: ${result.text}`)
  }
}
