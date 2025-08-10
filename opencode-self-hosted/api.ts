import { Hono } from "hono"
import { DurableObject } from "cloudflare:workers"
// @ts-ignore - crypto available in Cloudflare Workers
const randomUUID = () => crypto.randomUUID()

type Env = {
  SYNC_SERVER: DurableObjectNamespace<SyncServer>
  Bucket: R2Bucket
  WEB_DOMAIN: string
}

export class SyncServer extends DurableObject<Env> {
  constructor(ctx: DurableObjectState, env: Env) {
    super(ctx, env)
  }

  async fetch() {
    console.log("SyncServer subscribe")

    const webSocketPair = new WebSocketPair()
    const [client, server] = Object.values(webSocketPair)

    this.ctx.acceptWebSocket(server)

    const data = await this.ctx.storage.list()
    Array.from(data.entries())
      .filter(([key, _]) => key.startsWith("session/"))
      .map(([key, content]) => server.send(JSON.stringify({ key, content })))

    return new Response(null, {
      status: 101,
      webSocket: client,
    })
  }

  async webSocketMessage(_ws: WebSocket, _message: string | ArrayBuffer) { }

  async webSocketClose(ws: WebSocket, code: number, _reason: string, _wasClean: boolean) {
    ws.close(code, "Durable Object is closing WebSocket")
  }

  async publish(key: string, content: any) {
    const sessionID = await this.getSessionID()
    if (
      !key.startsWith(`session/info/${sessionID}`) &&
      !key.startsWith(`session/message/${sessionID}/`) &&
      !key.startsWith(`session/part/${sessionID}/`)
    )
      return new Response("Error: Invalid key", { status: 400 })

    await this.env.Bucket.put(`share/${key}.json`, JSON.stringify(content), {
      httpMetadata: { contentType: "application/json" },
    })
    await this.ctx.storage.put(key, content)

    const clients = this.ctx.getWebSockets()
    for (const client of clients) {
      client.send(JSON.stringify({ key, content }))
    }
  }

  async share(sessionID: string) {
    let secret = await this.getSecret()
    if (secret) return secret
    secret = randomUUID()

    await this.ctx.storage.put("secret", secret)
    await this.ctx.storage.put("sessionID", sessionID)
    return secret
  }

  async getData() {
    const data = (await this.ctx.storage.list()) as Map<string, any>
    return Array.from(data.entries())
      .filter(([key, _]) => key.startsWith("session/"))
      .map(([key, content]) => ({ key, content }))
  }

  async assertSecret(secret: string) {
    if (secret !== (await this.getSecret())) throw new Error("Invalid secret")
  }

  private async getSecret() {
    return this.ctx.storage.get<string>("secret")
  }

  private async getSessionID() {
    return this.ctx.storage.get<string>("sessionID")
  }

  async clear() {
    const sessionID = await this.getSessionID()
    const list = await this.env.Bucket.list({
      prefix: `session/message/${sessionID}/`,
      limit: 1000,
    })
    for (const item of list.objects) {
      await this.env.Bucket.delete(item.key)
    }
    await this.env.Bucket.delete(`session/info/${sessionID}`)
    await this.ctx.storage.deleteAll()
  }

  static shortName(id: string) {
    return id.substring(id.length - 8)
  }
}

export default new Hono<{ Bindings: Env }>()
  .get("/", (c) => c.text("Hello, world!"))
  .post("/share_create", async (c) => {
    const body = await c.req.json<{ sessionID: string }>()
    const sessionID = body.sessionID
    const short = SyncServer.shortName(sessionID)
    const id = c.env.SYNC_SERVER.idFromName(short)
    const stub = c.env.SYNC_SERVER.get(id)
    const secret = await stub.share(sessionID)
    return c.json({
      secret,
      url: `https://${c.env.WEB_DOMAIN}/s/${short}`,
    })
  })
  .post("/share_delete", async (c) => {
    const body = await c.req.json<{ sessionID: string; secret: string }>()
    const sessionID = body.sessionID
    const secret = body.secret
    const id = c.env.SYNC_SERVER.idFromName(SyncServer.shortName(sessionID))
    const stub = c.env.SYNC_SERVER.get(id)
    await stub.assertSecret(secret)
    await stub.clear()
    return c.json({})
  })
  .post("/share_sync", async (c) => {
    const body = await c.req.json<{
      sessionID: string
      secret: string
      key: string
      content: any
    }>()
    const name = SyncServer.shortName(body.sessionID)
    const id = c.env.SYNC_SERVER.idFromName(name)
    const stub = c.env.SYNC_SERVER.get(id)
    await stub.assertSecret(body.secret)
    await stub.publish(body.key, body.content)
    return c.json({})
  })
  .get("/share_poll", async (c) => {
    const upgradeHeader = c.req.header("Upgrade")
    if (!upgradeHeader || upgradeHeader !== "websocket") {
      return c.text("Error: Upgrade header is required", { status: 426 })
    }
    const id = c.req.query("id")
    if (!id) return c.text("Error: Share ID is required", { status: 400 })
    const stub = c.env.SYNC_SERVER.get(c.env.SYNC_SERVER.idFromName(id))
    return stub.fetch(c.req.raw)
  })
  .get("/share_data", async (c) => {
    const id = c.req.query("id")
    if (!id) return c.text("Error: Share ID is required", { status: 400 })
    const stub = c.env.SYNC_SERVER.get(c.env.SYNC_SERVER.idFromName(id))
    const data = await stub.getData()

    let info
    const messages: Record<string, any> = {}
    data.forEach((d: any) => {
      const [root, type] = d.key.split("/")
      if (root !== "session") return
      if (type === "info") {
        info = d.content
        return
      }
      if (type === "message") {
        messages[d.content.id] = {
          parts: [],
          ...d.content,
        }
      }
      if (type === "part") {
        messages[d.content.messageID].parts.push(d.content)
      }
    })

    return c.json({ info, messages })
  })
  .all("*", (c) => c.text("Not Found"))

