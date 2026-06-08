<!-- Issue: SGribanov/RunnerMonitor#37 -->
# Plans

## Source
- Task: Allow unregistering remote-only self-hosted runners after their local folder was deleted.
- Concrete report: `deltag-win-2` remains listed and cannot be unbound.
- Repo context: Go CLI/TUI in `internal/app`; remote-only rows are built from GitHub runner API data.
- Last updated: 2026-06-08

## Assumptions
- The row is a repository-level self-hosted runner returned by GitHub, not a GitHub-hosted workflow job.
- The operator has `gh` auth with repository admin permission for the target repo.
- Deleting a remote-only row should unregister GitHub registration only and must not touch local files.

## Milestone Order
| ID | Title | Depends on | Status |
| --- | --- | --- | --- |
| M1 | Preserve GitHub runner ID | - | [x] |
| M2 | Remote-only unregister path | M1 | [x] |
| M3 | Docs, validation, and handoff | M2 | [x] |

## M1. Preserve GitHub runner ID `[x]`
### Goal
- Carry `runner_id` from GitHub API into remote-only inventory rows.

### Tasks
- [x] Parse `id` from `repos/{owner}/{repo}/actions/runners`.
- [x] Store it in `GitHubRunnerStatus` and `Runner`.
- [x] Preserve existing local/remote matching behavior.

### Definition of Done
- A remote-only row has enough data to call repository runner delete API.

### Validation
```sh
go test ./internal/app
```

### Known Risks
- Runner IDs are required for deletion; missing IDs must fail clearly.

### Stop-and-Fix Rule
- If a confirmed remote-only removal can run without a runner ID, stop and add a guard.

## M2. Remote-only unregister path `[x]`
### Goal
- Let confirmed `remove`/`delete` unregister remote-only self-hosted runners through GitHub API.

### Tasks
- [x] Keep GitHub-hosted rows read-only.
- [x] Add remote-only dry-run explaining GitHub API deletion and no local actions.
- [x] Call `DELETE /repos/{owner}/{repo}/actions/runners/{runner_id}` only after confirmation.
- [x] Add regression tests with stubbed `gh api`.

### Definition of Done
- `remove N confirm` and `delete N confirm` work for `github-remote` rows without local folder access.

### Validation
```sh
go test ./internal/app
go test ./...
```

### Known Risks
- Busy remote-only runners should still require `--force` to avoid removing active capacity by accident.

### Stop-and-Fix Rule
- If GitHub-hosted rows become removable, restore the read-only guard before continuing.

## M3. Docs, validation, and handoff `[x]`
### Goal
- Document behavior, update project knowledge, validate, and publish handoff.

### Tasks
- [x] Update README, README_RU, CHANGELOG, and technology insights.
- [x] Sync technology insights to IdeaBox vault.
- [x] Build local binary.
- [x] Run validation gates.
- [x] Publish GitHub issue handoff comment.

### Definition of Done
- User can run the local binary and unregister `deltag-win-2` without recreating its folder.

### Validation
```sh
go test ./...
git diff --check
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
```

### Known Risks
- Live removal depends on `gh` auth and repository admin permissions.

### Stop-and-Fix Rule
- If build or tests fail, fix before handoff.
