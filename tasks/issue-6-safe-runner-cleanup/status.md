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
  - `--clear-runner NAME`
  - `--clear-idle`
- Cleanup refuses busy runners, stops/restarts running runners, and preserves
  runner registration files.
- WSL systemd cleanup now falls back to sudo by reading the password file from
  `RUNNER_MONITOR_WSL_SUDO_FILE` or `C:\Users\gsv777\Desktop\WSL_sudo.txt`.
- Windows service cleanup launched from a non-elevated TUI now requests an
  elevated PowerShell helper through UAC for the selected runner.
- Added tests for busy-runner refusal, work/archive cleanup, and repo filtering.

## Validation

- `go test ./...` passes.
- `runner-monitor.ps1 --once` prints the current inventory successfully.
- `runner-monitor.ps1 --clear-repo SGribanov/MyCloneOsEngine` clears both the
  Windows manual runner and WSL systemd runner successfully.
- `runner-monitor.ps1 --help` includes `--clear-runner`.
- `runner-monitor.ps1 --clear-runner definitely-not-a-runner` reports a clean
  not-found message.

## Notes

- Windows service cleanup still requires elevated PowerShell when the selected
  runner is service-managed, but the non-elevated TUI now opens a UAC helper
  instead of deleting anything without stopping the service.
- WSL systemd cleanup requires a readable sudo password file if plain
  `systemctl` needs interactive authentication.
