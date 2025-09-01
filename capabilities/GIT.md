Here’s a practical, low-drama way to manage your heavy-divergence fork while keeping upstream pullable.

Goals

- Keep a clean “upstream-tracking” line for easy rebases/merges.
- Isolate your invasive changes (subproject swaps, build overhaul, removing SST/Stainless).
- Preserve a PR-friendly history to your own GitHub upstream (manno23/dev).
- Avoid future merge hell as you add agents/sandboxing.

Recommended branch model

- upstream/main (remote): Original project’s default branch.
- origin/main (remote): Your fork’s default branch (kept close to upstream).
- origin/manno23/dev (remote): Your published integration branch.
- local branches:
  - track-upstream: Mirror of upstream/main. No changes.
  - dev: Your working branch with all changes.
  - topic/\*: Focused feature branches merged into dev via PRs.

Setup (one-time)

- Add remotes:
  - upstream -> original project
  - origin -> your fork
- Create tracker:
  - track-upstream: purely follows upstream/main
  - dev: based off track-upstream (initially)

Update flow (ongoing)

1. Sync upstream and rebase dev regularly

```bash
git fetch upstream
git checkout track-upstream
git reset --hard upstream/main
git checkout dev
git rebase track-upstream   # or: git merge --no-ff track-upstream if you prefer merges
```

2. Publish your work

- git push -u origin dev:manno23/dev

3. Keep origin/main “clean”

- Periodically fast-forward origin/main from track-upstream (so your fork’s main stays close to upstream):
  - git checkout track-upstream
  - git push origin track-upstream:main

Structural strategy for large changes

- Early, apply deletions/moves as separate commits with clear messages:
  - Remove SST, remove Stainless generation, rip out obsolete CI, replace subprojects
- Then apply build system introduction and new wiring in follow-up commits.
- Then apply new features (agents, sandboxing) in small topic branches merged into dev.
- This sequencing minimizes conflicts and clarifies intent.

Subprojects replacement best practices

- Use git subtree or Git sparse directories where useful:
  - If pulling in external code you plan to vendor: git subtree add --prefix subproj <remote> <branch> --squash
  - If fully replacing with your own, do it in one commit labelled “Replace subproject X with Y (rationale)”
- If history churn is massive, consider:
  - git mv to preserve rename detection before edits
  - .gitattributes with:
    - renames: enable rename detection in diffs/reviews
    - linguist-vendored for large vendored directories to reduce review noise

Removing SST and Stainless cleanly

- Do removals in discrete commits with exhaustive cleanup:
  - Remove dependency from package.json/pnpm-lock, tsconfig paths, build scripts, CI, env files, infra directories, and types
  - Replace any generated artifacts with checked-in equivalents or new generation path
  - Add a MIGRATION.md noting replacements and env changes
- Run linters/tests after each removal commit to avoid “stuck” states.

Rebase vs merge policy

- Prefer rebase for dev onto track-upstream to keep a linear, readable history.
- Use merge commits only when upstream did sweeping refactors that make rebase painful.
- If rebase churn is heavy:
  - Create periodic integration points: dev-integration-YYYYMMDD that merge upstream once, then continue feature work on dev rebased onto that tag. Keeps conflict blasts scoped.

Handling recurring conflicts

- Maintain a conflict resolution doc (CONTRIBUTING-fork.md):
  - e.g., “If upstream reinstates SST hooks, remove them; our canonical build lives in scripts/build.ts”
- Add .gitattributes merge strategies where safe (e.g., ours for lockfiles, generated bundles):
  - package-lock.json merge=ours
  - dist/\* merge=ours
  - generated/\*\* merge=ours

CI strategy

- Two pipelines:
  - upstream-tracking CI: lightweight, ensures your fork can still build in upstream shape (track-upstream branch).
  - dev CI: your real pipeline (no SST, no Stainless), tests agents/sandboxing.

Release management

- Tag dev milestones (v0.x-dev.N) to create restore points.
- If you need to backport select upstream fixes:
  - Cherry-pick from track-upstream to dev, not the other way around.
  - Keep a small “backports.md” with upstream commit SHAs and rationale.

GitHub workflow

- Protect dev, require CI.
- Use PRs from topic/\* into dev for reviewability, even if you’re solo.
- Keep a draft PR “Fork: ongoing divergence” from dev to main (in your fork) describing the high-level migration; update its body as a running log.

Documentation hygiene

- Maintain MIGRATION.md explaining:
  - Removed: SST, Stainless
  - New: build system, envs, scripts, how to run tests/build
  - Subproject replacements and paths
- Add docs/architecture-fork.md for the new layout.

Quick command cheatsheet

- Sync upstream into dev:

```bash
  - git fetch upstream
  - git checkout track-upstream
  - git reset --hard upstream/main
  - git checkout dev
  - git rebase track-upstream
  - git push --force-with-lease origin dev:manno23/dev
```
- Publish clean upstream mirror to your fork:
```bash
  - git push --force-with-lease origin track-upstream:main
```
- Create topic branch:
```bash
  - git checkout -b topic/agents-sandbox dev
```
- Merge topic with a PR (preferred) or locally:
```bash
  - git checkout dev
  - git merge --no-ff topic/agents-sandbox
```

When to cut losses and stop rebasing

- If upstream drifts fundamentally (e.g., adopts a rebuild you discarded), consider freezing sync:
  - Tag last sync point.
  - Switch to cherry-picking only security/critical fixes.
  - Document the freeze in MIGRATION.md.

This keeps upstream pulls manageable, your fork coherent, and history readable while you overhaul subsystems.
