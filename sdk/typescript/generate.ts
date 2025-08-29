#!/usr/bin/env bun
import { existsSync } from "fs"
import { mkdir, writeFile } from "fs/promises"
import { execSync } from "child_process"

const PROJECT_ROOT = "../.."
const SPEC_PATH = `${PROJECT_ROOT}/schema/openapi.json`
const OUTPUT_DIR = "src/generated"

async function generateOpenAPISpec() {
  console.log("🔄 Generating OpenAPI spec from opencode server...")

  try {
    // Change to project root and generate spec
    process.chdir(PROJECT_ROOT)
    execSync("bun run --conditions=development packages/opencode/src/index.ts generate", {
      stdio: "pipe",
    })

    if (!existsSync("schema/openapi.json")) {
      throw new Error("OpenAPI spec generation failed - no output file")
    }

    console.log("✅ OpenAPI spec generated successfully")
  } catch (error) {
    console.error("❌ Failed to generate OpenAPI spec:", error)
    throw error
  } finally {
    // Change back to typescript directory
    process.chdir("sdk/typescript")
  }
}

async function generateTypeScriptTypes() {
  console.log("🔄 Generating TypeScript types...")

  if (!existsSync(OUTPUT_DIR)) {
    await mkdir(OUTPUT_DIR, { recursive: true })
  }

  try {
    execSync(`npx openapi-typescript ${SPEC_PATH} --output ${OUTPUT_DIR}/types.ts`, {
      stdio: "inherit",
    })
    console.log("✅ TypeScript types generated successfully")
  } catch (error) {
    console.error("❌ Failed to generate TypeScript types:", error)
    throw error
  }
}

async function generateSDKClient() {
  console.log("🔄 Generating SDK client...")

  const clientCode = `import createClient from 'openapi-fetch'
import type { paths } from './types.js'

export type OpenCodeClient = ReturnType<typeof createClient<paths>>

export function createOpenCodeClient(options: {
  baseUrl?: string
  headers?: Record<string, string>
  fetch?: typeof fetch
} = {}): OpenCodeClient {
  const {
    baseUrl = 'http://localhost:3000',
    headers = {},
    fetch: customFetch
  } = options

  return createClient<paths>({
    baseUrl,
    headers: {
      'User-Agent': 'opencode-sdk-typescript/0.1.0',
      ...headers
    },
    ...(customFetch && { fetch: customFetch })
  })
}

export * from './types.js'
`

  await writeFile(`${OUTPUT_DIR}/client.ts`, clientCode)
  console.log("✅ SDK client generated successfully")
}

async function generateIndexFile() {
  console.log("🔄 Generating index file...")

  const indexCode = `// Modern TypeScript SDK for OpenCode API
// Optimized for Cloudflare Workers and edge runtimes

export { createOpenCodeClient } from './generated/client.js'
export type { OpenCodeClient } from './generated/client.js'
export type * from './generated/types.js'

// Re-export commonly used types for convenience
export type {
  paths,
  components,
  operations,
  webhooks,
  external
} from './generated/types.js'
`

  if (!existsSync("src")) {
    await mkdir("src", { recursive: true })
  }

  await writeFile("src/index.ts", indexCode)
  console.log("✅ Index file generated successfully")
}

async function generateTypeScriptConfig() {
  console.log("🔄 Generating TypeScript config...")

  const tsConfig = {
    compilerOptions: {
      target: "ES2022",
      module: "ESNext",
      moduleResolution: "bundler",
      declaration: true,
      declarationMap: true,
      esModuleInterop: true,
      allowSyntheticDefaultImports: true,
      strict: true,
      skipLibCheck: true,
      forceConsistentCasingInFileNames: true,
      resolveJsonModule: true,
      isolatedModules: true,
      noEmit: false,
      outDir: "dist",
    },
    include: ["src/**/*"],
    exclude: ["node_modules", "dist", "**/*.test.ts"],
  }

  await writeFile("tsconfig.json", JSON.stringify(tsConfig, null, 2))
  console.log("✅ TypeScript config generated successfully")
}

async function main() {
  console.log("🚀 Starting TypeScript SDK generation...")

  try {
    await generateOpenAPISpec()
    await generateTypeScriptTypes()
    await generateSDKClient()
    await generateIndexFile()
    await generateTypeScriptConfig()

    console.log("🎉 TypeScript SDK generation completed successfully!")
    console.log("")
    console.log("Next steps:")
    console.log("  1. cd sdk/typescript")
    console.log("  2. bun install")
    console.log("  3. bun run build")
    console.log("  4. bun test (once tests are added)")
  } catch (error) {
    console.error("💥 SDK generation failed:", error)
    process.exit(1)
  }
}

if (import.meta.main) {
  main()
}
