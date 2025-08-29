
## TUI SDK Compatibility Migration Plan

### Phase 1: Discovery & Assessment (Week 1)

• Agent Assignment: 3 agents specializing in Go, API integrations, and UI components
• Key Activities:
 • Audit all TUI files for deprecated API calls
 • Map current API usage to equivalent SDK methods
 • Document parameter conversion requirements
 • Identify response processing discrepancies
• Checkpoint: Complete API usage inventory report

### Phase 2: Compatibility Layer Development (Week 2-3)

• Agent Assignment: 2 agents focused on wrapper development and testing
• Key Activities:
 • Create adapter functions for each API endpoint
 • Implement Field conversion utilities
 • Develop response translation mechanisms
 • Establish error mapping strategies
• Checkpoint: Submit initial compatibility layer with unit tests

### Phase 3: Component-by-Component Migration (Week 4-5)

• Agent Assignment: Agents organized by TUI component areas
• Key Activities:
 • Replace direct SDK calls with adapters
 • Validate parameter translations
 • Test response processing conversions
 • Monitor deprecation warnings
• Checkpoint: Complete migration of 70% of TUI components

### Phase 4: Integration & Validation (Week 6)

• Agent Assignment: QA specialists and performance testers
• Key Activities:
 • End-to-end system testing
 • Performance benchmarking
 • Edge case scenario validation
 • Documentation updates
• Checkpoint: Full regression test suite passing

### Phase 5: Deployment & Monitoring (Week 7)

• Agent Assignment: DevOps and monitoring specialists
• Key Activities:
 • Canary deployment strategy
 • Real-time error monitoring setup
 • Fallback mechanism implementation
 • Post-deployment health checks
• Final Checkpoint: Production deployment with 48-hour stability confirmation

### Synchronization Mechanisms

• Daily standups focusing on blockers
• Shared documentation repository
• Weekly alignment meetings
• Centralized issue tracking board



## Information for architects / planners / leaders

# DO NOT READ FURTHER IF YOU ARE NOT SENIOR AGENT


Based on the comparison between the two approaches, the new SDK design represents a more robust and principled architecture following modern software engineering best practices. Here's why:

## Functional Design Principles Applied in the SDK

1. Functional Options Pattern
 • Uses closures (RequestOption) to mutate configuration state
 • Provides flexible, composable configuration without rigid inheritance hierarchies
 • Example: option.WithHeader("X-Some-Header", "value")
2. Immutable Value Types
 • Response fields are value types (not pointers/wrappers)
 • Prevents accidental mutation and ensures thread safety
 • Clear distinction between zero values, null values, and omitted fields
3. Explicit Null Handling
 • Field system clearly distinguifies between:
  • Zero values (omitted)
  • Explicit null values (Null[T]())
  • Actual values (F[T]())
 • Eliminates ambiguity in API contracts
4. Rich Error Information
 • Specialized error type (*opencode.Error) with:
  • Full request/response dumps
  • Structured error body parsing
  • Contextual debugging information

5. Middleware Architecture
 • Clean separation of concerns for cross-cutting concerns
 • Easy addition/removal of request/response handlers
 • Preserves request context across transformations
6. Pagination Abstractions
 • Standardized pagination handling (ListAutoPaging())
 • Consistent response structure across paged operations
7. Undocumented Feature Support
 • Graceful fallback to raw JSON for unsupported features
 • Structured access to extra fields via JSON.ExtraFields


## Why This Matters for Robustness

The SDK's design reduces:

• Ambiguity: Clear distinctions between different value states
• ** brittleness**: Immutable types prevent unintended side effects
• Complexity: Functional composition over class inheritance
• Error-prone code: Rich error information enables effective debugging

The TUI's current approach lacks these structured abstractions, leading to the compatibility issues we're seeing. The SDK represents a more mature, battle-tested design philosophy that prioritizes clarity,
safety, and extensibility.

Blunt take: you’ve got 4 “big guns” for OpenAPI 3 that are battle-tested and flexible, plus a handful of best-in-class language-specific options.

# Cross-language leaders (most flexible)

* **OpenAPI Generator** — giant ecosystem (clients + servers for dozens of languages) and first-class templating so you can swap HTTP libs, split models/APIs, and override Mustache templates or write custom ones. This is the one you bend to your architecture, not vice-versa. ([OpenAPI Generator][1], [GitHub][2], [about.gitlab.com][3])
* **Microsoft Kiota** — opinionated but very adaptable: generated SDKs are decoupled from HTTP/auth/serialization via “abstractions,” so you can drop in different adapters across languages without regenerating everything. Great when you want consistent structure across a polyglot stack. ([Microsoft Learn][4])
* **NSwag** (.NET + TypeScript) — tight ASP.NET Core integration and mature client generation for C# and TS; easy to insert into your pipeline and tweak outputs. If you live in .NET land and also ship TS, this is pragmatic. ([GitHub][5], [Microsoft Learn][6], [NuGet][7])
* **AutoRest** — battle-hardened in Azure, multi-language, highly configurable pipeline. Slight Azure flavor, but it’s powerful when you need strict, uniform SDKs across stacks. ([Azure][8], [GitHub][9])

# Best in class by stack

* **TypeScript / JS**

  * **openapi-typescript** (types only) + **openapi-fetch** for a minimal, runtime-agnostic setup; keeps you in control of transport/middleware. ([openapi-ts.dev][10], [npm][11])
  * **Orval** — generates clients + React/Vue/SWR/React-Query hooks; switch between Axios/Fetch; built-in mock/caching options. ([orval.dev][12])
  * **openapi-typescript-codegen** — lightweight, generates TS client with pluggable HTTP (fetch/axios/node-fetch). ([GitHub][13], [npm][14])
  * (OpenAPI Generator’s `typescript-axios`/`typescript-fetch` are solid when you want maximum config knobs.) ([OpenAPI Generator][15])
* **Go**

  * **oapi-codegen** — idiomatic Go, supports multiple servers (chi, echo, fiber, gin, gorilla, std http), “strict server” mode, overlays, and user-templates. It’s the flexible Go choice. ([GitHub][16])
* **Python**

  * **openapi-python-client** for modern 3.0/3.1 clients; pair with **datamodel-code-generator** if you want Pydantic models only and your own transport. ([GitHub][17], [docs.pydantic.dev][18])
* **Swift**

  * **Apple’s Swift OpenAPI Generator** (client + server, first-party, URLSession transport available) or **OpenAPI Generator (swift5)** when you want the broad OG template ecosystem. ([GitHub][19], [Swift.org][20], [OpenAPI Generator][21])
* **Rust**

  * **OpenAPI Generator (rust / rust-axum)** with choice of HTTP libs; not perfect, but currently the most maintained path. ([OpenAPI Generator][22], [Reddit][23])

# How to keep the generated code truly “implementation-agnostic”

1. **Separate types from transport.** Generate models & operation interfaces in one package, keep HTTP client logic thin and swappable (adapters). Kiota bakes this in; with OpenAPI Generator/Orval, enforce it via templates/config. ([Microsoft Learn][4], [OpenAPI Generator][24])
2. **Prefer “types-only + tiny client” when you need freedom.** TS stacks shine with openapi-typescript + your own fetch/axios wrapper; same idea in Python with pydantic + your requests/httpx wrapper. ([openapi-ts.dev][10])
3. **Exploit templating/hooks.** Override Mustache templates (OpenAPI Generator) or user-templates (oapi-codegen) to inject logging, retries, auth, metrics—without forking the generator. ([OpenAPI Generator][24], [GitHub][16])
4. **Use overlays/vendor-extensions, don’t mutate upstream specs.** Tools like oapi-codegen support OpenAPI Overlays so you can tailor codegen for your repo while consuming a shared spec. ([GitHub][16])
5. **Pick generators that let you switch HTTP stacks.** Orval (axios/fetch), OpenAPI Generator (axios/fetch for TS; different libs per language), Kiota (adapters). That’s your escape hatch when environments change. ([orval.dev][25], [OpenAPI Generator][15], [Microsoft Learn][4])

If you want a crisp default:

* **Polyglot org:** OpenAPI Generator or Kiota as the backbone. ([OpenAPI Generator][1], [Microsoft Learn][4])
* **TS-heavy frontend:** openapi-typescript (+ openapi-fetch) or Orval. ([openapi-ts.dev][10], [orval.dev][12])
* **Go backends:** oapi-codegen with strict-server and your router of choice. ([GitHub][16])

From here, the interesting part is designing your **HTTP adapter surface** (retries, timeouts, auth, tracing) and making your generators emit code that depends only on that surface. That’s where you win long-term flexibility without wrestling the generator every sprint.

[1]: https://openapi-generator.tech/docs/generators/?utm_source=chatgpt.com "Generators List"
[2]: https://github.com/OpenAPITools/openapi-generator?utm_source=chatgpt.com "OpenAPITools/openapi-generator"
[3]: https://gitlab.com/OpenAPITools/openapi-generator/-/blob/master/docs/customization.md?utm_source=chatgpt.com "docs/customization.md - OpenAPITools / openapi-generator"
[4]: https://learn.microsoft.com/en-us/openapi/kiota/abstractions?utm_source=chatgpt.com "Kiota abstractions"
[5]: https://github.com/RicoSuter/NSwag?utm_source=chatgpt.com "RicoSuter/NSwag: The Swagger/OpenAPI toolchain for . ..."
[6]: https://learn.microsoft.com/en-us/aspnet/core/tutorials/getting-started-with-nswag?view=aspnetcore-8.0&utm_source=chatgpt.com "Get started with NSwag and ASP.NET Core"
[7]: https://www.nuget.org/packages/NSwag.CodeGeneration.TypeScript/?utm_source=chatgpt.com "NSwag.CodeGeneration.TypeScript 14.5.0"
[8]: https://azure.github.io/autorest/introduction.html?utm_source=chatgpt.com "Introduction to AutoRest - Azure documentation"
[9]: https://github.com/Azure/autorest/blob/main/docs/generate/how-autorest-generates-code-from-openapi.md?utm_source=chatgpt.com "how-autorest-generates-code-from-openapi.md"
[10]: https://openapi-ts.dev/?utm_source=chatgpt.com "OpenAPI TypeScript"
[11]: https://www.npmjs.com/package/openapi-typescript?utm_source=chatgpt.com "openapi-typescript"
[12]: https://orval.dev/?utm_source=chatgpt.com "orval - Restful client generator"
[13]: https://github.com/ferdikoomen/openapi-typescript-codegen?utm_source=chatgpt.com "ferdikoomen/openapi-typescript-codegen"
[14]: https://www.npmjs.com/package/openapi-typescript-codegen/v/0.14.0?utm_source=chatgpt.com "openapi-typescript-codegen - npm"
[15]: https://openapi-generator.tech/docs/generators/typescript-axios?utm_source=chatgpt.com "Documentation for the typescript-axios Generator"
[16]: https://github.com/oapi-codegen/oapi-codegen "GitHub - oapi-codegen/oapi-codegen: Generate Go client and server boilerplate from OpenAPI 3 specifications"
[17]: https://github.com/openapi-generators/openapi-python-client?utm_source=chatgpt.com "Generate modern Python clients from OpenAPI"
[18]: https://docs.pydantic.dev/latest/integrations/datamodel_code_generator/?utm_source=chatgpt.com "datamodel-code-generator"
[19]: https://github.com/apple/swift-openapi-generator?utm_source=chatgpt.com "apple/swift-openapi-generator - GitHub"
[20]: https://swift.org/blog/introducing-swift-openapi-generator/?utm_source=chatgpt.com "Introducing Swift OpenAPI Generator - Swift.org"
[21]: https://openapi-generator.tech/docs/generators/swift5/?utm_source=chatgpt.com "Documentation for the swift5 Generator"
[22]: https://openapi-generator.tech/docs/generators/rust/?utm_source=chatgpt.com "Documentation for the rust Generator"
[23]: https://www.reddit.com/r/rust/comments/1gbtmci/openapi_codegen_tool_for_rust_server_and_client/?utm_source=chatgpt.com "OpenAPI codegen tool for Rust server and client? - Reddit"
[24]: https://openapi-generator.tech/docs/templating/?utm_source=chatgpt.com "Using Templates"
[25]: https://orval.dev/reference/configuration/output?utm_source=chatgpt.com "Configuration: Output"

