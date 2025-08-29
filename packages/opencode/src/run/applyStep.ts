import { $ } from "bun"
import * as fs from "fs/promises"
import { minimatch } from "minimatch"
import type { StepIntent, StepGate, StepResult, Artifacts, Metrics } from "../session/protocol/types"
import { Log } from "../util/log"

const log = Log.create({ service: "apply-step" })

// Permission policies from the spec
const DENY_PATTERNS = ["**/*.env*", "**/*.key", "**/*.secret", "**/node_modules/**", "**/.git/**"]

const ALLOW_PATTERNS = ["**/*.go", "**/*.ts", "**/*.js", "**/*.json", "**/*.txt", "**/*.md", "**/*.yml", "**/*.yaml"]

/**
 * Check if a file path is allowed based on permission policies
 */
function isAllowedFile(filePath: string): boolean {
  // Check deny patterns first
  for (const pattern of DENY_PATTERNS) {
    if (minimatch(filePath, pattern)) {
      return false
    }
  }

  // Check allow patterns
  for (const pattern of ALLOW_PATTERNS) {
    if (minimatch(filePath, pattern)) {
      return true
    }
  }

  return false
}

/**
 * Apply a unified diff patch to the filesystem
 */
async function applyPatch(patchText: string): Promise<void> {
  // Parse the unified diff and apply changes
  const lines = patchText.split("\n")
  let currentFile = ""
  let inHunk = false
  let hunkLines: string[] = []
  let newContent = ""

  for (const line of lines) {
    if (line.startsWith("+++ ")) {
      // New file
      currentFile = line.substring(4)
      if (!isAllowedFile(currentFile)) {
        throw new Error(`File not allowed: ${currentFile}`)
      }
      inHunk = false
      hunkLines = []
      newContent = ""
    } else if (line.startsWith("@@ ")) {
      // Start of hunk
      inHunk = true
      hunkLines = []
    } else if (line.startsWith(" ")) {
      // Context line
      if (inHunk) {
        hunkLines.push(line.substring(1))
      }
    } else if (line.startsWith("-")) {
      // Removed line - skip in new content
      if (inHunk) {
        // This would be more complex in a real implementation
        // For now, we'll assume the patch is well-formed
      }
    } else if (line.startsWith("+")) {
      // Added line
      if (inHunk) {
        hunkLines.push(line.substring(1))
      }
    } else if (line === "") {
      // End of hunk
      if (inHunk && currentFile) {
        newContent = hunkLines.join("\n")
        await fs.writeFile(currentFile, newContent, "utf-8")
        inHunk = false
      }
    }
  }
}

/**
 * Run tests and checks
 */
async function runChecks(tests: string[], checks: string[]): Promise<{ logs: string; metrics: Metrics }> {
  const logs: string[] = []
  const metrics: Metrics = {}

  const startTime = Date.now()

  // Run tests
  for (const test of tests) {
    try {
      log.info("Running test", { test })
      const result = await $`${test}`.nothrow()
      logs.push(`Test: ${test}\n${result.stdout}\n${result.stderr}`)
      if (result.exitCode !== 0) {
        throw new Error(`Test failed: ${test}`)
      }
    } catch (error) {
      logs.push(`Test failed: ${test}\n${error}`)
      throw error
    }
  }

  // Run checks
  for (const check of checks) {
    try {
      log.info("Running check", { check })
      const result = await $`${check}`.nothrow()
      logs.push(`Check: ${check}\n${result.stdout}\n${result.stderr}`)
      if (result.exitCode !== 0) {
        throw new Error(`Check failed: ${check}`)
      }
    } catch (error) {
      logs.push(`Check failed: ${check}\n${error}`)
      throw error
    }
  }

  metrics.compile_ms = Date.now() - startTime

  // Mock test results (in real implementation, parse actual test output)
  metrics.tests = {
    passed: tests.length,
    failed: 0,
  }

  return { logs: logs.join("\n\n"), metrics }
}

/**
 * Apply a step with policy enforcement and return results
 */
export async function applyStep(intent: StepIntent, gate: StepGate): Promise<StepResult> {
  const { id, changes, tests, checks } = intent
  const { constraints } = gate

  log.info("Applying step", { id, changes: changes.length, tests: tests.length, checks: checks.length })

  try {
    // Enforce constraints
    if (constraints) {
      for (const constraint of constraints) {
        log.info("Enforcing constraint", { constraint })
        // Implement constraint checking logic here
        // For now, just log them
      }
    }

    // Validate all files are allowed
    for (const change of changes) {
      if (!isAllowedFile(change.file)) {
        throw new Error(`File not allowed by policy: ${change.file}`)
      }
    }

    // Apply changes (in real implementation, this would apply the actual patch)
    // For now, we'll simulate applying changes
    const mockPatch = `+++ ${changes[0]?.file || "example.txt"}
@@ -1,1 +1,1 @@
- old content
+ new content`

    await applyPatch(mockPatch)

    // Run tests and checks
    const { logs, metrics } = await runChecks(tests, checks)

    const artifacts: Artifacts = {
      patch: mockPatch,
      logs,
    }

    return {
      type: "step.result",
      id,
      status: "pass",
      artifacts,
      metrics,
    }
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : String(error)
    log.error("Step failed", { id, error: errorMessage })

    return {
      type: "step.result",
      id,
      status: "fail",
      artifacts: {
        patch: "",
        logs: `Error: ${errorMessage}`,
      },
      metrics: {},
    }
  }
}
