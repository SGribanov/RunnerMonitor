# RunnerMonitor

[Русская версия](README_RU.MD)

RunnerMonitor is a lightweight terminal UI and command-line tool for monitoring
and controlling self-hosted CI runners. It was built for a workstation that has
multiple GitHub Actions runners across Windows and WSL/Linux, with a planned
move to a dedicated runner host on the local network.

The tool combines local runner lifecycle state with GitHub runner state, queue
information, repository ownership, and safe lifecycle commands.

## What It Does

- Discovers GitHub Actions runner folders from configured Windows, WSL, and
  Linux runner roots.
- Shows the project/repository each runner belongs to.
- Merges local service/process state with GitHub status from `gh api`.
- Displays busy state, queued workflow count, and stale queued workflow count.
- Starts, stops, restarts, clears, removes, and reprovisions selected runners.
- Keeps destructive operations guarded by dry-runs, busy-runner checks, and
  explicit confirmation.
- Provides project-scoped commands that Codex or an operator can run before a
  push or CI wait.
- Supports saved SSH remote-runner host profiles for the future dedicated
  runner machine.
- Keeps OneDev support as a future provider-oriented direction.

## Stack

- Go 1.26+
- Charmbracelet Bubble Tea for the TUI
- Charmbracelet Bubbles for table and text input widgets
- Charmbracelet Lip Gloss for terminal styling
- GitHub CLI (`gh`) for GitHub Actions runner and workflow data
- Windows PowerShell for Windows service discovery/control
- WSL and Linux systemd for Linux runner discovery/control
- SSH for remote runner-host access

## Requirements

- Go 1.26+
- GitHub CLI authenticated with access to monitored repositories:

  ```powershell
  gh auth status
  ```

- Windows PowerShell for local Windows service discovery.
- WSL for current-workstation Linux runner discovery.
- `git` for current-project repository detection.

## Quick Start

Build the executable and create the app-local config:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
.\runner-monitor.ps1 --show-config
```

Open the TUI:

```powershell
.\runner-monitor.ps1
```

Print inventory once:

```powershell
.\runner-monitor.ps1 --once
```

Audit runner cleanup candidates:

```powershell
.\runner-monitor.ps1 --audit
```

Start runners for the current GitHub project:

```powershell
.\runner-monitor.ps1 --start-current
```

## Configuration

RunnerMonitor reads host-specific settings from `runner-monitor.json` next to
the executable. For the default local build, the config path is:

```text
C:\Repos\RunnerMonitor\bin\runner-monitor.json
```

Create or inspect it:

```powershell
.\runner-monitor.ps1 --init-config
.\runner-monitor.ps1 --show-config
```

`--show-config` masks `wslSudoPassword` as `<set>` or `<empty>` and never prints
the actual value.

Default config:

```json
{
  "projectsRoot": "C:\\Repos",
  "windowsRunnerRoots": [
    "C:\\Runners"
  ],
  "wslRunnerRoots": [
    "/home/gsv777/Runners"
  ],
  "linuxRunnerRoots": [
    "/opt/Runners",
    "/srv/Runners"
  ],
  "wslSudoPassword": ""
}
```

Settings fields:

| Field | Purpose |
|---|---|
| `projectsRoot` | Root directory used by `--project`, for example `C:\Repos`. |
| `windowsRunnerRoots` | Windows runner root folders to scan and treat as safe for runner-folder deletion. |
| `wslRunnerRoots` | WSL runner root folders to scan. |
| `linuxRunnerRoots` | Linux runner root folders used on a native Linux runner host. |
| `wslSudoPassword` | Direct WSL sudo password value for fallback service control. Keep it only in the app-local config. |

For tests or special launch contexts, override the config path:

```powershell
$env:RUNNER_MONITOR_CONFIG = "D:\Temp\runner-monitor.json"
```

Do not commit a real `runner-monitor.json` or sudo password.

## TUI Usage

Inside the TUI:

| Command | Description |
|---|---|
| `refresh` | Refresh local and GitHub runner state. |
| Arrow keys | Select a runner row. |
| `start [N]` | Start runner `N`, or the selected runner if `N` is omitted. |
| `stop [N]` | Stop runner `N`, or the selected runner if `N` is omitted. |
| `restart [N]` | Restart runner `N`, or the selected runner if `N` is omitted. |
| `force-stop [N]` | Stop even when GitHub reports the runner as busy. Use carefully. |
| `force-restart [N]` | Restart even when GitHub reports the runner as busy. Use carefully. |
| `clear [N]` | Safely clear idle runner work files for runner `N` or selected row. |
| `clear idle` | Clear all idle runners. Busy runners are skipped. |
| `auto-clear on` | Run safe idle cleanup after refresh. |
| `auto-clear off` | Disable refresh-triggered auto cleanup. |
| `remove [N]` | Dry-run GitHub runner unregistration for runner `N` or selected row. |
| `remove [N] confirm` | Execute runner unregistration after confirmation. |
| `delete [N] confirm` | Unregister and delete a safe runner folder. |
| `logs [N]` | Open runner logs for runner `N` or selected row. |
| `connect remote NAME` | Open the saved remote RunnerMonitor TUI over SSH. |
| `q`, `quit`, `exit`, `Esc`, `Ctrl+C` | Exit the TUI. |

The table is resize-aware. On narrow terminals, low-priority columns are hidden
before the main project, runner, status, busy, and queue columns are allowed to
drift.

## CLI Command Reference

Inventory and audit:

```powershell
.\runner-monitor.ps1 --once
.\runner-monitor.ps1 --audit
```

Project-scoped lifecycle:

```powershell
.\runner-monitor.ps1 --start-repo SGribanov/DeltaG
.\runner-monitor.ps1 --stop-repo SGribanov/DeltaG
.\runner-monitor.ps1 --restart-repo SGribanov/DeltaG
```

Current Git project lifecycle:

```powershell
.\runner-monitor.ps1 --start-current
.\runner-monitor.ps1 --stop-current
.\runner-monitor.ps1 --restart-current
```

Safe cleanup:

```powershell
.\runner-monitor.ps1 --clear-repo SGribanov/DeltaG
.\runner-monitor.ps1 --clear-current
.\runner-monitor.ps1 --clear-idle
.\runner-monitor.ps1 --clear-runner ideabox-runner
```

Autostart policy:

```powershell
.\runner-monitor.ps1 --disable-autostart
```

Config:

```powershell
.\runner-monitor.ps1 --init-config
.\runner-monitor.ps1 --init-config --overwrite-config
.\runner-monitor.ps1 --show-config
```

Remote host profiles:

```powershell
.\runner-monitor.ps1 --configure-remote runnerbox
.\runner-monitor.ps1 --connect-remote runnerbox
```

Runner removal is dry-run by default:

```powershell
.\runner-monitor.ps1 --remove-runner ideabox-runner --repo SGribanov/IdeaBox
.\runner-monitor.ps1 --remove-runner ideabox-runner --repo SGribanov/IdeaBox --confirm
.\runner-monitor.ps1 --remove-runner ideabox-runner --repo SGribanov/IdeaBox --confirm --delete-folder
```

Runner reprovisioning configures an existing prepared runner distribution
folder. The `--project` value is a project folder name under configured
`projectsRoot`.

```powershell
.\runner-monitor.ps1 --add-runner runner-monitor-win --project RunnerMonitor --runner-folder C:\Runners\SGribanov-RunnerMonitor\runner-monitor-win --labels "self-hosted,Windows,X64"
.\runner-monitor.ps1 --add-runner runner-monitor-win --project RunnerMonitor --runner-folder C:\Runners\SGribanov-RunnerMonitor\runner-monitor-win --labels "self-hosted,Windows,X64" --confirm --replace
```

Install and start the runner service after configuration:

```powershell
.\runner-monitor.ps1 --add-runner runner-monitor-win --project RunnerMonitor --runner-folder C:\Runners\SGribanov-RunnerMonitor\runner-monitor-win --labels "self-hosted,Windows,X64" --confirm --replace --service
```

## Remote Runner Host

When runners move to a dedicated machine on the network, run RunnerMonitor on
that machine and connect to it over SSH.

Configure or update a remote host profile:

```powershell
.\runner-monitor.ps1 --configure-remote runnerbox
```

The prompt asks for:

- remote name;
- SSH host or alias;
- host OS: `windows` or `linux`;
- remote RunnerMonitor path;
- default remote project path.

Open the saved remote TUI:

```powershell
.\runner-monitor.ps1 --connect-remote runnerbox
```

Equivalent Windows SSH command:

```powershell
ssh -t runnerbox "powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1"
```

Start runners for the current remote project before Codex pushes or waits for
CI:

```powershell
ssh runnerbox "cd C:/Repos/DeltaG; powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1 --start-current"
```

Linux remote host example:

```powershell
ssh -t runnerbox "cd /opt/RunnerMonitor && ./runner-monitor"
ssh runnerbox "cd /srv/DeltaG && /opt/RunnerMonitor/runner-monitor --start-current"
```

Saved remote profiles are stored separately from the app-local runner settings
under the current user's config directory as `RunnerMonitor\remote-hosts.json`.

## Runner Folder Layout

The preferred current layout is:

```text
C:\Runners\<owner>-<repo>\<runner-name>
/home/gsv777/Runners/<owner>-<repo>/<runner-name>
```

For a future dedicated Linux host, use a stable shared root such as:

```text
/opt/Runners/<owner>-<repo>/<runner-name>
/srv/Runners/<owner>-<repo>/<runner-name>
```

Runner folder migration is tracked separately and should be performed
runner-by-runner. Do not move or delete busy runners without explicit approval.

## Safety Model

- Busy runners are protected by default.
- Removal and reprovisioning are dry-run by default.
- Folder deletion requires `--delete-folder` and is limited to configured safe
  runner roots.
- Cleanup removes safe generated content such as `_work` contents and runner
  installer archives, while preserving runner registration and binaries.
- `--show-config` never prints the real WSL sudo password.
- Local app config files and secrets must not be committed.

## Build And Test

Run tests:

```powershell
go test ./...
```

Build:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
```

Generate icon assets:

```powershell
uv run --with pillow python .\scripts\generate-icon-assets.py
go run github.com/akavel/rsrc@v0.10.2 -arch amd64 -ico .\assets\runner-monitor-hourglass.ico -o .\cmd\runner-monitor\runner-monitor_windows_amd64.syso
```

## Repository Layout

```text
cmd/runner-monitor/     Application entry point
internal/app/           Discovery, GitHub integration, lifecycle, TUI, cleanup
scripts/                Build, migration, cleanup, and automation helpers
assets/                 Hourglass icon assets
reports/                Runner audit reports
research/               Long-lived project insights
tasks/                  Per-issue implementation plans and status files
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Contributions should preserve the safety
model: no unguarded destructive runner operations, no committed secrets, and no
silent changes to runner registrations. Please also follow the
[Code of Conduct](CODE_OF_CONDUCT.md).

## Security

See [SECURITY.md](SECURITY.md). Do not open public issues with secrets, runner
tokens, sudo passwords, or private hostnames.

## License

RunnerMonitor is licensed under the [MIT License](LICENSE).
