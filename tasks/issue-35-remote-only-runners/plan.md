<!-- Issue: SGribanov/RunnerMonitor#35 -->
# Plans

## Source
- Task: Show partner/remote self-hosted runners registered in monitored GitHub repositories.
- Canonical input: User request and GitHub issue #35.
- Repo context: Go CLI/TUI in `internal/app`, refresh combines local discovery with GitHub runner status from `gh api`.
- Last updated: 2026-06-07

## Assumptions
- Remote-only self-hosted runners are GitHub-visible but not locally controllable from this machine.
- They should appear as read-only inventory rows with GitHub status, labels, OS, version, and busy state.
- Existing local runners keep current lifecycle, cleanup, remove, and logs behavior.

## Milestone Order
| ID | Title | Depends on | Status |
| --- | --- | --- | --- |
| M1 | Remote-only row integration | - | [x] |
| M2 | Read-only safety and validation | M1 | [x] |
| M3 | Docs, research, and handoff | M2 | [~] |

## M1. Remote-only row integration `[x]`
### Goal
- Add GitHub API runners that do not match a local runner as read-only rows.

### Tasks
- [x] Detect unmatched entries returned by `repos/{owner}/{repo}/actions/runners`.
- [x] Convert unmatched entries into `Runner` rows with `Host=github`, `Local=remote`, `Transport=github-remote`, and non-local path text.
- [x] Preserve existing local matching and status enrichment behavior.

### Definition of Done
- Refresh output includes GitHub-visible self-hosted runners even without a local runner folder.

### Validation
```sh
go test ./internal/app
```

### Known Risks
- Remote-only rows must not accidentally be treated as local runners by command handlers.

### Stop-and-Fix Rule
- If a remote-only row can trigger local filesystem/service actions, add guards before continuing.

## M2. Read-only safety and validation `[x]`
### Goal
- Make remote-only rows safe and covered by regression tests.

### Tasks
- [x] Add read-only checks for lifecycle, cleanup, logs, removal, and reprovisioning paths where needed.
- [x] Add unit tests for unmatched GitHub runners and matched local runners.
- [x] Run full Go test suite.

### Definition of Done
- Remote-only rows are visible but cannot be controlled locally.

### Validation
```sh
go test ./...
```

### Known Risks
- Existing commands mostly guard GitHub-hosted rows only, not all non-local rows.

### Stop-and-Fix Rule
- Any failing test or unsafe command path blocks docs/handoff.

## M3. Docs, research, and handoff `[~]`
### Goal
- Record the behavior and publish session handoff.

### Tasks
- [x] Update user-facing docs or changelog if behavior changes are visible.
- [x] Update technology insights in repo and IdeaBox vault if a meaningful engineering finding emerges.
- [ ] Publish issue handoff comment.

### Definition of Done
- User can tell that partner runners are visible but read-only.

### Validation
```sh
git diff --check
```

### Known Risks
- Documentation can overpromise control of partner machines.

### Stop-and-Fix Rule
- If docs imply remote runners can be controlled, correct wording before handoff.
