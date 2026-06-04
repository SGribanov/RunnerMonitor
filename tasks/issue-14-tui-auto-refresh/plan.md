<!-- Issue: SGribanov/RunnerMonitor#14 -->
# TUI Auto-Refresh Plan

## Scope

Refresh RunnerMonitor TUI inventory automatically every 5 seconds while keeping
manual refresh, startup loading, and cleanup behavior predictable.

## Milestone 1 -- Refresh Scheduler `[x]`

Goal: add a Bubble Tea timer that schedules non-overlapping inventory refreshes.

Tasks:
- [x] Add a settings-backed auto-refresh interval with a 5-second default.
- [x] Schedule recurring auto-refresh ticks while the TUI is open.
- [x] Avoid overlapping manual and automatic refresh commands.
- [x] Show a concise status message when auto-refresh updates the table.
- [x] Keep existing table data visible during manual and automatic refresh.

Validation:
- `go test ./...`

## Milestone 2 -- Documentation And Research `[x]`

Goal: reflect the auto-refresh behavior in user docs and durable project notes.

Tasks:
- [x] Update English and Russian README command/TUI/config sections.
- [x] Update changelog for the next release.
- [x] Update repo research insights and IdeaBox vault copy.

Validation:
- Review generated markdown diffs.

## Milestone 3 -- Release `[ ]`

Goal: build, tag, publish, and hand off the release.

Tasks:
- [ ] Run build and smoke checks.
- [ ] Commit using conventional commit format.
- [ ] Create and push `v0.2.0`.
- [ ] Publish GitHub release with binaries.
- [ ] Close issue #14 with handoff.

Validation:
- `scripts\build.ps1`
- `runner-monitor.ps1 --once`
- `gh release view v0.2.0`
