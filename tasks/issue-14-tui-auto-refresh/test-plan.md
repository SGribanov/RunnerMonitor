<!-- Issue: SGribanov/RunnerMonitor#14 -->
# TUI Auto-Refresh Test Plan

## Unit Tests

- Model initialization schedules startup refresh and spinner.
- Automatic refresh tick starts a refresh only when no refresh is already in
  progress.
- Settings default the TUI refresh interval to 5 seconds and load custom values.
- Manual refresh starts immediately and does not lose the next scheduled
  automatic refresh.
- Auto-refresh completion updates the table and leaves the model ready for the
  next tick.

## CLI Smoke

- `go test ./...`
- `scripts\build.ps1`
- `runner-monitor.ps1 --once`

## Manual TUI Smoke

- Start `runner-monitor.ps1`.
- Confirm startup loading appears, then table appears.
- Wait at least 6 seconds and confirm the status message updates without manual
  input.
- Run manual `refresh` and confirm it still works.
- Quit with `q`.

## Acceptance Gates

- TUI information refreshes every 5 seconds by default.
- TUI refresh interval can be changed through `runner-monitor.json`.
- Refresh calls do not overlap.
- Manual commands remain responsive.
- Release artifacts are attached to GitHub release `v0.2.0`.
