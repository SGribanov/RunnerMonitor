# RunnerMonitor -- Technology Insights

| Field | Value |
|---|---|
| Project | RunnerMonitor |
| Type | technology-research |
| Last updated | 2026-06-03 |
| Status | active |
| Tags | go, bubble-tea, github-actions, wsl, windows-services |

## 2026-06-03 -- Dual-state runner monitoring

GitHub runner health is not enough on its own. The current workstation has
service-managed runners, configured-but-manual runner directories, WSL systemd
units, and at least one stale queued `DeltaG` workflow even while matching
runners are online. RunnerMonitor therefore needs to display local lifecycle
state, GitHub `online/offline/busy`, labels, version, and queued/stale workflow
counts together.

## 2026-06-03 -- First implementation stack

Go `1.26.2` is available locally and works for a small cross-platform TUI.
Bubble Tea dependencies resolve successfully. The first smoke command,
`go run ./cmd/runner-monitor --once`, discovers 11 runner records across local
Windows and WSL and reports `SGribanov/DeltaG` queue health as `1/1 stale`.

## 2026-06-03 -- WSL runner background behavior

WSL runners should run through systemd units, not visible terminal windows.
Starting a repository through RunnerMonitor should call `systemctl start` inside
WSL so the runner stays in the background. Disabling boot autostart for current
Windows and WSL services requires elevated/admin or root authentication; the
non-elevated session cannot change those startup policies.

## 2026-06-03 -- Wrapper command

`runner-monitor.ps1` builds `bin\runner-monitor.exe` on first use and forwards
arguments to it. This gives Codex a stable command path:
`powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --start-current`.

## 2026-06-03 -- Already-running services

Starting an already running Windows service or active WSL unit can still require
permissions if the command blindly calls `Start-Service` or `systemctl start`.
RunnerMonitor should short-circuit `start` when local state is already
`running` or `active`, returning `already running` instead of invoking privileged
service control.
