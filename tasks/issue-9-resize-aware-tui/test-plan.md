<!-- Issue: SGribanov/RunnerMonitor#9 -->
# Resize-Aware TUI Test Plan

## Unit Tests

- Layout helpers:
  - narrow terminal keeps all visible columns at positive widths;
  - wide terminal allocates extra room to project/runner/path columns;
  - table height remains at least a small usable value.
- Runner rows:
  - project name is shown from the repository field;
  - queue count includes stale queue count when present;
  - busy state and labels are rendered as strings.
- Command parsing:
  - existing numeric commands still work;
  - selected-row commands work when no number is supplied.

## CLI Smoke

- `go test ./...`
- `scripts\build.ps1`
- `runner-monitor.ps1 --once`

## Manual TUI Smoke

- Start `runner-monitor.ps1`.
- Resize terminal from narrow to wide and back.
- Confirm table remains aligned and command input stays visible.
- Use arrow keys to select a runner.
- Run at least one non-destructive command such as `logs N` or `refresh`.

## Acceptance Gates

- Existing runner lifecycle commands are not regressed.
- Terminal resize does not break row alignment.
- Long project/runner/path values truncate inside their columns.
