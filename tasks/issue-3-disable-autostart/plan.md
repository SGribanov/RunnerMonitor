<!-- Issue: SGribanov/RunnerMonitor#3 -->
# Disable Runner Autostart Plan

## Milestone 1: Project-scoped commands

Status: [x]

Goal: let Codex/operator start runners for a repository with one command.

Tasks:
- [x] Add `--start-repo owner/repo`.
- [x] Add `--stop-repo owner/repo`.
- [x] Add `--restart-repo owner/repo`.
- [x] Add `--start-current`, `--stop-current`, and `--restart-current`.
- [x] Add PowerShell wrapper that builds the binary on first use.
- [x] Add parent `C:\Repos\AGENTS.md` policy for Codex startup.
- [x] Skip manual runners clearly.

Validation:
- `go test ./...`
- `go run ./cmd/runner-monitor --start-repo SGribanov/NoSuchRepo`
- `powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --start-current`

## Milestone 2: Disable boot autostart

Status: [~]

Goal: make service-managed runners manual-start only.

Tasks:
- [x] Add `--disable-autostart`.
- [x] Attempt current disable operation.
- [ ] Re-run from elevated Windows PowerShell and root-authorized WSL session.

Required elevated commands:

```powershell
sc.exe config actions.runner.SGribanov-DeltaG.deltag-win start= demand
sc.exe config actions.runner.SGribanov-IdeaBox.ideabox-runner start= demand
wsl.exe sudo systemctl disable actions.runner.SGribanov-DeltaG.deltag-linux-wsl.service
wsl.exe sudo systemctl disable actions.runner.SGribanov-NewGenOsEngine.newgen-wsl-linux.service
```

## WSL background rule

Use `systemd` units for WSL runners. Do not rely on a visible terminal window.
`runner-monitor --start-repo ...` starts the unit in the background.
