<!-- Issue: SGribanov/RunnerMonitor#35 -->
# Test Plan

## Source
- Task: Show remote-only self-hosted runners.
- Plan file: `tasks/issue-35-remote-only-runners/plan.md`
- Status file: `tasks/issue-35-remote-only-runners/status.md`
- Last updated: 2026-06-07

## Validation Scope
- In scope: refresh integration, local/remote matching, read-only command safety, table/render behavior.
- Out of scope: org/enterprise runner group inventory beyond repositories already monitored.

## Environment / Fixtures
- Unit fixtures: synthetic local runners and GitHub runner status maps.
- External dependency: `gh` CLI for live refresh, not required for unit tests.

## Test Levels

### Unit
- Existing local runner is enriched from matching GitHub status and not duplicated.
- Unmatched GitHub runner becomes a `github-remote` read-only row.
- Remote-only row carries status, busy flag, labels, OS, and version.
- Lifecycle, cleanup, logs, remove, and audit behavior does not perform local actions for remote-only rows.

### Integration
- `go test ./...` verifies all packages.

### Smoke
- `runner-monitor --once` can list remote-only rows when GitHub API returns them for monitored repos.

## Negative / Edge Cases
- GitHub API returns zero runners.
- GitHub runner name differs only by case from local runner name.
- Remote-only runner is busy.
- Remote-only runner has no labels or version.

## Acceptance Gates
- [x] `go test ./internal/app`
- [x] `go test ./...`
- [x] `git diff --check`
- [x] `.\scripts\build.ps1`
- [x] `.\bin\runner-monitor.exe --once`

## Smoke Result
- [x] `go run ./cmd/runner-monitor --once` showed DeltaG remote-only self-hosted rows as `Host=github`, `Local=remote`, and `Path=(not local)`.
- [x] v0.5.0 ZIP was built and extracted successfully with `runner-monitor.exe`, wrapper, README files, license, and sanitized config.

## Command Matrix
```sh
go test ./internal/app
go test ./...
git diff --check
```
