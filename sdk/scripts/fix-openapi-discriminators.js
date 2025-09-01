#!/usr/bin/env node

/**
 * Fix OpenAPI discriminator mappings by replacing inline objects with $ref references
 * This resolves oapi-codegen compatibility issues with discriminator mappings.
 */

import fs from "fs"
import path from "path"

function fixDiscriminatorMappings(spec) {
  const fixes = []

  // Recursive function to find and fix discriminator mappings
  function fixDiscriminatorsRecursive(obj, path = "") {
    if (!obj || typeof obj !== "object") return

    // Check if this is a discriminator with oneOf that needs fixing
    if (obj.discriminator && obj.oneOf && Array.isArray(obj.oneOf)) {
      const mapping = obj.discriminator.mapping || {}

      // Check if any oneOf items are inline objects that should be $refs
      const newOneOf = obj.oneOf.map((item) => {
        // Skip if already a $ref
        if (item["$ref"]) return item

        // Check if this inline object matches a discriminator mapping
        if (item.properties && obj.discriminator.propertyName) {
          const discriminatorProp = item.properties[obj.discriminator.propertyName]
          if (discriminatorProp && discriminatorProp.const) {
            const discriminatorValue = discriminatorProp.const
            const mappedRef = mapping[discriminatorValue]

            if (mappedRef) {
              fixes.push(`Fixed discriminator at ${path}: replaced inline ${discriminatorValue} with ${mappedRef}`)
              return { $ref: mappedRef }
            }
          }
        }

        return item
      })

      obj.oneOf = newOneOf
    }

    // Recursively process all properties
    if (Array.isArray(obj)) {
      obj.forEach((item, index) => {
        fixDiscriminatorsRecursive(item, `${path}[${index}]`)
      })
    } else {
      Object.keys(obj).forEach((key) => {
        const newPath = path ? `${path}.${key}` : key
        fixDiscriminatorsRecursive(obj[key], newPath)
      })
    }
  }

  // Start recursive fixing from the root
  fixDiscriminatorsRecursive(spec, "")

  return { spec, fixes }
}

function main() {
  const inputFile = process.argv[2] || "schema/openapi.fixed.json"
  const outputFile = process.argv[3] || "schema/openapi.oapi-compatible.json"

  console.log("🔧 Fixing OpenAPI discriminator mappings for oapi-codegen compatibility...")
  console.log(`Input: ${inputFile}`)
  console.log(`Output: ${outputFile}`)

  if (!fs.existsSync(inputFile)) {
    console.error(`❌ Input file not found: ${inputFile}`)
    process.exit(1)
  }

  try {
    // Read and parse the spec
    const specContent = fs.readFileSync(inputFile, "utf8")
    const spec = JSON.parse(specContent)

    // Apply fixes
    const { spec: fixedSpec, fixes } = fixDiscriminatorMappings(spec)

    // Write the fixed spec
    fs.writeFileSync(outputFile, JSON.stringify(fixedSpec, null, 2))

    console.log("✅ Fixed OpenAPI specification saved")
    console.log("\nApplied fixes:")
    fixes.forEach((fix) => console.log(`  - ${fix}`))

    console.log("\nNext steps:")
    console.log(`  1. Test: oapi-codegen -generate types -package opencode ${outputFile}`)
    console.log("  2. If successful, use this spec for SDK generation")
  } catch (error) {
    console.error("❌ Error processing specification:", error.message)
    process.exit(1)
  }
}

if (import.meta.url === `file://${process.argv[1]}`) {
  main()
}
