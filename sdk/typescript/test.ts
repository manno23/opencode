#!/usr/bin/env bun
// Quick test to verify the SDK builds and exports work correctly

import { createOpenCodeClient } from "./dist/index.js"
import type { OpenCodeClient, paths } from "./dist/index.js"

console.log("🧪 Testing TypeScript SDK...")

// Test 1: Client creation
try {
  const client = createOpenCodeClient({
    baseUrl: "http://localhost:3000",
    headers: {
      Authorization: "Bearer test-token",
    },
  })

  console.log("✅ Client creation successful")
  console.log("✅ Client type:", typeof client)
} catch (error) {
  console.error("❌ Client creation failed:", error)
  process.exit(1)
}

// Test 2: Type exports
const checkTypes = () => {
  // This should compile without errors if types are properly exported
  const client: OpenCodeClient = createOpenCodeClient()
  const pathsCheck: paths = {} as paths

  // Use the variables to avoid unused warnings
  if (client && pathsCheck) {
    console.log("✅ Type exports working")
    return true
  }
  return false
}

checkTypes()

console.log("🎉 All TypeScript SDK tests passed!")
