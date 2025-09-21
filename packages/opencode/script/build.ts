#!/usr/bin/env bun
import { $ } from "bun"
import { SemVer } from "semver"

const dir = new URL("..", import.meta.url).pathname
process.chdir(dir)

import pkg from "../package.json"

async function checkGoVersion() {
  const goMod = await Bun.file("../tui/go.mod").text()
  const requiredVersion = goMod.match(/go (\d+\.\d+\.\d+)/)?.[1]

  if (!requiredVersion) {
    console.error("Could not find required Go version in go.mod")
    process.exit(1)
  }

  try {
    const { stdout } = await $`go version`.text()
    const installedVersion = stdout.match(/go version go(\d+\.\d+\.\d+)/)?.[1]

    if (!installedVersion) {
      console.error("Could not determine installed Go version.")
      process.exit(1)
    }

    if (new SemVer(installedVersion) < new SemVer(requiredVersion)) {
      console.error(`Go version ${requiredVersion} or higher is required.`)
      console.error(`You have ${installedVersion}.`)
      process.exit(1)
    }
  } catch (e) {
    console.error("Go is not installed or not in your PATH.")
    process.exit(1)
  }
}

await checkGoVersion()

const GOARCH: Record<string, string> = {
  arm64: "arm64",
  x64: "amd64",
  "x64-baseline": "amd64",
}

const targets = [
  ["windows", "x64"],
  ["linux", "arm64"],
  ["linux", "x64"],
  ["linux", "x64-baseline"],
  ["darwin", "x64"],
  ["darwin", "x64-baseline"],
  ["darwin", "arm64"],
]

await $`rm -rf dist`

const binaries: Record<string, string> = {}
const version = process.env["OPENCODE_VERSION"] ?? "dev"
for (const [os, arch] of targets) {
  console.log(`building ${os}-${arch}`)
  const name = `${pkg.name}-${os}-${arch}`
  await $`mkdir -p dist/${name}/bin`
  try {
    await $`CGO_ENABLED=0 GOOS=${os} GOARCH=${GOARCH[arch]} go build -ldflags="-s -w -X main.Version=${version}" -o ../opencode/dist/${name}/bin/tui ../tui/cmd/opencode/main.go`.cwd(
      "../tui",
    )
  } catch (e) {
    console.error("Failed to build Go TUI.")
    console.error(e)
    process.exit(1)
  }
  await Bun.build({
    compile: {
      target: `bun-${os}-${arch}` as any,
      outfile: `dist/${name}/bin/opencode`,
      execArgv: [`--user-agent=opencode/${version}`, `--env-file=""`, `--`],
      windows: {},
    },
    entrypoints: ["./src/index.ts"],
    define: {
      OPENCODE_VERSION: `'${version}'`,
      OPENCODE_TUI_PATH: `'../../../dist/${name}/bin/tui'`,
    },
  })
  await $`rm -rf ./dist/${name}/bin/tui`
  await Bun.file(`dist/${name}/package.json`).write(
    JSON.stringify(
      {
        name,
        version,
        os: [os === "windows" ? "win32" : os],
        cpu: [arch],
      },
      null,
      2,
    ),
  )
  binaries[name] = version
}

export { binaries }
