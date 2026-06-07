<!-- Issue: SGribanov/RunnerMonitor#34 -->
# Plans

## Source
- Task: Monitor GitHub-hosted Actions runners and publish a release.
- Canonical input: User request and GitHub issue #34.
- Repo context: Go CLI/TUI in `internal/app`, existing GitHub status loading via `gh api`.
- Last updated: 2026-06-07

## Assumptions
- GitHub-hosted runners are monitored as read-only workflow/job execution rows because they cannot be started, stopped, cleaned, or reprovisioned from RunnerMonitor.
- Repository-level workflow jobs are enough for the first hosted monitoring release; org/enterprise larger-runner fleet inventory can be added later.
- Existing dirty working-tree changes are preserved and incorporated only where they overlap with this task.

## Milestone Order
| ID | Title | Depends on | Status |
| --- | --- | --- | --- |
| M1 | Hosted job discovery and model integration | - | [x] |
| M2 | TUI/CLI safety and tests | M1 | [x] |
| M3 | Docs, research, binary, and release | M2 | [~] |

## M1. Hosted job discovery and model integration `[x]`
### Goal
- Add GitHub-hosted queued/in-progress workflow job monitoring alongside existing self-hosted runner rows.

### Tasks
- [x] Query workflow runs/jobs through GitHub REST API using the existing `gh api` wrapper.
- [x] Classify jobs that do not target `self-hosted` labels as GitHub-hosted.
- [x] Convert hosted jobs into read-only `Runner` rows with repo, workflow/job name, status, busy, labels, queue age, and URL/path context.
- [x] Keep existing self-hosted runner status loading unchanged.

### Definition of Done
- Refresh shows hosted execution rows even when no local runner folder exists.
- Hosted rows are visually distinguishable from local/self-hosted rows.

### Validation
```sh
go test ./internal/app
```

### Known Risks
- GitHub job API pagination and run volume can be noisy; first release should cap requests conservatively.

### Stop-and-Fix Rule
- If hosted discovery slows normal refresh noticeably or breaks self-hosted rows, fix before continuing.

## M2. TUI/CLI safety and tests `[x]`
### Goal
- Make hosted rows safely read-only while keeping normal self-hosted lifecycle behavior.

### Tasks
- [x] Block lifecycle, cleanup, remove, delete, and logs commands for hosted rows with clear messages.
- [x] Add unit tests for hosted job parsing/classification and read-only command behavior.
- [x] Verify existing lifecycle tests still pass.

### Definition of Done
- Hosted rows cannot trigger local filesystem or service operations.
- Regression tests cover the new behavior.

### Validation
```sh
go test ./...
```

### Known Risks
- The existing command handlers assume table rows are local runners.

### Stop-and-Fix Rule
- If any command can mutate a hosted row, stop and add guard coverage before release work.

## M3. Docs, research, binary, and release `[~]`
### Goal
- Ship the feature with updated docs, release notes, and a fresh binary.

### Tasks
- [x] Update README, README_RU, CHANGELOG, and technology insights in repo and IdeaBox vault.
- [x] Run release build script and confirm binary output.
- [ ] Commit, push, and publish a GitHub release with the binary artifact.
- [ ] Publish GitHub issue handoff and close only after acceptance gates pass.

### Definition of Done
- Documentation explains hosted monitoring limits and setup.
- Release artifact exists and tests/build pass.
- Issue #34 has a final handoff comment.

### Validation
```sh
go test ./...
.\scripts\build.ps1
```

### Known Risks
- Release publishing may fail if GitHub auth lacks release permissions.

### Stop-and-Fix Rule
- If release creation fails due to auth or tag conflict, report the exact blocker and leave the code/build validated.
