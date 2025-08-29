// Test file for the generated OpenAPI client
// This validates the actual client structure and methods

import { createOpenCodeClient } from "./src/index"

// Test the generated client API
async function testGeneratedClient() {
  console.log("🧪 Testing Generated Client...")

  // Create API client
  const client = createOpenCodeClient({
    baseUrl: "https://api.opencode.ai",
    headers: {
      Authorization: "Bearer test-token",
      "User-Agent": "opencode-sdk-test/1.0.0",
    },
  })

  try {
    // Test simple operations using the actual client methods
    console.log("📋 Testing app info...")
    const appResponse = await client.GET("/app")
    console.log("✅ App response:", appResponse)

    console.log("⚙️ Testing config providers...")
    const configResponse = await client.GET("/config/providers")
    console.log("✅ Config response:", configResponse)

    console.log("📝 Testing session list...")
    const sessionsResponse = await client.GET("/session")
    console.log("✅ Sessions response:", sessionsResponse)

    console.log("➕ Testing session creation...")
    const createResponse = await client.POST("/session")
    console.log("✅ Create response:", createResponse)

    console.log("🎉 Generated client test completed successfully!")
    return true
  } catch (error) {
    console.error("❌ Generated client test failed:", error)
    return false
  }
}

// Run all tests
export async function runFunctionalTests() {
  console.log("🚀 Running Generated Client Tests...")

  const results = [
    testClientCreation(),
    await testGeneratedClient(),
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
}

// Test client creation and configuration
function testClientCreation() {
  console.log("🔧 Testing client creation...")

  try {
    // Test default client
    const defaultClient = createOpenCodeClient()
    console.log("✅ Default client created")

    // Test client with custom options
    const customClient = createOpenCodeClient({
      baseUrl: "https://custom.api.com",
      headers: {
        "X-Custom": "test",
      },
    })
    console.log("✅ Custom client created")

    return true
  } catch (error) {
    console.error("❌ Client creation test failed:", error)
    return false
  }
}

// Run all tests
export async function runFunctionalTests() {
  console.log("🚀 Running Generated Client Tests...")

  const results = [
    testClientCreation(),
    await testGeneratedClient(),
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

    console.log("🎉 All functional API tests passed!")
    return true
  } catch (error) {
    console.error("❌ Functional API test failed:", error)
    return false
  }
}

// Test Cloudflare-optimized API
async function testCloudflareAPI() {
  console.log("☁️ Testing Cloudflare API...")

  // Create Cloudflare context
  const cfCtx: CloudflareContext = {
    baseUrl: "https://api.opencode.ai",
    env: {
      API_TOKEN: "cf-test-token",
      CF_ACCESS_CLIENT_ID: "test-client-id",
      CF_ACCESS_CLIENT_SECRET: "test-client-secret",
    },
  }

  // Create Cloudflare API client
  const cfApi = createCloudflareAPI(cfCtx)

  try {
    console.log("📋 Testing Cloudflare app info...")
    const appInfo = await cfApi.app.get()
    console.log("✅ Cloudflare app info:", appInfo)

    console.log("🎉 Cloudflare API test passed!")
    return true
  } catch (error) {
    console.error("❌ Cloudflare API test failed:", error)
    return false
  }
}

// Test compatibility layer
async function testCompatibilityLayer() {
  console.log("🔄 Testing Compatibility Layer...")

  const { createOpencodeClient } = await import("./src/compat")

  try {
    // Create compatibility client
    const client = createOpencodeClient({
      baseUrl: "https://api.opencode.ai",
      headers: {
        Authorization: "Bearer test-token",
      },
    })

    console.log("📋 Testing compatibility app info...")
    const appInfo = await client.app.get()
    console.log("✅ Compatibility app info:", appInfo)

    console.log("📝 Testing compatibility session list...")
    const sessions = await client.session.list()
    console.log("✅ Compatibility sessions:", sessions)

    console.log("🎉 Compatibility layer test passed!")
    return true
  } catch (error) {
    console.error("❌ Compatibility layer test failed:", error)
    return false
  }
}

// Run all tests
async function runAllTests() {
  console.log("🚀 Starting OpenCode SDK Tests...\n")

  const results = {
    functionalAPI: await testFunctionalAPI(),
    cloudflareAPI: await testCloudflareAPI(),
    compatibilityLayer: await testCompatibilityLayer(),
  }

  console.log("\n📊 Test Results:")
  console.log("Functional API:", results.functionalAPI ? "✅ PASS" : "❌ FAIL")
  console.log("Cloudflare API:", results.cloudflareAPI ? "✅ PASS" : "❌ FAIL")
  console.log("Compatibility Layer:", results.compatibilityLayer ? "✅ PASS" : "❌ FAIL")

  const allPassed = Object.values(results).every((result) => result)

  if (allPassed) {
    console.log("\n🎉 All tests passed! SDK is ready for production.")
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
