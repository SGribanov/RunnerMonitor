<!-- Issue: SGribanov/RunnerMonitor#10 -->
# GitHub-Ready Documentation Plan

## Scope

Prepare RunnerMonitor for GitHub presentation with an open license, English-first
README, Russian README, and baseline repository hygiene.

## Milestone 1 -- Repository Standards `[x]`

Goal: add files GitHub uses for license and project health signals.

Tasks:
- [x] Add a root `LICENSE` file.
- [x] Add `.gitignore` for Go, local config, build outputs, and OS/editor files.
- [x] Add contribution, security, code of conduct, issue, PR, and ownership
  guidance where appropriate.
- [x] Update GitHub repository description and topics.

Validation:
- `git status --short --branch`
- `gh repo view SGribanov/RunnerMonitor --json description,repositoryTopics,licenseInfo`

## Milestone 2 -- README English And Russian `[x]`

Goal: make `README.md` the complete English entry point and add equivalent
Russian documentation in `README_RU.MD`.

Tasks:
- [x] Rewrite `README.md` in English with purpose, stack, quick start,
  configuration, command reference, remote host usage, cleanup, reprovisioning,
  and safety notes.
- [x] Add `README_RU.MD` with Russian guidance and cross-links.
- [x] Keep sensitive local config guidance explicit.

Validation:
- Manual markdown review.
- `rg` checks for stale hard-coded user-secret guidance.

## Milestone 3 -- Validation And Handoff `[~]`

Goal: verify the docs changes do not break the repo and publish GitHub handoff.

Tasks:
- [x] Run Go tests and build.
- [x] Update research insights in repo and IdeaBox vault.
- [ ] Publish issue handoff.

Validation:
- `go test ./...`
- `scripts\build.ps1`
- `runner-monitor.ps1 --show-config`
