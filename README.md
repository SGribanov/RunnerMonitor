# RunnerMonitor

Lightweight TUI for monitoring and controlling self-hosted CI runners.

RunnerMonitor is designed around a simple topology: run the TUI on the operator
machine and control local, WSL, and future dedicated runner hosts through
standard service mechanisms such as Windows Services, Linux systemd, and SSH.

## Current scope

- Auto-discover GitHub Actions runner directories from local Windows and WSL.
- Merge local lifecycle state with GitHub runner status through `gh api`.
- Show queued and stale queued workflow counts per repository.
- Control service-managed runners with short commands such as `start 1`,
  `stop 1`, and `restart 1`.
- Keep runner folder migration and OneDev support as follow-up phases.

## Requirements

- Go 1.26+
- GitHub CLI authenticated with access to the monitored repositories
- Windows PowerShell for local Windows service discovery
- WSL for Linux runner discovery on the current workstation

## Run

```powershell
go run ./cmd/runner-monitor
go run ./cmd/runner-monitor --once
go run ./cmd/runner-monitor --audit
go run ./cmd/runner-monitor --start-repo SGribanov/DeltaG
go run ./cmd/runner-monitor --start-current
go run ./cmd/runner-monitor --disable-autostart
```

Inside the TUI:

- `refresh`
- `start 1`
- `stop 1`
- `restart 1`
- `force-stop 1`
- `force-restart 1`
- `logs 1`
- `q`

Codex/operator automation can start runners for a project without opening the
TUI:

```powershell
runner-monitor --start-repo SGribanov/DeltaG
runner-monitor --start-current
runner-monitor --stop-repo SGribanov/DeltaG
runner-monitor --restart-repo SGribanov/DeltaG
```

From any project root with a GitHub `origin`, Codex can run:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --start-current
```

Optional local git hook for a repository:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\scripts\install-prepush-hook.ps1 -RepoPath C:\Repos\DeltaG
```

Disable runner autostart from an elevated PowerShell session:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\scripts\disable-autostart-elevated.ps1
```
