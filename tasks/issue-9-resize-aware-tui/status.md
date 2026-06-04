<!-- Issue: SGribanov/RunnerMonitor#9 -->
# Resize-Aware TUI Status

## Current Phase

Milestones 1-3 implemented; validation passed and handoff is ready.

## Done

- Created GitHub issue #9 for resize-aware TUI.
- Created branch `codex/9-resize-aware-tui`.
- Refreshed the TUI decision against the locally installed Charmbracelet
  packages: `bubbles/table` provides bounded table rendering and selection;
  Bubble Tea emits `tea.WindowSizeMsg` for resize handling.
- Replaced the interactive fixed-width table with a `bubbles/table` component.
- Added terminal width/height tracking and table/input resizing through
  `tea.WindowSizeMsg`.
- Added adaptive column sizing: wide terminals expand useful columns, narrow
  terminals hide low-priority columns and keep the table bounded.
- Added selected-row commands while preserving numeric commands.
- Updated README and research insights in the repo and IdeaBox vault.
- Passed final build, `--once` smoke, and TUI PTY smoke.

## In Progress

- None.

## Next

- Commit and push branch `codex/9-resize-aware-tui`.
- Publish GitHub issue handoff.

## Decisions

- Keep Bubble Tea and the Charmbracelet stack.
- Do not replace the whole TUI framework with `tview`.
- Keep text `RenderInventory` for non-interactive `--once` output.

## Validation Log

- Passed: `go test ./...`.
- Passed: `scripts\build.ps1`.
- Passed: `runner-monitor.ps1 --once`.
- Passed: manual PTY TUI smoke: loading screen, table render, compact help,
  single input prompt, and `q` exit.
