// Protocol message types for the two-agent plan-step-review loop

export interface Step {
  id: string
  title: string
  intent: string
  files: string[]
  exit: string
  risk: "low" | "med" | "high"
}

export interface Change {
  file: string
  kind: "edit" | "create" | "delete"
  reason: string
}

export interface Patch {
  replace: {
    id: string
    with: Step
  }
}

export interface Artifacts {
  patch: string // unified diff
  logs: string
}

export interface Metrics {
  compile_ms?: number
  tests?: {
    passed: number
    failed: number
  }
}

// Message types
export interface PlanPropose {
  type: "plan.propose"
  goal: string
  steps: Step[]
}

export interface PlanDecision {
  type: "plan.decision"
  status: "approve" | "revise" | "reject"
  patch?: Patch[]
  notes?: string
}

export interface StepIntent {
  type: "step.intent"
  id: string
  changes: Change[]
  tests: string[]
  checks: string[]
}

export interface StepGate {
  type: "step.gate"
  id: string
  decision: "approve" | "revise" | "hold"
  constraints?: string[]
  notes?: string
}

export interface StepResult {
  type: "step.result"
  id: string
  status: "pass" | "fail" | "partial"
  artifacts: Artifacts
  metrics: Metrics
}

// Union type for all messages
export type Message = PlanPropose | PlanDecision | StepIntent | StepGate | StepResult

// Type guards
export function isPlanPropose(msg: Message): msg is PlanPropose {
  return msg.type === "plan.propose"
}

export function isPlanDecision(msg: Message): msg is PlanDecision {
  return msg.type === "plan.decision"
}

export function isStepIntent(msg: Message): msg is StepIntent {
  return msg.type === "step.intent"
}

export function isStepGate(msg: Message): msg is StepGate {
  return msg.type === "step.gate"
}

export function isStepResult(msg: Message): msg is StepResult {
  return msg.type === "step.result"
}
