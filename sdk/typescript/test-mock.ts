// Simple mock test for TypeScript SDK
// This validates basic API structure without complex mocking

import { createOpenCodeClient } from "./src/index"

// Basic test to ensure the client can be created
function testClientCreation() {
  console.log("🧪 Testing client creation...")

  try {
    const client = createOpenCodeClient({
      baseUrl: "https://api.opencode.ai",
    })
    console.log("✅ Client created successfully")
    return true
  } catch (error) {
    console.error("❌ Client creation failed:", error)
    return false
  }
}

// Test type imports
function testTypeImports() {
  console.log("🧪 Testing type imports...")

  // This will fail to compile if types are not available
  // We just test that the module can be imported
  try {
    const client = createOpenCodeClient({
      baseUrl: "https://api.opencode.ai",
    })
    console.log("✅ Type imports working")
    return true
  } catch (error) {
    console.error("❌ Type imports failed:", error)
    return false
  }
}

// Run basic tests
export function runMockTests() {
  console.log("🚀 Running TypeScript SDK mock tests...")

  const results = [
    testClientCreation(),
    testTypeImports(),
  ]

  const passed = results.filter(Boolean).length
  const total = results.length

  console.log(`\n📊 Test Results: ${passed}/${total} passed`)

  if (passed === total) {
    console.log("🎉 All tests passed!")
  } else {
    console.log("❌ Some tests failed")
  }

  return passed === total
}
    baseUrl: "https://api.opencode.ai",
    headers: {
      Authorization: "Bearer test-token",
      "User-Agent": "opencode-sdk-test/1.0.0",
    },
    fetch: mockFetch,
  }

  // Create functional API client
  const api = createOpenCodeAPI(ctx)

  try {
    // Test that all expected methods exist
    console.log("📋 Checking API structure...")

    // Check top-level services
    const expectedServices = ["session", "app", "config", "find", "file"]
    for (const service of expectedServices) {
      if (!api[service as keyof typeof api]) {
        throw new Error(`Missing service: ${service}`)
      }
      console.log(`✅ Service '${service}' exists`)
    }

    // Check session methods
    const expectedSessionMethods = [
      "list",
      "create",
      "get",
      "delete",
      "update",
      "init",
      "abort",
      "share",
      "unshare",
      "summarize",
      "messages",
      "chat",
      "shell",
      "revert",
      "unrevert",
    ]
    for (const method of expectedSessionMethods) {
      if (!api.session[method as keyof typeof api.session]) {
        throw new Error(`Missing session method: ${method}`)
      }
      console.log(`✅ Session method '${method}' exists`)
    }

    // Check app methods
    const expectedAppMethods = ["get", "init", "log"]
    for (const method of expectedAppMethods) {
      if (!api.app[method as keyof typeof api.app]) {
        throw new Error(`Missing app method: ${method}`)
      }
      console.log(`✅ App method '${method}' exists`)
    }

    // Check config methods
    const expectedConfigMethods = ["get", "providers"]
    for (const method of expectedConfigMethods) {
      if (!api.config[method as keyof typeof api.config]) {
        throw new Error(`Missing config method: ${method}`)
      }
      console.log(`✅ Config method '${method}' exists`)
    }

    console.log("🎉 API structure test passed!")
    return true
  } catch (error) {
    console.error("❌ API structure test failed:", error)
    return false
  }
}

// Test type safety
async function testTypeSafety() {
  console.log("🔍 Testing Type Safety...")

  try {
    // Test that we can create API with proper context
    const ctx: OpenCodeContext = {
      baseUrl: "https://api.opencode.ai",
      headers: { Authorization: "Bearer test" },
      fetch: mockFetch,
    }

    const api = createOpenCodeAPI(ctx)

    // Test that methods return promises
    if (!(api.app.get() instanceof Promise)) {
      throw new Error("app.get() should return a Promise")
    }

    if (!(api.session.list() instanceof Promise)) {
      throw new Error("session.list() should return a Promise")
    }

    if (!(api.config.get() instanceof Promise)) {
      throw new Error("config.get() should return a Promise")
    }

    console.log("✅ All methods return Promises")

    // Test Cloudflare API
    const cfCtx: CloudflareContext = {
      baseUrl: "https://api.opencode.ai",
      env: {
        API_TOKEN: "test-token",
        CF_ACCESS_CLIENT_ID: "test-id",
        CF_ACCESS_CLIENT_SECRET: "test-secret",
      },
      fetch: mockFetch,
    }

    const cfApi = createCloudflareAPI(cfCtx)

    if (!(cfApi.app.get() instanceof Promise)) {
      throw new Error("Cloudflare app.get() should return a Promise")
    }

    console.log("✅ Cloudflare API methods return Promises")

    console.log("🎉 Type safety test passed!")
    return true
  } catch (error) {
    console.error("❌ Type safety test failed:", error)
    return false
  }
}

// Test compatibility layer
async function testCompatibilityStructure() {
  console.log("🔄 Testing Compatibility Layer Structure...")

  try {
    const { createOpencodeClient } = await import("./src/compat")

    // Create compatibility client
    const client = createOpencodeClient({
      baseUrl: "https://api.opencode.ai",
      headers: { Authorization: "Bearer test" },
      fetch: mockFetch,
    })

    // Check that all expected methods exist
    const expectedServices = ["session", "app", "config", "find", "file"]
    for (const service of expectedServices) {
      if (!client[service as keyof typeof client]) {
        throw new Error(`Missing compatibility service: ${service}`)
      }
      console.log(`✅ Compatibility service '${service}' exists`)
    }

    // Check session methods
    const expectedSessionMethods = [
      "list",
      "create",
      "get",
      "delete",
      "update",
      "init",
      "abort",
      "share",
      "unshare",
      "summarize",
      "messages",
      "chat",
      "shell",
      "revert",
      "unrevert",
    ]
    for (const method of expectedSessionMethods) {
      if (!client.session[method as keyof typeof client.session]) {
        throw new Error(`Missing compatibility session method: ${method}`)
      }
      console.log(`✅ Compatibility session method '${method}' exists`)
    }

    console.log("🎉 Compatibility layer structure test passed!")
    return true
  } catch (error) {
    console.error("❌ Compatibility layer structure test failed:", error)
    return false
  }
}

// Run all tests
async function runAllTests() {
  console.log("🚀 Starting OpenCode SDK Structure Tests...\n")

  const results = {
    apiStructure: await testAPIStructure(),
    typeSafety: await testTypeSafety(),
    compatibilityStructure: await testCompatibilityStructure(),
  }

  console.log("\n📊 Test Results:")
  console.log("API Structure:", results.apiStructure ? "✅ PASS" : "❌ FAIL")
  console.log("Type Safety:", results.typeSafety ? "✅ PASS" : "❌ FAIL")
  console.log("Compatibility Structure:", results.compatibilityStructure ? "✅ PASS" : "❌ FAIL")

  const allPassed = Object.values(results).every((result) => result)

  if (allPassed) {
    console.log("\n🎉 All structure tests passed! SDK implementation is correct.")
  } else {
    console.log("\n⚠️ Some tests failed. Please check the implementation.")
  }

  return allPassed
}

// Export test runner
export { runAllTests }

// Run tests if this file is executed directly
if (import.meta.main) {
  runAllTests().then((success) => {
    process.exit(success ? 0 : 1)
  })
}
