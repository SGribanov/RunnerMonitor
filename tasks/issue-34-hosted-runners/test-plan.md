<!-- Issue: SGribanov/RunnerMonitor#34 -->
# Test Plan

## Source
- Task: Monitor GitHub-hosted Actions runners and publish a release.
- Plan file: `tasks/issue-34-hosted-runners/plan.md`
- Status file: `tasks/issue-34-hosted-runners/status.md`
- Last updated: 2026-06-07

## Validation Scope
- In scope: hosted workflow job API parsing, hosted/self-hosted classification, read-only hosted command handling, docs/build/release checks.
- Out of scope: starting/stopping GitHub-hosted runners, enterprise larger-runner inventory management, billing/cost reporting.

## Environment / Fixtures
- Unit fixtures: JSON responses from GitHub workflow runs/jobs APIs.
- External dependencies: `gh` CLI for live refresh and release publishing.
- Build: repository `scripts/build.ps1`.

## Test Levels

### Unit
- Parse workflow job responses with queued and in-progress jobs.
- Ignore jobs with `self-hosted` labels for hosted-job rows.
- Convert GitHub-hosted jobs into read-only rows.
- Reject lifecycle/cleanup/remove/delete/log commands for hosted rows.

### Integration
- `go test ./...` verifies all packages and current command behavior.

### Release Smoke
- `.\scripts\build.ps1` produces the Windows binary.
- `gh release view <tag>` confirms release metadata and uploaded artifact after publish.

## Negative / Edge Cases
- GitHub jobs endpoint fails for one run while other repo data is still usable.
- Workflow runs have no jobs.
- Hosted row is selected and a mutating command is typed.
- Job labels are empty or mixed case.

## Acceptance Gates
- [x] `go test ./internal/app`
- [x] `go test ./...`
- [x] `.\scripts\build.ps1`
- [ ] Release artifact uploaded.

## Release / Demo Readiness
- [x] README and README_RU document hosted monitoring and read-only limits.
- [x] CHANGELOG includes the hosted runner feature.
- [x] Technology insights updated in repo and IdeaBox vault.
- [ ] Issue #34 handoff comment published.

## Command Matrix
```sh
go test ./internal/app
go test ./...
.\scripts\build.ps1
gh release view <tag>
```
