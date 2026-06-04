<!-- Issue: SGribanov/RunnerMonitor#6 -->
# Safe Runner Cleanup Test Plan

## Automated

- `go test ./...`
- Verify busy runners are not cleaned.
- Verify `_work` contents and `actions-runner*.zip` archives are removed.
- Verify `.runner` is preserved.
- Verify repo-scoped cleanup ignores other repositories.

## Manual

- Start TUI and run `clear N` on an idle manual runner.
- Start TUI and run `clear idle`.
- Toggle `auto-clear on`, run `refresh`, and verify only idle runners are
  cleaned.
- For Windows service runners, run from elevated PowerShell and verify the
  service returns to running state.
- For WSL systemd runners, verify the unit returns to active state.
