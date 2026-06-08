## What's New

### Remote-Only Runner Unregister

RunnerMonitor can now unregister remote-only self-hosted runners through the
GitHub REST API after explicit confirmation. This fixes the case where a local
runner folder was deleted manually but the runner, such as `deltag-win-2`,
still remained registered and visible in GitHub.

Key points:
- `remove [N] confirm` and `delete [N] confirm` work for `github` / `remote`
  self-hosted runner rows.
- GitHub-hosted workflow rows remain read-only.
- Remote-only `delete` unregisters the GitHub runner only; no local folder or
  service action is attempted.

### TUI Terminal Restore Hardening

Interactive mutating TUI commands now run through Bubble Tea terminal
release/restore handling. This strengthens the existing frame-level text-mode
reset and helps preserve the terminal font/display state after PowerShell, cmd,
WSL, and service-control operations.

### Validation

- `go test ./internal/app`
- `go test ./...`
- `git diff --check`
- `powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1`
- `.\bin\runner-monitor.exe --once`
- Dry-run unregister check for `deltag-win-2`

### Assets

- `RunnerMonitor-v0.5.1-windows-x64.zip`
- `RunnerMonitor-v0.5.1-windows-x64.zip.sha256`

**Full Changelog**: https://github.com/SGribanov/RunnerMonitor/compare/v0.5.0...v0.5.1
