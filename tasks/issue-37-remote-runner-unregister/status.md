<!-- Issue: SGribanov/RunnerMonitor#37 -->
# Status

## Snapshot
- Current phase: M3 handoff
- Plan file: `tasks/issue-37-remote-runner-unregister/plan.md`
- Status: validated
- Last updated: 2026-06-08

## Done
- Synchronized repository state with `git fetch --all --prune`, `git status --short --branch`, and `git pull --rebase`.
- Checked GitHub issue/project context.
- Created canonical GitHub issue #37 and added it to Project #25.
- Set GitHub Project Status to `In Progress`.
- Created branch `gsv777/issue-37-remote-runner-unregister`.
- Refreshed GitHub self-hosted runner deletion API via Exa against official GitHub docs.
- Implemented GitHub runner ID parsing and remote-only unregister flow.
- Added regression tests for dry-run, confirmed DELETE endpoint, missing ID, and ID propagation.
- Updated README, README_RU, CHANGELOG, task files, and technology insights.
- Synced technology insights to `D:\Repos\IdeaBox\vault\my-research\RunnerMonitor_technology_insights.md`.
- Verified IdeaBox watcher is running.
- Verified `go test ./internal/app`, `go test ./...`, `git diff --check`, `scripts/build.ps1`, `bin\runner-monitor.exe --once`, and remote-only dry-run for `deltag-win-2`.
- Published GitHub issue handoff comment: https://github.com/SGribanov/RunnerMonitor/issues/37#issuecomment-4645764111.
- Bumped release docs/version to `v0.5.1` and built release ZIP/checksum.
- Latest smoke no longer finds `deltag-win-2` in GitHub runner inventory.

## In Progress
- Awaiting push/merge/release.

## Next
- Push/merge/release `v0.5.1`.

## Decisions Made
- `github-hosted` rows remain read-only.
- `github-remote` rows can be unregistered only after explicit confirmation.
- `delete ... confirm` on a `github-remote` row unregisters GitHub only and does not attempt local folder deletion.

## Assumptions In Force
- `deltag-win-2` is a repository self-hosted runner registration visible through GitHub API.
- GitHub API deletion requires repository admin access through `gh`.

## Commands
```sh
git fetch --all --prune
git status --short --branch
git pull --rebase
gh issue create --title "Allow unregistering remote-only self-hosted runners" ...
gh project item-edit ...
git checkout -b gsv777/issue-37-remote-runner-unregister
go test ./internal/app
go test ./...
git diff --check
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
.\bin\runner-monitor.exe --once
.\bin\runner-monitor.exe --remove-runner deltag-win-2 --repo SGribanov/DeltaG
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1
```

## Current Blockers
- None

## Audit Log
| Date | Milestone | Files | Commands | Result | Next |
| --- | --- | --- | --- | --- | --- |
| 2026-06-08 | Setup | GitHub issue/project | `git fetch`; `git pull --rebase`; `gh issue create`; `gh project item-edit` | issue #37 and branch ready | implementation |
| 2026-06-08 | M1-M2 | `internal/app/*` | `gofmt`; `go test ./internal/app` | pass | docs and full validation |
| 2026-06-08 | M3 | docs, research, binary | `go test ./...`; `git diff --check`; `scripts/build.ps1`; `bin\runner-monitor.exe --once`; dry-run remove | pass; `deltag-win-2` found as GitHub id 26 | issue handoff |
| 2026-06-08 | Handoff | GitHub issue #37 | `gh issue comment 37` | published handoff comment | commit/push when ready |
| 2026-06-08 | Release prep | version/docs/assets | `scripts/build.ps1`; package script; dry-run remove | v0.5.1 ZIP/checksum built; `deltag-win-2` no longer in inventory | push/merge/release |

## Smoke / Demo Checklist
- [x] Remote-only dry-run explains GitHub API unregister.
- [x] Confirmed remote-only removal calls repository runner DELETE endpoint by ID.
- [x] GitHub-hosted rows remain read-only.
- [x] Full Go tests pass.
- [x] Local binary rebuilt.
- [x] Handoff comment published.
