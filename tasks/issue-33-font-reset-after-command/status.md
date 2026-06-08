<!-- Issue: SGribanov/RunnerMonitor#33 -->
# Status

## Snapshot
- Current phase: M3 handoff
- Plan file: `tasks/issue-33-font-reset-after-command/plan.md`
- Status: validated
- Last updated: 2026-06-08

## Done
- Synchronized repository state with `git fetch --all --prune`.
- Reported dirty worktree from active branch `gsv777/issue-37-remote-runner-unregister`; no pull/rebase was run.
- Checked relevant GitHub issue #33 and linked it to Project #25.
- Set Project Status for #33 to `In Progress`.
- Verified `v0.5.0` source and release asset still contain `ESC[0m` + `ESC(B`.
- Refreshed Bubble Tea external command behavior through Exa using official Go package/source docs.
- Identified root cause: TUI mutating commands were executing external work outside Bubble Tea terminal release/restore.
- Implemented terminal-managed function actions through custom `tea.Exec`.
- Updated lifecycle command test for result-driven refresh.
- Verified `go test ./internal/app`.
- Synced technology insights to `D:\Repos\IdeaBox\vault\my-research\RunnerMonitor_technology_insights.md`.
- Verified IdeaBox watcher is running.
- Verified `go test ./...`, `git diff --check`, `scripts/build.ps1`, and `bin\runner-monitor.exe --once`.
- Published GitHub issue handoff comment: https://github.com/SGribanov/RunnerMonitor/issues/33#issuecomment-4645813958.
- Bumped release docs/version to `v0.5.1` and built release ZIP/checksum.

## In Progress
- Awaiting push/merge/release.

## Next
- Push/merge/release `v0.5.1`, then run real interactive terminal smoke.

## Decisions Made
- Keep frame-level reset as a rendered-frame guard.
- Add Bubble Tea terminal release/restore around interactive mutating actions.
- Do not change one-shot CLI commands in this fix.

## Assumptions In Force
- The reported regression happens in interactive TUI after commands such as `start`, `stop`, `clear`, `remove`, or `delete`.
- Live visual validation must be done in a real Windows terminal.

## Commands
```sh
git fetch --all --prune
git status --short --branch
gh issue view 33 --json number,title,projectItems,comments
gh issue edit 33 --add-project RunnerMonitor
gh project item-edit ...
go test ./internal/app
go test ./...
git diff --check
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
.\bin\runner-monitor.exe --once
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
```

## Current Blockers
- None

## Audit Log
| Date | Milestone | Files | Commands | Result | Next |
| --- | --- | --- | --- | --- | --- |
| 2026-06-08 | Setup | GitHub issue/project | `git fetch`; `gh issue view`; `gh issue edit`; `gh project item-edit` | #33 active on project board | inspect v0.5.0 |
| 2026-06-08 | M1-M2 | `internal/app/model.go`, tests | release asset check; Exa refresh; `go test ./internal/app` | reset present, terminal-managed action fix passes unit package | full validation |
| 2026-06-08 | M3 | docs, research, binary | `go test ./...`; `git diff --check`; `scripts/build.ps1`; `bin\runner-monitor.exe --once` | pass | issue handoff |
| 2026-06-08 | Handoff | GitHub issue #33 | `gh issue comment 33` | published handoff comment | manual visual smoke |
| 2026-06-08 | Release prep | version/docs/assets | `scripts/build.ps1`; package script | v0.5.1 ZIP/checksum built | push/merge/release |

## Smoke / Demo Checklist
- [x] `v0.5.0` release asset contains previous text-mode reset.
- [x] TUI mutating commands now use Bubble Tea `tea.Exec`.
- [x] Package tests pass.
- [x] Full Go tests pass.
- [x] Local binary rebuilt.
- [ ] Real terminal visual smoke completed by operator.
- [x] Handoff comment published.
