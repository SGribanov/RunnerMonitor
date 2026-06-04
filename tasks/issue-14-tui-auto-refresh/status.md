<!-- Issue: SGribanov/RunnerMonitor#14 -->
# TUI Auto-Refresh Status

## Current Phase

Milestone 3 in progress.

## Done

- Created GitHub issue #14.
- Created branch `codex/14-tui-auto-refresh`.
- Refreshed the Bubble Tea timer decision through `search_ai_mcp_default`; the
  MCP search results were not relevant, so the implementation is anchored to
  local package docs for the installed Bubble Tea API.
- Added settings-backed auto-refresh scheduling in the Bubble Tea model with a
  5-second default.
- Added `refreshing` state so manual and automatic refreshes do not overlap.
- Kept existing TUI table visible during manual and automatic refreshes.
- Added unit tests for auto-refresh scheduling and manual refresh visibility.
- Updated README, README_RU, CHANGELOG, task docs, and technology insights in
  the repo and IdeaBox vault.
- Confirmed IdeaBox watcher is running with one supervised watcher process.

## In Progress

- Build, smoke, commit, tag, and publish release `v0.2.0`.

## Next

- Run `scripts\build.ps1`.
- Run `runner-monitor.ps1 --once`.
- Commit and publish release.

## Decisions

- Use Bubble Tea command scheduling instead of a goroutine-owned ticker so the
  update loop remains message-driven.
- Keep `auto-clear` tied to refresh completion; an enabled `auto-clear` applies
  after both manual and automatic refreshes.
- Manual `refresh` only enters the loading-only screen when no previous runner
  inventory exists.
- Invalid or omitted `tuiRefreshIntervalSeconds` values normalize to 5 seconds.

## Validation Log

- Passed: `go test ./...`.
- Passed: `scripts\build.ps1`.
- Passed: `.\runner-monitor.ps1 --show-config`; normalized config includes
  `"tuiRefreshIntervalSeconds": 5`.
- Passed: `.\runner-monitor.ps1 --once`.
- Passed: PTY TUI smoke; table stayed visible and an `auto-refreshed at ...`
  message appeared after the configured interval.
