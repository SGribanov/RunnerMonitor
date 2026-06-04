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
- WSL systemd cleanup now falls back to sudo by reading the password file from
  `RUNNER_MONITOR_WSL_SUDO_FILE` or `C:\Users\gsv777\Desktop\WSL_sudo.txt`.
- Windows service cleanup reports a UTF-8 elevated-PowerShell hint instead of
  mojibake when the TUI is not elevated.
- Added tests for busy-runner refusal, work/archive cleanup, and repo filtering.

## Validation

- `go test ./...` passes.
- `runner-monitor.ps1 --once` prints the current inventory successfully.
- `runner-monitor.ps1 --clear-repo SGribanov/MyCloneOsEngine` clears both the
  Windows manual runner and WSL systemd runner successfully.
- Non-elevated `runner-monitor.ps1 --clear-repo SGribanov/IdeaBox` fails safely
  before cleanup and reports that service control may require elevated
  PowerShell.

## Notes

- Windows service cleanup still requires elevated PowerShell when the selected
  runner is service-managed.
- WSL systemd cleanup requires a readable sudo password file if plain
  `systemctl` needs interactive authentication.
