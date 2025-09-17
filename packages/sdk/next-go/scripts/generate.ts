#!/usr/bin/env bun
import { $ } from "bun"
import path from "path"
const dir = new URL(".", import.meta.url).pathname
process.chdir(dir)

console.log("=== Generating Go SDK ===")
console.log(process.cwd())

// Updater logic to modify openapi.json
import { readFileSync, writeFileSync, } from "fs"

const openapiPath = path.join(dir, "openapi.json")


await $`bun run --conditions=development ../../opencode/src/index.ts generate > openapi.json`

await $`echo "Updating the OpenAPI spec to meet criteria of the generator..."`
try {
  // Read and parse the JSON
  const data: any = JSON.parse(readFileSync(openapiPath, "utf8"))

  // Function to recursively find and update "exclusiveMinimum"
  function updateExclusiveMinimum(obj: any): void {
    if (typeof obj !== "object" || obj === null) return

    for (const key in obj) {
      if (key === "exclusiveMinimum") {
        // Insert "minimum": 0 above "exclusiveMinimum"
        const newObj: any = {}
        // Copy other properties
        for (const k in obj) {
          if (k !== "exclusiveMinimum") {
            newObj[k] = obj[k]
          }
        }

        newObj.type = "number"
        newObj.minimum = 0
        newObj.exclusiveMinimum = true // Ensure it's true

        // Replace the object
        Object.assign(obj, newObj)
        console.log("Updated exclusiveMinimum in object:", JSON.stringify(obj, null, 2))
      } else if (typeof obj[key] === "object") {
        updateExclusiveMinimum(obj[key])
      }
    }
  }

  // Update the data
  updateExclusiveMinimum(data)

  // Write back the updated JSON
  writeFileSync(openapiPath, JSON.stringify(data, null, 2))
  console.log("openapi.json updated successfully.")
} catch (error) {
  console.error("Error updating openapi.json:", (error as Error).message)
  process.exit(1)
}

await $`echo "Regenerating the Go client side SDK. This only runs if there has been a change."`
