<!-- Issue: SGribanov/RunnerMonitor#6 -->
# Safe Runner Cleanup Status

## Current State

- Added backend cleanup logic for local Windows paths and WSL runner paths.
- Added TUI commands:
  - `clear N`
  - `clear idle`
  - `auto-clear on`
  - `auto-clear off`
- Added CLI commands:
  - `--clear-current`
  - `--clear-repo owner/repo`
  - `--clear-idle`
- Cleanup refuses busy runners, stops/restarts running runners, and preserves
  runner registration files.
- Added tests for busy-runner refusal, work/archive cleanup, and repo filtering.

## Validation

- `go test ./...` passes.
- `runner-monitor.ps1 --once` prints the current inventory successfully.

## Notes

- Windows service cleanup still requires elevated PowerShell when the selected
  runner is service-managed.
- WSL systemd cleanup requires the current user to be able to stop/start the
  unit; otherwise cleanup fails before deleting `_work`.
