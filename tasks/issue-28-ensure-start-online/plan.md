<!-- Issue: SGribanov/RunnerMonitor#28 -->
# Plans

## Source
- Task: Ensure `start` brings disabled/inactive WSL runners fully online.
- Canonical input: GitHub issue #28 and user report: WSL systemd unit existed but was `disabled` and `inactive (dead)`.
- Repo context: `internal/app/lifecycle.go`, WSL systemd controls, GitHub runner status polling.
- Last updated: 2026-06-05

## Assumptions
- A service-managed WSL runner should be enabled before start so a previously disabled unit is corrected.
- `online` in GitHub runner API is the app-level signal that the runner is listening for jobs.
- The CLI/TUI can wait briefly during a direct lifecycle command to report a truthful result.

## Milestone Order
| ID | Title | Depends on | Status |
| --- | --- | --- | --- |
| M1 | Start path enforces enable/start/verify | - | [x] |
| M2 | Regression tests and docs/research note | M1 | [x] |

## M1. Start path enforces enable/start/verify `[x]`
### Goal
- `start` no longer reports success for a service-managed runner until the service is active locally and online on GitHub.

### Tasks
- [x] Enable systemd autostart before starting WSL/Linux systemd runner services.
- [x] Wait for local service active state after start.
- [x] Poll GitHub runner status until the named runner is `online`.
- [x] Keep stop/restart busy protection behavior unchanged.

### Definition of Done
- Disabled + inactive WSL systemd units are enabled and started.
- If GitHub does not report `online` in time, the result says so instead of claiming success.

### Validation
```sh
go test ./internal/app
go test ./...
```

### Known Risks
- GitHub status polling depends on `gh` authentication and network availability.

### Stop-and-Fix Rule
- If unit tests or full Go tests fail, fix the failure before moving on.

## M2. Regression tests and docs/research note `[x]`
### Goal
- Capture the behavior and validation trail for future sessions.

### Tasks
- [x] Add unit tests for WSL disabled/inactive start behavior.
- [x] Add negative test for GitHub never reaching `online`.
- [x] Update technology insights in repo and IdeaBox vault.

### Definition of Done
- Tests fail on the old behavior and pass on the new behavior.
- Research note records the systemd + GitHub online verification finding.

### Validation
```sh
go test ./...
```

### Known Risks
- Existing uncommitted issue #27 changes are present and should be preserved.

### Stop-and-Fix Rule
- If validation fails from unrelated dirty work, isolate and report it without reverting user changes.
