import type { ModelMessage } from "ai"
import { unique } from "remeda"
import { Provider } from "./provider"

export namespace ProviderTransform {
  export function normalizeToolCallIds(msgs: ModelMessage[]): ModelMessage[] {
    return msgs.map((msg) => {
      if ((msg.role === "assistant" || msg.role === "tool") && Array.isArray(msg.content)) {
        msg.content = msg.content.map((part) => {
          if ((part.type === "tool-call" || part.type === "tool-result") && "toolCallId" in part) {
            return {
              ...part,
              toolCallId: part.toolCallId.replace(/[^a-zA-Z0-9_-]/g, "_"),
            }
          }
          return part
        })
      }
      return msg
    })
  }

  export function applyCaching(msgs: ModelMessage[], providerID: string): ModelMessage[] {
    const system = msgs.filter((msg) => msg.role === "system").slice(0, 2)
    const final = msgs.filter((msg) => msg.role !== "system").slice(-2)

    const providerOptions = {
      anthropic: {
        cacheControl: { type: "ephemeral" },
      },
      openrouter: {
        cache_control: { type: "ephemeral" },
      },
      bedrock: {
        cachePoint: { type: "ephemeral" },
      },
      openaiCompatible: {
        cache_control: { type: "ephemeral" },
      },
    }

    for (const msg of unique([...system, ...final])) {
      const shouldUseContentOptions = providerID !== "anthropic" && Array.isArray(msg.content) && msg.content.length > 0

      if (shouldUseContentOptions) {
        const lastContent = msg.content[msg.content.length - 1]
        if (lastContent && typeof lastContent === "object") {
          lastContent.providerOptions = {
            ...lastContent.providerOptions,
            ...providerOptions,
          }
          continue
        }
      }

      msg.providerOptions = {
        ...msg.providerOptions,
        ...providerOptions,
      }
    }

    return msgs
  }

  export async function optimizePrompts(
    msgs: ModelMessage[],
    _providerID: string,
    _modelID: string,
  ): Promise<ModelMessage[]> {
    // Check if prompt optimization is enabled in config
    const config = await import("../config/config").then((m) => m.Config.get())
    const optimizationEnabled = config.experimental?.promptOptimization?.enabled

    if (!optimizationEnabled) {
      return msgs
    }

    // Get the optimizing model
    let optimizingModel
    try {
      const optimizerModel = config.experimental?.promptOptimization?.model
      if (optimizerModel) {
        optimizingModel = await Provider.getModel(optimizerModel.providerID, optimizerModel.modelID)
      } else {
        // Default to GPT-4 if available
        optimizingModel = await Provider.getModel("openai", "gpt-4")
      }
    } catch {
      // If optimizing model not available, skip optimization
      return msgs
    }

    const optimizedMsgs = await Promise.all(
      msgs.map(async (msg) => {
        if (msg.role === "user" && typeof msg.content === "string") {
          try {
            const optimizationPrompt = `You are an expert prompt engineer. Optimize the following prompt to make it more effective for an AI assistant:

Original prompt: ${msg.content}

Provide an optimized version that is clearer, more specific, and more likely to produce high-quality responses. Return only the optimized prompt without any explanation.`

            const result = await optimizingModel.language.generateText({
              prompt: optimizationPrompt,
              temperature: 0.3,
              maxTokens: 1000,
            })

            return {
              ...msg,
              content: result.text.trim(),
            }
          } catch (error) {
            // If optimization fails, return original message
            return msg
          }
        }
        return msg
      }),
    )

    return optimizedMsgs
  }

  export async function message(msgs: ModelMessage[], providerID: string, modelID: string) {
    if (modelID.includes("claude")) {
      msgs = normalizeToolCallIds(msgs)
    }
    if (providerID === "anthropic" || modelID.includes("anthropic") || modelID.includes("claude")) {
      msgs = applyCaching(msgs, providerID)
    }

    // Apply prompt optimization if enabled
    msgs = await optimizePrompts(msgs, providerID, modelID)

    return msgs
  }

  export function temperature(_providerID: string, modelID: string) {
    if (modelID.toLowerCase().includes("qwen")) return 0.55
    if (modelID.toLowerCase().includes("claude")) return 1
    return 0
  }

  export function topP(_providerID: string, modelID: string) {
    if (modelID.toLowerCase().includes("qwen")) return 1
    return undefined
  }

  export function options(providerID: string, modelID: string, sessionID: string): Record<string, any> | undefined {
    const result: Record<string, any> = {}

    if (providerID === "openai") {
      result["promptCacheKey"] = sessionID
    }

    if (modelID.includes("gpt-5") && !modelID.includes("gpt-5-chat")) {
      result["reasoningEffort"] = "high"
      if (providerID !== "azure") {
        result["textVerbosity"] = "low"
      }
    }
    return result
  }
}
