<!-- Issue: SGribanov/RunnerMonitor#34 -->
# Status

## Snapshot
- Current phase: M3 release
- Plan file: `tasks/issue-34-hosted-runners/plan.md`
- Status: validating release
- Last updated: 2026-06-07

## Done
- Synchronized repository state with `git fetch --all --prune`.
- Confirmed local branch was aligned with `origin/main` but worktree had existing local modifications.
- Refreshed GitHub Actions hosted/self-hosted API context via Exa against official GitHub documentation.
- Created canonical GitHub issue #34.
- Added issue #34 to GitHub Project board #25.
- Created branch `codex/issue-34-hosted-runners`.
- Implemented GitHub-hosted workflow job discovery through queued/in-progress workflow runs/jobs.
- Added read-only guards for hosted rows across lifecycle, cleanup, remove, audit, and logs paths.
- Updated README, README_RU, CHANGELOG, version, build defaults, and technology insights.
- Synced technology insights to `D:\Repos\IdeaBox\vault\my-research\RunnerMonitor_technology_insights.md`.
- Verified IdeaBox watcher is running.
- Built `bin\runner-monitor.exe` and `dist\RunnerMonitor-v0.4.0-windows-x64.zip`.
- Verified `go test ./internal/app`, `go test ./...`, and `bin\runner-monitor.exe --once`.

## In Progress
- Commit, push, and publish the v0.4.0 GitHub release.

## Next
- Review git status/diff, commit, push, create tag/release, and publish issue handoff.

## Decisions Made
- Hosted runners will be represented through workflow job state, not local lifecycle state.
- Hosted rows are read-only and must not invoke service, folder, or runner-registration operations.

## Assumptions In Force
- First release focuses on repository workflow jobs visible through `gh api`.
- Existing uncommitted changes are treated as user work and preserved.

## Commands
```sh
git fetch --all --prune
git status --short --branch
gh issue create ...
gh project item-add 25 --owner SGribanov --url https://github.com/SGribanov/RunnerMonitor/issues/34
go test ./internal/app
go test ./...
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
.\bin\runner-monitor.exe --once
```

## Current Blockers
- None

## Audit Log
| Date | Milestone | Files | Commands | Result | Next |
| --- | --- | --- | --- | --- | --- |
| 2026-06-07 | Setup | `tasks/issue-34-hosted-runners/*` | `gh issue create`; `gh project item-add` | issue/project created | M1 implementation |
| 2026-06-07 | M1-M3 | `internal/app/*`, docs, research, build output | `go test ./...`; `.\bin\runner-monitor.exe --once`; `scripts/build.ps1` | pass | commit/release |

## Smoke / Demo Checklist
- [x] Hosted queued/in-progress jobs appear in refresh output.
- [x] Hosted rows are read-only.
- [x] Go tests pass.
- [x] Release build passes.
- [ ] Release is published.
