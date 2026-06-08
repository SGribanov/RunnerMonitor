<!-- Issue: SGribanov/RunnerMonitor#33 -->
# Plans

## Source
- Task: Check and repair `v0.5.0` regression where TUI font/display reset protection breaks after command execution.
- Canonical issue: GitHub issue #33.
- Repo context: Go Bubble Tea TUI in `internal/app/model.go`; lifecycle commands can spawn PowerShell, cmd, WSL, systemd, and service-control processes.
- Last updated: 2026-06-08

## Assumptions
- The symptom is in the interactive TUI command path, not one-shot CLI output.
- The previous frame reset is useful but insufficient around external process execution.
- Fix should preserve existing post-command status refresh.

## Milestone Order
| ID | Title | Depends on | Status |
| --- | --- | --- | --- |
| M1 | Verify v0.5.0 reset state | - | [x] |
| M2 | Add terminal-managed TUI actions | M1 | [x] |
| M3 | Validate and handoff | M2 | [x] |

## M1. Verify v0.5.0 reset state `[x]`
### Goal
- Confirm whether `v0.5.0` lost the text-mode reset or whether the gap is elsewhere.

### Tasks
- [x] Inspect `v0.5.0` source and release asset for `ESC[0m` + `ESC(B`.
- [x] Compare `v0.4.0..v0.5.0` around TUI command execution.
- [x] Refresh Bubble Tea external execution behavior through Exa using official docs/source.

### Definition of Done
- Root cause hypothesis is grounded in code and upstream docs.

### Validation
```sh
git show v0.5.0:internal/app/model.go
gh release download v0.5.0 --pattern RunnerMonitor-v0.5.0-windows-x64.zip
```

### Known Risks
- Visual terminal regressions need manual smoke in a real Windows terminal.

### Stop-and-Fix Rule
- If the release asset lacks the reset sequence, fix packaging before code behavior.

## M2. Add terminal-managed TUI actions `[x]`
### Goal
- Run mutating interactive TUI commands under Bubble Tea terminal release/restore.

### Tasks
- [x] Add a custom `tea.Exec` command wrapper for function-based RunnerMonitor actions.
- [x] Use it for lifecycle, clear, remove, delete, clear idle, and auto-clear actions.
- [x] Preserve lifecycle result messages and post-command refresh.
- [x] Update regression tests for the new result/refresh order.

### Definition of Done
- External operations no longer run directly inside `Update` or plain command callbacks.

### Validation
```sh
go test ./internal/app
```

### Known Risks
- `tea.Exec` releases the terminal while the operation runs, so the command path may briefly leave the alternate screen.

### Stop-and-Fix Rule
- If lifecycle messages disappear behind refresh, restore message preservation before handoff.

## M3. Validate and handoff `[x]`
### Goal
- Run full validation, rebuild local binary, and publish issue handoff.

### Tasks
- [x] Update changelog, task docs, and technology insights.
- [x] Sync technology insights to IdeaBox vault.
- [x] Run full Go tests and diff check.
- [x] Rebuild `bin\runner-monitor.exe`.
- [x] Publish issue handoff comment.

### Definition of Done
- User has a rebuilt local binary with stronger TUI terminal-state protection.

### Validation
```sh
go test ./...
git diff --check
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
```

### Known Risks
- A real terminal visual smoke is still needed for final confidence.

### Stop-and-Fix Rule
- If tests fail or the binary does not build, fix before handoff.
