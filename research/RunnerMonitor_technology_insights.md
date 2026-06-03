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

## 2026-06-03 -- TUI startup loading state

The first runner refresh can take long enough that a blank terminal feels like a
hang. Interactive TUI startup now begins immediately with a loading model that
animates `⌛`/`⏳` and shows `Ожидайте, идет опрос раннеров...` while `Refresh()`
runs as a Bubble Tea command. CLI modes such as `--once`, `--audit`, and
project lifecycle commands still refresh synchronously and do not print spinner
frames, preserving machine-friendly output.

## 2026-06-03 -- Remote SSH command shape

Remote runner host support does not require a daemon for the first migration
step. The operator can allocate a TTY and run the TUI on the remote Windows host:

```powershell
ssh -t runnerbox "powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1"
```

Non-interactive Codex/project startup can run:

```powershell
ssh runnerbox "cd C:/Repos/DeltaG; powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1 --start-current"
```

Future `--host` support can wrap this command shape rather than changing the TUI
command model.

## 2026-06-03 -- Saved remote host configuration

RunnerMonitor now has a minimal saved SSH host workflow. `--configure-remote
NAME` prompts for remote name, SSH host/alias, host OS, remote RunnerMonitor
path, and default project path, then writes `remote-hosts.json` under the user's
config directory. `--connect-remote NAME` reuses that config and opens the
remote TUI over `ssh -t`. Inside the local TUI, `connect remote NAME` uses the
same saved config through Bubble Tea `ExecProcess`, temporarily handing the
terminal to SSH and returning to RunnerMonitor after the remote session closes.

## 2026-06-03 -- Runner folder migration constraints

Moving GitHub Actions runner folders is path-sensitive. Manual Windows runners
can be moved with backup, stop, move, and restart if idle. Windows service and
WSL systemd runners require service/unit reinstallation or re-registration from
the new path, often with elevated PowerShell or sudo. RunnerMonitor issue #5
tracks the migration into `C:\Runners` and `/home/gsv777/Runners`; each runner
move needs a runner-specific rollback plan and validation with `--audit` plus
GitHub online status.

## 2026-06-03 -- BackTester runner move

`SGribanov/BackTester backtester-runner` was moved from
`C:\actions-runner-backtester` to
`C:\Runners\SGribanov-BackTester\backtester-runner` after backup
`C:\Runners-backup\actions-runner-backtester-backtester-runner-2026-06-03.zip`.
The first post-move start failed because `bin` and `externals` were Windows
junctions still targeting `C:\actions-runner-backtester\bin.2.334.0` and
`C:\actions-runner-backtester\externals.2.334.0`. Retargeting those junctions
to the versioned folders under the new root fixed startup; GitHub then reported
the runner online and `--start-current` from `C:\Repos\BackTester` returned
`backtester-runner already running`.

## 2026-06-03 -- MyClone Windows runner move

`SGribanov/MyCloneOsEngine mycloneosengine-windows-local` was moved from
`C:\actions-runner-mycloneosengine` to
`C:\Runners\SGribanov-MyCloneOsEngine\mycloneosengine-windows-local` after
backup
`C:\Runners-backup\actions-runner-mycloneosengine-mycloneosengine-windows-local-2026-06-03.zip`.
As with BackTester, `bin` and `externals` were junctions targeting the old root.
They were retargeted to `bin.2.334.0` and `externals.2.334.0` under the new
root before starting the runner. GitHub reported the runner online/busy=false,
and `--start-current` from `C:\Repos\MyCloneOsEngine` returned both MyClone
runners already running.
