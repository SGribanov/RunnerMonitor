<!-- Issue: SGribanov/RunnerMonitor#1 -->
# RunnerMonitor Status

## Current phase

Milestone 2: TUI commands.

## Done

- Created GitHub repository `SGribanov/RunnerMonitor`.
- Created GitHub Project `RunnerMonitor` #25.
- Created canonical issue #1 and moved it to `In Progress`.
- Created branch `codex/1-runner-monitor-tui`.
- Confirmed Go `1.26.2` is available.
- Added Go/Bubble Tea application skeleton.
- Added local Windows runner discovery from `.runner` files and Windows services.
- Added WSL runner discovery from `.runner` files and `.service` hints.
- Added GitHub API status merge and queued/stale workflow counts.
- TUI/audit tables now show an explicit `Project` column.
- Windows manual runner processes now appear as `Local=running` while remaining `ControlMode=manual`.
- Manual Windows `run.cmd` runners can now be started/stopped/restarted by RunnerMonitor in a hidden background process.
- Added generated hourglass icon assets and embedded the static `.ico` into the Windows exe through a `.syso` resource.
- TUI startup now shows animated hourglass symbols and `Ожидайте, идет опрос раннеров...` while the first runner refresh runs asynchronously.
- Documented SSH commands for opening RunnerMonitor on a future dedicated remote runner host and starting remote project runners.
- Added saved remote host setup via `--configure-remote NAME`, SSH connect via `--connect-remote NAME`, and TUI command `connect remote NAME`.
- Added TUI commands: `refresh`, `start N`, `stop N`, `restart N`, `logs N`, `q`.
- Added `--once` smoke mode.
- Added `--audit`, `--start-repo`, `--stop-repo`, `--restart-repo`, and `--disable-autostart`.
- `go test ./...` passes.
- `go run ./cmd/runner-monitor --once` lists 11 current runner records and shows `DeltaG` queue as `1/1 stale`.
- `go run ./cmd/runner-monitor --audit` classifies cleanup candidates.

## In progress

- Hardening DeltaG queued-job diagnostics and preparing commit/push.
- Autostart disable requires elevated Windows/root WSL permissions.

## Next

- Investigate the DeltaG stale queued workflow and label/routing.
- Keep CLI commands machine-friendly; the loading animation is limited to interactive TUI mode.
- Run autostart disable from elevated context or document handoff commands.

## Decisions

- Use Go for a small Windows/Linux binary.
- Use local TUI over SSH for the future dedicated runner machine.
- Use GitHub CLI authentication for v1 instead of introducing app credentials.

## Remote runner host commands

```powershell
ssh -t runnerbox "powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1"
ssh runnerbox "cd C:/Repos/DeltaG; powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1 --start-current"
ssh runnerbox "powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1 --start-repo SGribanov/DeltaG"
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --configure-remote runnerbox
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --connect-remote runnerbox
```

## Commands

```powershell
go test ./...
go run ./cmd/runner-monitor
go run ./cmd/runner-monitor --once
go run ./cmd/runner-monitor --audit
go run ./cmd/runner-monitor --start-repo SGribanov/DeltaG
go run ./cmd/runner-monitor --disable-autostart
```

## Blockers

- None currently.

## Audit log

- 2026-06-03: Initialized repo/project/issue and started implementation.
- 2026-06-03: Completed first working inventory smoke with queued job counts.
- 2026-06-03: Attempted autostart disable. Windows services failed with access denied; WSL systemd disable failed with interactive authentication required.
