## What's New

### All-Runner Lifecycle Commands

RunnerMonitor TUI now accepts `start all` and `stop all` to run lifecycle
control for every controllable runner in the current inventory. The commands
reuse the same guarded behavior as single-runner `start` and `stop`.

Key points:
- Read-only GitHub-hosted and remote-only rows are skipped.
- Runners without a supported local control path are reported as skipped.
- After the command completes, the TUI immediately refreshes status so the
  table catches up with local service and GitHub state.

### Documentation

The in-app help panel and English/Russian README command references now include
the new all-runner commands.

### Validation

- `go test ./internal/app`
- `go test ./...`
- `git diff --check`
- `powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1`

### Assets

- `RunnerMonitor-v0.6.0-windows-x64.zip`
- `RunnerMonitor-v0.6.0-windows-x64.zip.sha256`

**Full Changelog**: https://github.com/SGribanov/RunnerMonitor/compare/v0.5.1...v0.6.0
