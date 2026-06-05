<!-- Issue: SGribanov/RunnerMonitor#28 -->
# Status

## Snapshot
- Current phase: Implementation complete; live smoke optional/manual
- Plan file: `tasks/issue-28-ensure-start-online/plan.md`
- Status: green
- Last updated: 2026-06-05

## Done
- Created GitHub issue #28 for the new task.
- Refreshed GitHub Actions runner systemd behavior via `exa` using official GitHub docs.
- Implemented service-managed `start` as enable/start/local-active/GitHub-online readiness.
- Added regression tests for WSL systemd enable/start and GitHub offline timeout reporting.
- Updated technology insights in repo and IdeaBox vault.
- Verified IdeaBox watcher is running.

## In Progress
- None

## Next
- Optional live smoke: stop/disable a WSL runner unit, run `--start-current` or `start [N]`, and confirm GitHub reports `online`.

## Decisions Made
- Use GitHub runner `online` status as the application-level "listening for jobs" signal.
- Enable service autostart before starting systemd-managed runners so a disabled unit is corrected.

## Assumptions In Force
- The existing `gh` CLI integration remains the status source for GitHub runner checks.
- The fix must preserve existing uncommitted issue #27 changes.

## Commands
```sh
go test ./internal/app
go test ./...
```

## Current Blockers
- None

## Audit Log
| Date | Milestone | Files | Commands | Result | Next |
| --- | --- | --- | --- | --- | --- |
| 2026-06-05 | M1-M2 | `internal/app/lifecycle.go`, `internal/app/github.go`, `internal/app/app_test.go`, `research/RunnerMonitor_technology_insights.md` | `go test ./internal/app`; `go test ./...` | pass | Optional live smoke |

## Smoke / Demo Checklist
- [x] WSL systemd runner start enables disabled service.
- [x] Start waits for local active state.
- [x] Start waits for GitHub `online`.
- [x] Go tests pass.
- [ ] Optional live WSL smoke run.
