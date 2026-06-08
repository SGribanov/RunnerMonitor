<!-- Issue: SGribanov/RunnerMonitor#37 -->
# Test Plan

## Source
- Task: Allow unregistering remote-only self-hosted runners.
- Plan file: `tasks/issue-37-remote-runner-unregister/plan.md`
- Status file: `tasks/issue-37-remote-runner-unregister/status.md`
- Last updated: 2026-06-08

## Validation Scope
- In scope: GitHub runner ID parsing, remote-only dry-run, confirmed unregister endpoint, missing-ID guard, and GitHub-hosted read-only safety.
- Out of scope: live deletion of the user's `deltag-win-2` registration without explicit user command.

## Environment / Fixtures
- Unit fixtures use synthetic `Runner` values and stubbed `ghAPIMethod`.
- Live smoke requires `gh` auth with repository admin permission.

## Test Levels

### Unit
- `github-remote` dry-run includes GitHub API removal and no local action.
- `github-remote` confirmed removal calls `DELETE /repos/{owner}/{repo}/actions/runners/{runner_id}`.
- `github-remote` confirmed removal fails clearly when runner ID is missing.
- Remote-only rows preserve GitHub runner ID from status loading.
- `github-hosted` rows remain read-only.

### Integration
- `go test ./...` verifies all packages.

### Smoke
- Rebuild local binary with `scripts/build.ps1`.
- Run `bin\runner-monitor.exe --once` or equivalent to confirm startup and inventory still work.

## Negative / Edge Cases
- Remote-only runner is busy and removal is requested without force.
- Remote-only row has no GitHub ID.
- `delete ... confirm` is used on remote-only row with no local folder.
- GitHub API returns auth or permission failure.

## Acceptance Gates
- [x] `go test ./internal/app`
- [x] `go test ./...`
- [x] `git diff --check`
- [x] `.\scripts\build.ps1`

## Smoke Result
- [x] `bin\runner-monitor.exe --once` showed `deltag-win-2` as `github` / `remote` / `offline` with path `(not local)`.
- [x] `bin\runner-monitor.exe --remove-runner deltag-win-2 --repo SGribanov/DeltaG` produced a dry-run GitHub API removal plan with runner id `26`.

## Command Matrix
```sh
go test ./internal/app
go test ./...
git diff --check
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
```
