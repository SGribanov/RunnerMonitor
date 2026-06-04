<!-- Issue: SGribanov/RunnerMonitor#9 -->
# Resize-Aware TUI Plan

## Scope

Improve the interactive RunnerMonitor TUI so terminal resizing keeps the runner
table readable and command input visible.

## Milestone 1 -- TUI Layout Model `[x]`

Goal: make the Bubble Tea model track terminal dimensions and own a table
component.

Tasks:
- [x] Add window width/height fields to the model.
- [x] Replace the interactive hand-rendered table with `bubbles/table`.
- [x] Recompute table width, height, and columns on `tea.WindowSizeMsg`.

Validation:
- `go test ./...`

## Milestone 2 -- Runner Selection And Commands `[x]`

Goal: keep existing numeric commands while adding keyboard selection.

Tasks:
- [x] Route arrow/page keys to the table when command input is empty.
- [x] Let commands that need an index use the selected runner when no number is
  supplied.
- [x] Keep existing commands such as `start 1`, `stop 1`, `clear 1`, `logs 1`,
  and `connect remote NAME`.

Validation:
- `go test ./...`
- Manual TUI smoke test.

## Milestone 3 -- Documentation And Handoff `[x]`

Goal: document the layout decision and publish issue handoff.

Tasks:
- [x] Update README if command behavior changes.
- [x] Update research insights in repo and IdeaBox vault.
- [x] Publish GitHub issue handoff.

Validation:
- `scripts\build.ps1`
- Manual resize smoke test in a terminal.
