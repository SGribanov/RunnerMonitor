<!-- Issue: SGribanov/RunnerMonitor#35 -->
# Status

## Snapshot
- Current phase: M3 handoff
- Plan file: `tasks/issue-35-remote-only-runners/plan.md`
- Status: validated
- Last updated: 2026-06-07

## Done
- Synchronized repository state with `git fetch --all --prune`, `git status --short --branch`, and `git pull --rebase`.
- Refreshed GitHub self-hosted runner visibility through Exa using official GitHub docs.
- Confirmed current monitor only enriches locally discovered runners with GitHub status.
- Created canonical GitHub issue #35.
- Added issue #35 to GitHub Project board #25.
- Set GitHub Project Status to `In Progress`.
- Created branch `codex/issue-35-remote-only-runners`.
- Created task plan/status/test-plan files.
- Implemented `github-remote` rows for unmatched self-hosted runners returned by GitHub.
- Added read-only guards for lifecycle, cleanup, logs, removal, audit, and repo lifecycle paths.
- Added unit coverage for remote-only rows and local/remote matching.
- Updated README, README_RU, CHANGELOG, and technology insights.
- Synced technology insights to `D:\Repos\IdeaBox\vault\my-research\RunnerMonitor_technology_insights.md`.
- Verified IdeaBox watcher is running.
- Verified `go test ./internal/app`, `go test ./...`, `git diff --check`, and `go run ./cmd/runner-monitor --once`.
- Bumped `CurrentVersion` and download docs to `v0.5.0`.
- Rebuilt `bin\runner-monitor.exe`.
- Built `dist\RunnerMonitor-v0.5.0-windows-x64.zip` and `.sha256`.
- Verified the ZIP contents by extracting a temporary check folder.

## In Progress
- Commit, push, merge, and publish v0.5.0 release.

## Next
- Commit branch, push, merge into `main`, tag/release `v0.5.0`, and publish issue handoff.

## Decisions Made
- Remote-only self-hosted runners are read-only inventory rows, not controllable local runners.
- `github-remote` transport identifies GitHub-visible self-hosted runners without local folders.

## Assumptions In Force
- GitHub repository runner API contains the partner runners when our auth can access DeltaG runner administration data.
- The first implementation is repository-level, matching current monitored repo behavior.

## Commands
```sh
git fetch --all --prune
git status --short --branch
git pull --rebase
gh issue create ...
gh project item-add 25 --owner SGribanov --url https://github.com/SGribanov/RunnerMonitor/issues/35
gh project item-edit --project-id PVT_kwHOAL59L84BZm3n --id PVTI_lAHOAL59L84BZm3nzgu-Lts --field-id PVTSSF_lAHOAL59L84BZm3nzhUkWa8 --single-select-option-id 47fc9ee4
git switch -c codex/issue-35-remote-only-runners
go test ./internal/app
go test ./...
git diff --check
go run ./cmd/runner-monitor --once
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
.\bin\runner-monitor.exe --once
```

## Current Blockers
- None

## Audit Log
| Date | Milestone | Files | Commands | Result | Next |
| --- | --- | --- | --- | --- | --- |
| 2026-06-07 | Setup | `tasks/issue-35-remote-only-runners/*` | `gh issue create`; `gh project item-add`; `git switch -c` | issue/project/branch ready | M1 implementation |
| 2026-06-07 | Project sync | GitHub Project #25 | `gh project item-edit` | status set to In Progress | implementation |
| 2026-06-07 | M1-M2 | `internal/app/*` | `go test ./internal/app`; `go test ./...` | pass | docs/handoff |
| 2026-06-07 | M3 | README, CHANGELOG, research, IdeaBox vault | `git diff --check`; `go run ./cmd/runner-monitor --once` | pass; DeltaG shows `github`/`remote` rows | issue handoff |
| 2026-06-07 | Release prep | `internal/app/version.go`, `dist/*` | `scripts/build.ps1`; package script; `.\bin\runner-monitor.exe --once` | v0.5.0 exe/zip/sha built and smoke-tested | push/merge/release |

## Smoke / Demo Checklist
- [x] Remote-only self-hosted runner appears in inventory.
- [x] Local runner behavior is unchanged.
- [x] Remote-only rows are read-only.
- [x] Go tests pass.
- [x] v0.5.0 binary and ZIP are built.
- [ ] Issue handoff comment published.
