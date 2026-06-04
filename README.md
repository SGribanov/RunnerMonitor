# RunnerMonitor

Lightweight TUI for monitoring and controlling self-hosted CI runners.

RunnerMonitor is designed around a simple topology: run the TUI on the operator
machine and control local, WSL, and future dedicated runner hosts through
standard service mechanisms such as Windows Services, Linux systemd, and SSH.

## Current scope

- Auto-discover GitHub Actions runner directories from local Windows and WSL.
- Merge local lifecycle state with GitHub runner status through `gh api`.
- Show queued and stale queued workflow counts per repository.
- Show the project each runner belongs to in the TUI/audit tables.
- Control service-managed runners with short commands such as `start 1`,
  `stop 1`, and `restart 1`.
- Control manual Windows `run.cmd` runners in a hidden background process.
- Track runner folder migration into common `Runners` directories as a separate
  safety-gated phase.
- Keep OneDev support as a follow-up phase.

## Requirements

- Go 1.26+
- GitHub CLI authenticated with access to the monitored repositories
- Windows PowerShell for local Windows service discovery
- WSL for Linux runner discovery on the current workstation

## Settings

RunnerMonitor reads host-specific settings from `runner-monitor.json` next to
the executable:

```powershell
runner-monitor --init-config
runner-monitor --show-config
```

For the default local build this is:

```text
C:\Repos\RunnerMonitor\bin\runner-monitor.json
```

The build script creates this file beside `runner-monitor.exe` if it does not
already exist. The config can be moved or tested with `RUNNER_MONITOR_CONFIG`.
`--show-config` masks `wslSudoPassword`; it prints `<set>` or `<empty>`, never
the real value.

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

Use `wslSudoPassword` only in the app-local config file. Do not commit a real
password to the repository. Runner discovery, project resolution, safe folder
deletion, and WSL sudo fallback all use these settings.

## Run

```powershell
go run ./cmd/runner-monitor
go run ./cmd/runner-monitor --once
go run ./cmd/runner-monitor --audit
go run ./cmd/runner-monitor --start-repo SGribanov/DeltaG
go run ./cmd/runner-monitor --start-current
go run ./cmd/runner-monitor --disable-autostart
go run ./cmd/runner-monitor --configure-remote runnerbox
go run ./cmd/runner-monitor --connect-remote runnerbox
```

Inside the TUI:

- `refresh`
- arrow keys select a runner
- `start [N]`
- `stop [N]`
- `restart [N]`
- `force-stop [N]`
- `force-restart [N]`
- `clear [N]`
- `remove [N]`
- `remove [N] confirm`
- `delete [N] confirm`
- `logs [N]`
- `connect remote runnerbox`
- `q`

Commands with `[N]` can use either a runner number or the currently selected
runner row.

Codex/operator automation can start runners for a project without opening the
TUI:

```powershell
runner-monitor --start-repo SGribanov/DeltaG
runner-monitor --start-current
runner-monitor --stop-repo SGribanov/DeltaG
runner-monitor --restart-repo SGribanov/DeltaG
```

Safe cleanup keeps runner registration and binaries intact:

```powershell
runner-monitor --clear-repo SGribanov/DeltaG
runner-monitor --clear-current
runner-monitor --clear-idle
```

Removal and reprovisioning are dry-run by default. The project selector is the
folder name under configured `projectsRoot`. With the default config,
`RunnerMonitor` resolves `C:\Repos\RunnerMonitor` and reads its GitHub
`origin`.

```powershell
runner-monitor --remove-runner ideabox-runner --repo SGribanov/IdeaBox
runner-monitor --remove-runner ideabox-runner --repo SGribanov/IdeaBox --confirm
runner-monitor --remove-runner ideabox-runner --repo SGribanov/IdeaBox --confirm --delete-folder
```

Adding a runner configures an existing prepared runner distribution folder.
Quote labels in PowerShell so commas stay in one argument:

```powershell
runner-monitor --add-runner runner-monitor-win --project RunnerMonitor --runner-folder C:\Runners\SGribanov-RunnerMonitor\runner-monitor-win --labels "self-hosted,Windows,X64"
runner-monitor --add-runner runner-monitor-win --project RunnerMonitor --runner-folder C:\Runners\SGribanov-RunnerMonitor\runner-monitor-win --labels "self-hosted,Windows,X64" --confirm --replace
```

From any project root with a GitHub `origin`, Codex can run:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --start-current
```

## Remote Runner Host

When runners move to a dedicated machine on the network, connect to that machine
over SSH and run RunnerMonitor there. RunnerMonitor can prompt for host settings
and save them in the user config file, so the SSH command does not need to be
remembered every time.

Configure or update a remote host:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --configure-remote runnerbox
```

The prompt asks for:

- remote name
- SSH host or alias
- host OS: `windows` or `linux`
- remote RunnerMonitor path
- default remote project path

Open the saved remote TUI:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --connect-remote runnerbox
```

Inside the local TUI, use:

```text
connect remote runnerbox
```

The saved config lives under the current user's config directory as
`RunnerMonitor\remote-hosts.json`.

Under the hood, the saved Windows host command is equivalent to:

```powershell
ssh -t runnerbox "powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1"
```

Start runners for the current remote project before Codex pushes or waits for
CI:

```powershell
ssh runnerbox "cd C:/Repos/DeltaG; powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1 --start-current"
```

Start a specific repository on the remote host without relying on current
directory detection:

```powershell
ssh runnerbox "powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1 --start-repo SGribanov/DeltaG"
```

For a Linux runner host, use the same pattern with the compiled binary:

```powershell
ssh -t runnerbox "cd /opt/RunnerMonitor && ./runner-monitor"
ssh runnerbox "cd /srv/DeltaG && /opt/RunnerMonitor/runner-monitor --start-current"
```

Optional local git hook for a repository:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\scripts\install-prepush-hook.ps1 -RepoPath C:\Repos\DeltaG
```

Disable runner autostart from an elevated PowerShell session:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\scripts\disable-autostart-elevated.ps1
```

Runner cleanup proposals live in `reports/`; do not execute removal commands
without explicit runner-by-runner approval.

Manual audit decisions such as known keep runners live in
`runner-policy.json`.

## Runner Folder Migration

Runner folders should move into common `Runners` directories, but only
runner-by-runner after explicit approval. The target layout is:

```text
C:\Runners\<owner>-<repo>\<runner-name>
/home/gsv777/Runners/<owner>-<repo>/<runner-name>
```

The migration plan lives in `tasks/issue-5-runners-directory-migration/`.
Do not move busy runners; DeltaG Windows and WSL runners are currently deferred
while busy.

## Icon Assets

RunnerMonitor includes a generated hourglass icon set:

- `assets/runner-monitor-hourglass.png` -- transparent 512px source.
- `assets/runner-monitor-hourglass-spin.gif` -- animated spinning hourglass.
- `assets/runner-monitor-hourglass.ico` -- Windows executable icon.
- `cmd/runner-monitor/runner-monitor_windows_amd64.syso` -- Windows resource
  linked automatically by `go build`.

Regenerate the icon assets with:

```powershell
uv run --with pillow python scripts\generate-icon-assets.py
go run github.com/akavel/rsrc@v0.10.2 -arch amd64 -ico assets\runner-monitor-hourglass.ico -o cmd\runner-monitor\runner-monitor_windows_amd64.syso
```
