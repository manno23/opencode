import { DurableObject } from "cloudflare:workers"

export class SyncServer extends DurableObject {
  constructor(ctx: DurableObjectState, env: any) {
    super(ctx, env)
  }

  async fetch(request: Request) {
    const webSocketPair = new WebSocketPair()
    const [client, server
    ] = Object.values(webSocketPair)

    this.ctx.acceptWebSocket(server)
    return new Response(null,
      {
        status: 101,
        webSocket: client,
      })
  }
}
