
---

### Documentation

For more info on how to configure opencode [**head over to our docs**](https://opencode.ai/docs).

To run opencode locally you need.

- Bun
- Golang 1.24.x

And run.

```bash
$ bun install
$ bun dev

```


#### Development Notes

This fork is making use of the client / server architecture here, to edit code and arrange agents where the providers are running them.
Cloudflare specifically.

~~**API Client**: After making changes to the TypeScript API endpoints in `packages/opencode/src/server/server.ts`, you will need the opencode team to generate a new stainless sdk for the clients.~~
API Client generation is now done in tree, so the Go TUI application can be modified however you like and you dont have to ask noone for nutt'in.


### FAQ

#### How is this different than Opencode?

It's very similar to Opencode in terms of capability. Here are the key differences:

- 100% open source
- Not coupled to any provider. Although Anthropic is recommended, opencode can be used with OpenAI, Google or even local models. As models evolve the gaps between them will close and pricing will drop so being provider-agnostic is important.
- A focus on TUI. opencode is built by neovim users and the creators of [terminal.shop](https://terminal.shop); we are going to push the limits of what's possible in the terminal.
- A client/server architecture. This for example can allow opencode to run on your computer, while you can drive it remotely from a mobile app. Meaning that the TUI frontend is just one of the possible clients.

#### What's the other repo?

The other confusingly named repo has no relation to this one. You can [read the story behind it here](https://x.com/thdxr/status/1933561254481666466).

---
