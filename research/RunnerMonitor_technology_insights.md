# RunnerMonitor -- Technology Insights

| Field | Value |
|---|---|
| Project | RunnerMonitor |
| Type | technology-research |
| Last updated | 2026-06-05 |
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

## 2026-06-04 -- Elevated Windows TUI rendering

When RunnerMonitor is launched from an elevated Windows PowerShell window, the
normal Bubble Tea terminal buffer can show only the bottom of the rendered
frame. The interactive TUI should run with Bubble Tea alternate screen mode so
the app owns the full viewport during redraws. The model should also keep the
actual small terminal height instead of forcing it to 12 rows; otherwise a tiny
admin console tries to render a full table and pushes the status/input area to
the only visible lines. For heights of 8 rows or less, a compact view that shows
title, status, and input is preferable to a broken table.

## 2026-06-04 -- Startup update check

The TUI can check for a newer RunnerMonitor release once at startup without
blocking runner monitoring. Reusing `gh release view --repo SGribanov/RunnerMonitor
--json tagName,url` avoids adding a new SDK or HTTP client dependency and stays
consistent with the existing GitHub CLI requirement. The check should use a
short timeout and treat offline, unauthenticated, or failed release lookups as
silent no-ops. A newer version should render as a separate notice line rather
than replacing the normal runner status message.

## 2026-06-04 -- Review hardening follow-up

Security review found that destructive WSL/Linux folder deletion must normalize
slash paths with `path.Clean` before configured-root checks. String prefix
checks alone can accept traversal such as `/runnerbox/Runners/../danger`.
Runner root folders themselves should not be deletable; only children under
configured roots are valid targets. Runner registration/remove tokens should be
redacted from command errors and kept out of parent process argument lists where
the upstream `config.cmd`/`config.sh` calling convention allows it. TUI
auto-refresh should avoid repeated `gh` process fan-out every few seconds by
briefly caching GitHub status for automatic refreshes while keeping manual
refresh fresh.

## 2026-06-04 -- Windows discovery subprocess timeout

After the `v0.3.0` release, a PowerShell window with the inherited
`runner-monitor.exe` title remained visible after TUI exit. Process inspection
showed the lingering command was Windows service discovery:
`Get-CimInstance Win32_Service ... ConvertTo-Json`. Discovery commands should be
bounded with `context.WithTimeout`; otherwise a slow or stuck CIM provider can
outlive Bubble Tea shutdown and keep the launch window around. Race validation
is still a local toolchain issue rather than a code validation signal until a
working cgo/GCC toolchain is available.

## 2026-06-04 -- Windows race-test toolchain

Go `go test -race ./...` on Windows failed with the existing MCF MinGW GCC
16.1.1 toolchain even after `gcc.exe` was on `PATH`; `runtime/cgo` failed in the
assembler with `junk '(eax)' after expression`. Installing MSYS2 and
`mingw-w64-ucrt-x86_64-gcc` fixed the race-test toolchain. User `PATH` should
prefer `C:\msys64\ucrt64\bin` over `C:\Soft\mingw-w64-gcc-mcf\mingw64\bin`.
With MSYS2 UCRT64 GCC 16.1.0, `go test -race ./...` passes.

## 2026-06-04 -- TUI status visibility

The first Windows discovery timeout fix used a 5-second bound, but local
service/process discovery can exceed that during startup. The symptom was
incomplete TUI status for `ideabox-runner` even though GitHub returned
`online`/`busy=true`. A 30-second bound still prevents indefinite child
PowerShell hangs while preserving normal status discovery. Avoid ANSI-styled
strings inside `bubbles/table` cell values; plain `true`/`false` keeps table
width accounting and selected-row rendering predictable.

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

## 2026-06-03 -- MyClone WSL runner move

`SGribanov/MyCloneOsEngine mycloneosengine-linux` was moved from
`/home/gsv777/myclone-runner-linux` to
`/home/gsv777/Runners/SGribanov-MyCloneOsEngine/mycloneosengine-linux` after
backup
`/home/gsv777/runner-backups/myclone-runner-linux-mycloneosengine-linux-move-2026-06-03.tar.gz`.
The old `svc.sh uninstall` removed the systemd unit, then `svc.sh install
gsv777` from the new path recreated it with `ExecStart` and `WorkingDirectory`
under the new root. As on Windows, the moved runner had path-sensitive
`bin`/`externals` links; they needed to be recreated to point to `bin.2.334.0`
and `externals.2.334.0` under the new root before `svc.sh install` could find
`bin/actions.runner.service.template`. After `svc.sh start`, GitHub reported
the runner online/busy=false, `runner-monitor --audit` showed `keep`, and
`--start-current` from `C:\Repos\MyCloneOsEngine` returned both MyClone runners
already running.

## 2026-06-03 -- DeltaG WSL runner move

`SGribanov/DeltaG deltag-linux-wsl` was moved from
`/home/gsv777/actions-runner-deltag` to
`/home/gsv777/Runners/SGribanov-DeltaG/deltag-linux-wsl` after backup
`/home/gsv777/runner-backups/actions-runner-deltag-deltag-linux-wsl-move-2026-06-03.tar.gz`.
The runner folder was about 11 GB, so the backup tar step took several minutes;
checking the active `tar` process and archive size was the right way to
distinguish a long backup from a hung migration. After `svc.sh uninstall`,
moving the folder, recreating `bin` and `externals` symlinks to the versioned
folders under the new root, and reinstalling via `svc.sh install gsv777`, the
systemd unit started with `ExecStart` under the new path. GitHub reported both
DeltaG runners online/busy=false, and `--start-current` from `C:\Repos\DeltaG`
returned both DeltaG runners already running. The repo still reports a stale
queued run, but this is now decoupled from local runner folder placement.

## 2026-06-03 -- AU runner move and reattach

`SGribanov/AU windows-local` was moved from `C:\actions-runner` to
`C:\Runners\SGribanov-AU\windows-local` after backup
`C:\Runners-backup\actions-runner-windows-local-move-2026-06-03.zip`. The local
folder still contained `.runner` credentials for `SGribanov/AU`, but GitHub
reported zero registered runners for the repo, so the migration needed a
GitHub reattach rather than only a path move. Running `config.cmd remove` with a
fresh remove token cleaned the stale local registration, and `config.cmd
--replace` with a fresh registration token recreated `windows-local`. The
post-move `bin` and `externals` junctions again needed retargeting to the
versioned folders under the new root. After fixing RunnerMonitor's manual
Windows launch path handling, `--start-current` from `C:\Repos\AU` started the
runner; GitHub reported `windows-local` online/busy=false.

## 2026-06-03 -- Manual Windows PowerShell path passing

Passing the runner path as a positional argument to `powershell -Command` was
fragile: PowerShell executed the path as a follow-up command and left
`$RunnerPath` empty inside the script, breaking `Join-Path`. Manual Windows
runner start/stop now pass the runner path through a process environment
variable, which avoids quoting edge cases for paths such as
`C:\Runners\SGribanov-AU\windows-local`.

## 2026-06-03 -- Windows service migration blocker

The remaining Windows service-managed runners, `SGribanov/IdeaBox
ideabox-runner` and `SGribanov/DeltaG deltag-win`, require elevated PowerShell
for both autostart changes and folder migration. In a non-elevated shell
(`Administrator=False`), `sc.exe config ... start= demand` returns
`OpenService FAILED 5: Access is denied`. These runners should not be moved
until an elevated session can stop the service, update or reinstall the service
from the new `C:\Runners` path, set startup to manual, and restart it.

## 2026-06-04 -- Windows service runner migration

`SGribanov/IdeaBox ideabox-runner` and `SGribanov/DeltaG deltag-win` were moved
with elevated PowerShell. Updating only the service `binPath` to the moved
`bin\RunnerService.exe` is sufficient when the runner registration remains
valid and the path-sensitive `bin`/`externals` junctions are retargeted to the
versioned directories under the new root. `IdeaBox` moved cleanly with
`Move-Item`. `DeltaG` needed a recovery path: after a successful backup,
`Move-Item` failed with access denied on `C:\github-runners\deltag`, so the
folder was copied to `C:\Runners\SGribanov-DeltaG\deltag-win` with `robocopy`,
then the service was reconfigured and started from the new path. The old
`C:\github-runners\deltag` folder was removed only after GitHub reported the
new service online and idle.

## 2026-06-04 -- Runner folder cleanup

After all runner bindings were verified online/busy=false, backup archives were
removed from `C:\Runners-backup` and `/home/gsv777/runner-backups`. Windows and
WSL runner `_work` directories were cleared, and installer archives inside
runner roots were removed. Windows `C:\Runners` dropped from about 16.5 GB to
about 2.9 GB, `C:\github-runners` was removed, WSL runner `_work` directories
are about 4 KB each, and `/home/gsv777/runner-backups` is about 4 KB. Manual
Windows runners should be restarted from non-elevated project commands after
elevated cleanup; otherwise non-elevated discovery cannot see their executable
paths and the audit may temporarily show `manual` despite GitHub being online.

## 2026-06-04 -- Safe runner cleanup command

Runner cleanup should preserve registration state and runner binaries. The safe
surface is `_work` contents plus root-level runner installer archives such as
`actions-runner*.zip` and `actions-runner*.tar.gz`. RunnerMonitor now refuses
cleanup when GitHub reports `busy=true`; if a runner is locally running or
active, it stops only that runner, clears the safe targets, and restarts it.
The TUI supports `clear N`, `clear idle`, and opt-in `auto-clear on/off`.
Automatic cleanup is tied to refresh instead of a permanent background loop, so
cleanup happens only after a fresh idle check and avoids silently racing new
jobs.

## 2026-06-04 -- Cleanup service-control fallbacks

Non-elevated Windows PowerShell cannot stop service-managed GitHub runners, so
Windows service cleanup must fail before deleting `_work` and clearly say that
elevated PowerShell is required. Setting PowerShell output encoding to UTF-8
keeps that error readable in RunnerMonitor output. WSL systemd cleanup needs a
sudo fallback because `systemctl stop` can return interactive-authentication
errors. Passing the sudo password through stdin to `wsl.exe -- sudo -S ...`
works reliably; WSL shell argument and stdin forwarding through `sh -c` were
not reliable in this environment, so WSL folder cleanup uses `wsl.exe -- python3
-c` with the runner path passed as base64.

## 2026-06-04 -- Elevated Windows cleanup helper

For service-managed Windows runners, a non-elevated TUI cannot safely perform
the stop-clean-start sequence because `Stop-Service` needs administrator rights.
RunnerMonitor now detects that case before deleting anything and opens an
elevated PowerShell helper through UAC for only the selected runner. The helper
uses `--clear-runner NAME`, keeps the elevated window open with `-NoExit`, and
preserves the existing busy-runner and safe-target cleanup rules. WSL cleanup
does not use this path; it continues to rely on the sudo fallback.

## 2026-06-04 -- Runner removal and reprovisioning

GitHub Actions runner reprovisioning should use official repository
registration/remove tokens instead of editing `.runner` state directly. Tokens
are fetched with `gh api -X POST repos/<owner>/<repo>/actions/runners/
registration-token` or `remove-token`, then passed to `config.cmd`/`config.sh`
or `config remove`. RunnerMonitor v1 makes these flows dry-run by default and
requires `--confirm` before unregistering or configuring. Folder deletion is a
second explicit gate with `--delete-folder` and is limited to known runner roots.

The requested project selector maps a plain folder name to `C:\Repos\<Project>`
and reads that repository's GitHub `origin`. It deliberately rejects path-like
input (`..`, slashes, drive prefixes) so a destructive command cannot escape the
project registry by accident. Blindly cloning existing runner folders remains a
later milestone because current runner directories can contain auto-update
junctions or symlinks; v1 configures an existing prepared runner distribution
folder instead.

## 2026-06-04 -- App-local settings file

RunnerMonitor now uses `runner-monitor.json` next to the compiled executable,
not a roaming user-config path. This fits the planned dedicated runner machine:
the executable and config can be copied as one app folder. `RUNNER_MONITOR_CONFIG`
remains an override for tests or unusual launch contexts. The app-local config
contains `projectsRoot`, Windows/WSL/Linux runner roots, and the direct
`wslSudoPassword` value requested by the operator. `--show-config` masks that
password as `<set>`/`<empty>`, and `scripts\build.ps1` creates a default config
beside `runner-monitor.exe` without overwriting an existing one.

## 2026-06-04 -- TUI table direction

The current TUI problems come from a hand-rendered fixed-width table, not from
Bubble Tea itself. The better next step is to keep the Charmbracelet stack and
replace the table surface with `bubbles/table`, using `tea.WindowSizeMsg` to
track terminal dimensions and recompute columns. If the command/status area
needs scrolling, `bubbles/viewport` should own that bounded region. This keeps
the app lightweight while making resize behavior deterministic.

## 2026-06-04 -- Resize-aware TUI implementation

RunnerMonitor now keeps the non-interactive `--once` output as the existing
plain text table, but the interactive Bubble Tea model owns a `bubbles/table`
component. The model records terminal width/height from `tea.WindowSizeMsg`,
recomputes table columns and height on resize, and keeps the command input
width within the terminal. Wide terminals give more space to project, runner,
labels, and path; narrow terminals hide low-priority columns and shrink
remaining columns instead of letting rows drift. Existing numeric commands still
work, and commands such as `start`, `stop`, `clear`, and `logs` can now target
the selected row when no number is provided.

## 2026-06-04 -- GitHub community profile files

GitHub repository best-practice guidance was refreshed through GitHub Docs:
README, license, contribution guidelines, and code of conduct are the core
signals that communicate project expectations. RunnerMonitor now adds a root
MIT `LICENSE`, English `README.md`, Russian `README_RU.MD`, `CONTRIBUTING.md`,
`SECURITY.md`, `CODE_OF_CONDUCT.md`, issue templates, PR template,
`CODEOWNERS`, Dependabot config for Go modules, and an expanded `.gitignore`
that excludes generated builds and local config/secrets. GitHub repo metadata
was also updated with a description and topics for discoverability.

## 2026-06-04 -- TUI auto-refresh loop

The interactive Bubble Tea model now schedules inventory refresh with a
single-shot `tea.Tick` loop. The interval comes from
`tuiRefreshIntervalSeconds` in `runner-monitor.json` and defaults to 5 seconds
when omitted or invalid. The model tracks `refreshing` so manual and automatic
refresh commands cannot overlap. Automatic refresh does not enter the
loading-only screen, and manual `refresh` keeps the current table visible when
existing data is present; the first startup refresh still shows the wait screen
because there is no useful old inventory yet.

## 2026-06-04 -- Clickable update notice URL

Terminal hyperlinks should be rendered with OSC-8 escape sequences around the
visible URL label. The TUI must not truncate the full rendered string after
adding OSC-8 codes, because that can cut the closing sequence and break the
link. Instead, truncate only the visible label before wrapping it in the
hyperlink sequence.

## 2026-06-05 -- Start must verify WSL systemd and GitHub online

Official GitHub Actions runner service guidance confirms that Linux runners
managed by systemd should be controlled through the service layer and checked
with service status after start. RunnerMonitor now treats `start` for
service-managed runners as a full readiness operation, not a fire-and-forget
command: systemd-backed runners are enabled before start so a previously
`disabled` WSL unit is corrected, the local unit must become `active`, and the
matching GitHub runner must report `online` before the command reports success.

The bug case was a WSL unit that existed but was `disabled` and `inactive
(dead)`. Plain `systemctl start` could leave the user with a success-looking
message even though GitHub still showed the runner unavailable. Polling the
GitHub runner API after local activation gives the app-level readiness signal:
the runner is online and listening for jobs.

## 2026-06-05 -- Release packaging PowerShell wildcard gotcha

RunnerMonitor release ZIP packaging should use `Compress-Archive -Path
(Join-Path $Stage '*')`, not `-LiteralPath` with a wildcard. PowerShell treats
`-LiteralPath` literally, so `stage\*` is not expanded and the ZIP is not
created. The v0.3.3 local release artifact was built by staging
`runner-monitor.ps1`, `bin/runner-monitor.exe`, sanitized
`bin/runner-monitor.json`, README files, and `LICENSE`, then compressing the
expanded stage contents and writing a SHA256 file.
