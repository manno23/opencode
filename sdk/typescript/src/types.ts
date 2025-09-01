import type { components, paths } from "./types.gen"

export type Event = components["schemas"]["Event"]
export type Part = components["schemas"]["Part"]
export type UserMessage = components["schemas"]["UserMessage"]

// TODO: replace unknown once the spec defines these exact names
export type App = components["schemas"]["App"] | unknown
export type Model = components["schemas"]["Model"] | unknown
export type Provider = components["schemas"]["Provider"] | unknown
export type Permission = components["schemas"]["Permission"] | unknown
export type Auth = components["schemas"]["Auth"] | unknown
export type Config = components["schemas"]["Config"] | unknown

// Optional helpers
export type GetAppResponse = paths["/app"]["get"]["responses"]["200"] extends {
  content: { "application/json": infer T }
}
  ? T
  : unknown
export type InitAppBody = paths["/app/init"]["post"]["requestBody"] extends {
  content: { "application/json": infer T }
}
  ? T
  : unknown
