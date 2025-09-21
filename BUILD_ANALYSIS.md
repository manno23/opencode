# Build and Maintenance Infrastructure Analysis

This document provides an analysis of the build and maintenance infrastructure for the `opencode` project, focusing on the main `bun` application.

## 1. Monorepo Structure

The project is a large monorepo managed with `bun` as the package manager and `turbo` for task execution. It is composed of several packages, each with a specific role:

-   `packages/opencode`: The core of the project, a TypeScript application that forms the `opencode` command-line interface (CLI).
-   `packages/tui`: A Go-based project responsible for the terminal user interface (TUI).
-   `packages/app`: A web application, likely a dashboard or a companion to the CLI.
-   `packages/console`: A web-based console for the application.
-   `packages/sdk`: Contains software development kits (SDKs) for interacting with the `opencode` platform.
-   `sst.config.ts`: Indicates the use of the Serverless Stack (SST) for infrastructure deployment.

## 2. Build Process

The primary build process is for the `opencode` CLI, which is orchestrated by the `packages/opencode/script/build.ts` script. This script is responsible for creating distributable, cross-platform binaries for Windows, Linux, and macOS.

The build process is a hybrid of Go and Bun, and can be broken down into the following steps for each target platform:

1.  **Compile the TUI:** The Go TUI from `packages/tui/cmd/opencode/main.go` is compiled into a native binary.
2.  **Compile the CLI:** The main TypeScript application from `packages/opencode/src/index.ts` is compiled into a single, native executable using `Bun.build` with the `compile` option.
3.  **Embed the TUI:** The path to the compiled TUI binary is embedded as a variable within the `opencode` executable. This means the Go binary is effectively bundled with the main application.
4.  **Package the Binary:** The final executable is placed in a `dist` directory, along with a `package.json` file, making it ready for distribution.

## 3. Key Technologies

-   **Bun:** Serves as the package manager, runtime environment, and bundler. The `Bun.build` API is critical to the build process.
-   **Go:** Used to build the TUI component. The project requires Go version 1.24.0, as specified in `packages/tui/go.mod`.
-   **Turbo:** Used for running tasks across the monorepo, such as `typecheck`.
-   **SST:** Used for deploying the project's infrastructure.

## 4. Potential Problem Points

The current build and maintenance infrastructure has several potential points of failure that could make it difficult to build and maintain the application:

1.  **Hard Dependency on Go Version:** The build script implicitly depends on Go version 1.24.0, but it does not verify that the correct version is installed. If a developer has a different version of Go, the build may fail with cryptic error messages.
2.  **Cascading Build Failures:** The `opencode` build is tightly coupled to the `packages/tui` build. Any error in the Go project (e.g., dependency issues, compilation errors) will cause the entire build process to fail.
3.  **Fragile `Bun.build` Usage:** The build process relies on the `Bun.build` `compile` API, which is a powerful but complex feature. This could be a source of issues due to bugs in Bun itself or incorrect usage of the API.
4.  **Cross-Compilation Complexity:** Managing cross-compilation for both Go and Bun across multiple platforms can be challenging. The build script relies on specific toolchain behaviors and environment variables that may not be consistent across all development environments.
5.  **Lack of a Centralized Build Script:** There is no single, top-level script to build the entire project. Developers need to know to run the `build.ts` script within the `packages/opencode` directory, which is not immediately obvious. The `dev` script in the root `package.json` only runs the application in development mode. This can be a barrier for new contributors.
