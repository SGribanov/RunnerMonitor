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

## 2026-06-03 -- Project column

The TUI/audit table should show the project explicitly. A short `Project` column
derived from the GitHub repo name is easier to scan than only showing the full
`owner/repo` string.

## 2026-06-03 -- Audit policy

RunnerMonitor supports a small `runner-policy.json` keep list. This lets the
operator preserve known future-use runners, such as `SGribanov/AU windows-local`,
without hard-coding exceptions into audit logic.

## 2026-06-03 -- MyClone WSL service

The existing `/home/gsv777/myclone-runner-linux` folder could be safely reused:
back up local runner config, re-register with a GitHub registration token using
`config.sh --replace`, then install/start the official WSL systemd service with
`svc.sh`. GitHub reported `mycloneosengine-linux` online as runner id `24`.

## 2026-06-03 -- Manual Windows process state

Manual Windows runners can be online because `Runner.Listener.exe` is running
from the runner folder even when no Windows service exists. Discovery should
map `Runner.Listener.exe` executable paths back to runner roots so the TUI shows
`Local=running` instead of only `manual`.

## 2026-06-03 -- Manual Windows lifecycle control

Manual Windows runners do not need to be converted to Windows Services before
RunnerMonitor can control them. For `ControlMode=manual` and `Transport=windows`,
the app can start `run.cmd` with `Start-Process -WindowStyle Hidden` and stop
only runner processes whose executable paths are inside the specific runner
folder. This keeps `BackTester/backtester-runner` and
`MyCloneOsEngine/mycloneosengine-windows-local` usable through the same TUI and
`--start-current` workflow while preserving the busy-runner stop protection.

## 2026-06-03 -- DeltaG stale queue diagnosis

DeltaG's remaining queued run `26447257991` has no jobs according to the GitHub
jobs API and belongs to closed PR branch `codex/604-vertical-freshness-diagnostic`.
At diagnosis time, both DeltaG self-hosted runners were online/not busy. This
points to a stale GitHub Actions run rather than a local runner availability
problem. The next engineering choice is to add a queue-diagnostics path that can
distinguish closed-PR/no-job stale runs from real queued work.

An explicit GitHub cancel call returned HTTP `409` with `Cannot cancel a
workflow re-run that has not yet queued`; the run remained `queued`. RunnerMonitor
should therefore treat this class as a stale/anomalous GitHub run rather than
expecting the normal cancel endpoint to clear it.

## 2026-06-03 -- Windows executable icon resource

Windows shell executable icons are static, not animated. RunnerMonitor now keeps
both an animated GIF asset (`assets/runner-monitor-hourglass-spin.gif`) and a
static `.ico` for the executable (`assets/runner-monitor-hourglass.ico`). The
static icon is embedded into Go's Windows build through
`cmd/runner-monitor/runner-monitor_windows_amd64.syso`, generated by
`github.com/akavel/rsrc`. Knowledge refresh for this resource-embedding choice
was performed through `search_ai_mcp_default`; local validation confirmed the
generator's `.syso` output is automatically linked by `go build` when placed in
the package directory.
