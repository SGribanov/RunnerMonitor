## What's New

### Remote Self-Hosted Runner Visibility

RunnerMonitor now shows self-hosted runners that are registered in monitored
GitHub repositories even when their runner folders are not local to the current
machine. These partner or dedicated-host runners appear as read-only
`github`/`remote` rows with path `(not local)`.

Key points:
- Shows GitHub status, busy state, labels, OS, version, and repo queue counts.
- Keeps locally discovered runners controllable with the existing lifecycle
  commands.
- Blocks lifecycle, cleanup, logs, remove, delete, and reprovision operations
  for remote-only rows because RunnerMonitor cannot control another machine's
  runner service or folder.

### Validation

- `go test ./internal/app`
- `go test ./...`
- `git diff --check`
- `.\bin\runner-monitor.exe --once`

### Assets

- `RunnerMonitor-v0.5.0-windows-x64.zip`
- `RunnerMonitor-v0.5.0-windows-x64.zip.sha256`

**Full Changelog**: https://github.com/SGribanov/RunnerMonitor/compare/v0.4.0...v0.5.0
