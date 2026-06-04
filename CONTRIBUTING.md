# Contributing

Thank you for improving RunnerMonitor.

RunnerMonitor controls local and remote CI runners, so contributions must keep
operator safety ahead of convenience.

## Development Setup

```powershell
git clone https://github.com/SGribanov/RunnerMonitor.git
cd RunnerMonitor
go test ./...
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
```

## Working Rules

- Do not commit secrets, tokens, passwords, or local `runner-monitor.json`.
- Keep destructive runner operations dry-run by default.
- Require explicit confirmation for unregistering runners or deleting folders.
- Preserve busy-runner checks unless a command is explicitly marked as forceful.
- Keep Windows, WSL, Linux, and future remote-host behavior documented.
- Update `README.md` and `README_RU.MD` together when user-facing behavior
  changes.

## Validation

Before opening a pull request, run:

```powershell
go test ./...
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
.\runner-monitor.ps1 --show-config
```

For TUI changes, also run a manual smoke test:

```powershell
.\runner-monitor.ps1
```

Confirm that the loading screen, table, commands, and exit flow work in a
normal terminal.

## Pull Requests

Use the pull request template and include:

- what changed;
- how it was validated;
- safety impact;
- any manual steps needed by the operator.
