<!-- Issue: SGribanov/RunnerMonitor#28 -->
# Test Plan

## Source
- Task: Ensure `start` brings disabled/inactive WSL runners fully online.
- Plan file: `tasks/issue-28-ensure-start-online/plan.md`
- Status file: `tasks/issue-28-ensure-start-online/status.md`
- Repo context: lifecycle commands and GitHub runner status checks.
- Last updated: 2026-06-05

## Validation Scope
- In scope: `RunLifecycle("start", runner)` for WSL/systemd service-managed runners, local service state checks, GitHub online verification.
- Out of scope: Manual Windows runner spawning and full live GitHub Actions job dispatch.

## Environment / Fixtures
- Data fixtures: Unit-test stubs for service control, local service state, GitHub runner status, and sleep.
- External dependencies: `gh` CLI only for live/manual verification; unit tests stub it.
- Setup assumptions: Existing tests run with `go test`.

## Test Levels

### Unit
- Start calls `enable` then `start` for WSL systemd runners.
- Start waits until local state becomes `active`.
- Start waits until GitHub status becomes `online`.
- Timeout path reports GitHub not online.

### Integration
- `go test ./...` covers package integration without real service mutation.

### End-to-End / Smoke
- Manual smoke, when a WSL runner is available: disable/stop its unit, run `--start-current`, confirm GitHub shows `online`.

## Negative / Edge Cases
- GitHub status remains `offline` or `unknown`.
- Service start succeeds but local state never becomes `active`.
- GitHub status lookup returns an error.

## Acceptance Gates
- [x] `go test ./internal/app`
- [x] `go test ./...`
- [x] Manual command output is truthful when online verification fails through unit-level regression coverage.

## Release / Demo Readiness
- [x] Core start scenario is covered through service/GitHub stubs.
- [x] Primary regression checks are green.
- [x] No blocker-level known issue remains.

## Command Matrix
```sh
go test ./internal/app
go test ./...
```

## Open Risks
- Unit tests do not prove a live WSL systemd instance transitions online; that remains a manual smoke check.

## Deferred Coverage
- Automated live WSL/systemd integration test is deferred because it would mutate the host runner service.
