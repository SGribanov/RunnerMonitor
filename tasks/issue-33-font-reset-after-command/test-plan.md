<!-- Issue: SGribanov/RunnerMonitor#33 -->
# Test Plan

## Source
- Task: Repair TUI font/display reset after command execution.
- Plan file: `tasks/issue-33-font-reset-after-command/plan.md`
- Status file: `tasks/issue-33-font-reset-after-command/status.md`
- Last updated: 2026-06-08

## Validation Scope
- In scope: model command dispatch, terminal-managed command wrapper, lifecycle result refresh, frame reset retention.
- Out of scope: automated pixel/terminal font assertion; final visual check requires a real Windows terminal.

## Environment / Fixtures
- Unit fixtures use synthetic runners.
- Live smoke should run the rebuilt TUI in Windows Terminal or the user's normal shell.

## Test Levels

### Unit
- TUI `View()` still starts and ends with the text-mode reset.
- Lifecycle command dispatch returns a terminal-managed command and shows an in-progress message.
- Lifecycle command result triggers the status refresh and preserves the command result message.

### Integration
- `go test ./...` verifies all packages.

### Smoke
- Rebuild local binary with `scripts/build.ps1`.
- Run the TUI, execute a safe command such as dry-run `remove N` or a harmless unsupported command on a non-service runner, and confirm terminal/font display remains normal.

## Negative / Edge Cases
- Terminal restore returns an error.
- Lifecycle result requests refresh while a refresh is already active.
- Clear/remove/delete commands complete without a post-command refresh.

## Acceptance Gates
- [x] `go test ./internal/app`
- [x] `go test ./...`
- [x] `git diff --check`
- [x] `.\scripts\build.ps1`

## Smoke Result
- [x] `bin\runner-monitor.exe --once` starts and renders inventory successfully.
- [ ] Real interactive terminal visual smoke remains manual.

## Command Matrix
```sh
go test ./internal/app
go test ./...
git diff --check
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
```
