import { getRequestEvent } from "solid-js/web"
import { and, Database, eq, inArray } from "@opencode/cloud-core/drizzle/index.js"
import { WorkspaceTable } from "@opencode/cloud-core/schema/workspace.sql.js"
import { UserTable } from "@opencode/cloud-core/schema/user.sql.js"
import { redirect } from "@solidjs/router"
import { AccountTable } from "@opencode/cloud-core/schema/account.sql.js"
import { Actor } from "@opencode/cloud-core/actor.js"

import { createClient } from "@openauthjs/openauth/client"
import { useAuthSession } from "./auth.session"

export const AuthClient = createClient({
  clientID: "app",
  issuer: import.meta.env.VITE_AUTH_URL,
})

export const getActor = async (): Promise<Actor.Info> => {
  "use server"
  const evt = getRequestEvent()
  if (!evt) throw new Error("No request event")
  const url = new URL(evt.request.headers.has("x-server-id") ? evt.request.headers.get("referer")! : evt.request.url)
  const auth = await useAuthSession()
  auth.data.account = auth.data.account ?? {}
  const splits = url.pathname.split("/").filter(Boolean)
  if (splits[0] !== "workspace") {
    const current = auth.data.account[auth.data.current ?? ""]
    if (current) {
      return {
        type: "account",
        properties: {
          email: current.email,
          accountID: current.id,
        },
      }
    }
    if (Object.keys(auth.data.account ?? {}).length > 0) {
      const current = Object.values(auth.data.account)[0]
      await auth.update((val) => ({
        ...val,
        current: current.id,
      }))
      return {
        type: "account",
        properties: {
          email: current.email,
          accountID: current.id,
        },
      }
    }
    return {
      type: "public",
      properties: {},
    }
  }
  const workspaceHint = splits[1]
  const accounts = Object.keys(auth.data.account ?? {})
  if (accounts.length) {
    const result = await Database.transaction(async (tx) => {
      return await tx
        .select({
          user: UserTable,
        })
        .from(AccountTable)
        .innerJoin(UserTable, and(eq(UserTable.email, AccountTable.email)))
        .innerJoin(WorkspaceTable, eq(WorkspaceTable.id, UserTable.workspaceID))
        .where(and(inArray(AccountTable.id, accounts), eq(WorkspaceTable.id, workspaceHint)))
        .limit(1)
        .execute()
        .then((x) => x[0])
    })
    if (result) {
      return {
        type: "user",
        properties: {
          userID: result.user.id,
          workspaceID: result.user.workspaceID,
        },
      }
    }
  }
  throw redirect("/auth/authorize")
}
