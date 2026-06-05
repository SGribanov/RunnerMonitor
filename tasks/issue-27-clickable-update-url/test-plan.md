<!-- Issue: SGribanov/RunnerMonitor#27 -->
# Test Plan

## Source
- Task: Make update notice release URL clickable.
- Plan file: `tasks/issue-27-clickable-update-url/plan.md`
- Status file: `tasks/issue-27-clickable-update-url/status.md`
- Repo context: TUI update notice rendering.
- Last updated: 2026-06-05

## Validation Scope
- In scope: update notice rendering, URL readability, OSC-8 escape sequence placement.
- Out of scope: terminal-specific click behavior.

## Environment / Fixtures
- Data fixtures: Update notice containing a GitHub release URL.
- External dependencies: None for unit tests.
- Setup assumptions: Existing Go tests are sufficient for regression coverage.

## Test Levels

### Unit
- Verify OSC-8 opener is present for the release URL.
- Verify the visible URL remains in the rendered view.

### Integration
- `go test ./...` covers the TUI model with existing package tests.

### End-to-End / Smoke
- Optional manual TUI smoke in an OSC-8-capable terminal.

## Negative / Edge Cases
- Notices without URL parentheses fall back to normal truncation.
- Non-HTTP strings inside parentheses are not treated as links.

## Acceptance Gates
- [x] `go test ./...`

## Release / Demo Readiness
- [x] Primary regression checks are green.
- [x] No blocker-level known issue remains.

## Command Matrix
```sh
go test ./...
```

## Open Risks
- Unit tests verify escape output, not actual terminal click handling.

## Deferred Coverage
- Automated terminal-emulator click validation is deferred.
