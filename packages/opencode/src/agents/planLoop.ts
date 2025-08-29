import { callModel } from "../llm/models"
import { applyStep } from "../run/applyStep"
import type { Message, PlanPropose, PlanDecision, StepIntent, StepGate, Step } from "../session/protocol/types"
import { Log } from "../util/log"

import EXECUTOR_PROMPT from "../session/prompt/executor-small.txt"
import SUPERVISOR_PROMPT from "../session/prompt/supervisor-gpt5.txt"

const log = Log.create({ service: "plan-loop" })

/**
 * Run the two-agent plan-step-review loop
 */
export async function runTask(goal: string): Promise<void> {
  log.info("Starting plan loop", { goal })

  // 1) Get initial plan from executor
  const planMessage = await callModel("executor_small", {
    system: EXECUTOR_PROMPT,
    user: `Goal: ${goal}`,
  })

  if (!isPlanPropose(planMessage)) {
    throw new Error("Expected plan.propose message from executor")
  }

  let currentPlan = planMessage

  // 2) Supervisor reviews plan
  let decision = await callModel("supervisor_gpt5", {
    system: SUPERVISOR_PROMPT,
    tool: currentPlan,
  })

  // Review loop: keep revising until approved or rejected
  while (isPlanDecision(decision) && decision.status === "revise") {
    log.info("Supervisor requested revision", { notes: decision.notes })

    // Feed revision back to executor
    const revisedPlan = await callModel("executor_small", {
      system: EXECUTOR_PROMPT,
      tool: decision,
    })

    if (!isPlanPropose(revisedPlan)) {
      throw new Error("Expected revised plan.propose message from executor")
    }

    currentPlan = revisedPlan
    decision = await callModel("supervisor_gpt5", {
      system: SUPERVISOR_PROMPT,
      tool: currentPlan,
    })
  }

  if (!isPlanDecision(decision) || decision.status !== "approve") {
    throw new Error("Plan was rejected by supervisor")
  }

  log.info("Plan approved", { steps: currentPlan.steps.length })

  // 3) Execute steps
  let stepsToExecute = currentPlan.steps

  // If supervisor provided a revised plan, use that
  if (decision.patch && decision.patch.length > 0) {
    const revisedStep = decision.patch[0].replace.with
    // Find the step to replace
    const stepIndex = stepsToExecute.findIndex((s) => s.id === decision.patch![0].replace.id)
    if (stepIndex >= 0) {
      stepsToExecute = [...stepsToExecute]
      stepsToExecute[stepIndex] = revisedStep
    }
  }

  for (const step of stepsToExecute) {
    await executeStep(step)
  }

  log.info("Plan execution completed")
}

/**
 * Execute a single step with intent/gate/apply/result flow
 */
async function executeStep(step: Step): Promise<void> {
  log.info("Executing step", { id: step.id, title: step.title })

  // Executor proposes step intent
  const intentMessage = await callModel("executor_small", {
    system: EXECUTOR_PROMPT,
    user: `Execute step: ${step.title} (${step.intent})`,
  })

  if (!isStepIntent(intentMessage)) {
    throw new Error("Expected step.intent message from executor")
  }

  // Supervisor gates the step
  const gateMessage = await callModel("supervisor_gpt5", {
    system: SUPERVISOR_PROMPT,
    tool: intentMessage,
  })

  if (!isStepGate(gateMessage)) {
    throw new Error("Expected step.gate message from supervisor")
  }

  if (gateMessage.decision !== "approve") {
    if (gateMessage.decision === "revise") {
      log.info("Step revised by supervisor", { id: step.id, notes: gateMessage.notes })
      // In a full implementation, we'd handle revision here
      throw new Error("Step revision not implemented")
    } else {
      throw new Error(`Step ${step.id} was ${gateMessage.decision} by supervisor`)
    }
  }

  // Apply the step (trusted runner)
  const result = await applyStep(intentMessage, gateMessage)

  log.info("Step result", { id: step.id, status: result.status })

  // Supervisor judges the result
  await callModel("supervisor_gpt5", {
    system: SUPERVISOR_PROMPT,
    tool: result,
  })

  if (result.status !== "pass") {
    log.warn("Step failed", { id: step.id, status: result.status })
    // In a full implementation, supervisor could issue remediation
    // For now, we'll stop on failure
    throw new Error(`Step ${step.id} failed: ${result.status}`)
  }
}

// Type guards
function isPlanPropose(msg: Message): msg is PlanPropose {
  return msg.type === "plan.propose"
}

function isPlanDecision(msg: Message): msg is PlanDecision {
  return msg.type === "plan.decision"
}

function isStepIntent(msg: Message): msg is StepIntent {
  return msg.type === "step.intent"
}

function isStepGate(msg: Message): msg is StepGate {
  return msg.type === "step.gate"
}
