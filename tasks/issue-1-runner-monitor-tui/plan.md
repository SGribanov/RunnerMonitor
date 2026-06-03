<!-- Issue: SGribanov/RunnerMonitor#1 -->
# RunnerMonitor TUI Plan

## Milestone 1: Bootstrap and inventory

Status: [~]

Goal: create the repo skeleton and discover current GitHub Actions runners from Windows, WSL, and GitHub API.

Tasks:
- [x] Create GitHub repository, project, issue, and working branch.
- [x] Initialize Go module.
- [x] Discover local `.runner` files without reading credentials.
- [x] Discover Windows service state for `actions.runner.*`.
- [x] Discover WSL runner folders and systemd service hints.
- [x] Merge local records with GitHub runner status and queued workflow counts.

Definition of done:
- `go test ./...` passes.
- `go run ./cmd/runner-monitor` lists current local runners and GitHub status.

Validation:
- `go test ./...`
- `go run ./cmd/runner-monitor`

Known risks:
- Some configured runner folders are not service-managed and must be marked as manual.
- WSL/systemd service control may require privileges.

Stop-and-fix rule:
- If discovery reads credentials or fails on missing GitHub access, stop and narrow the provider behavior before adding lifecycle commands.

## Milestone 2: TUI commands

Status: [~]

Goal: provide a keyboard-first interface with numeric runner commands.

Tasks:
- [x] Render a compact table with host, repo, runner, local state, GitHub state, busy, queue, labels, and path.
- [x] Parse `start N`, `stop N`, `restart N`, `logs N`, `refresh`, and `q`.
- [x] Execute lifecycle commands only for service-managed runners.
- [x] Show clear command output for unsupported manual runners.

Definition of done:
- Commands work against service-managed Windows runners.
- Unsupported runners are safe and explicit.

Validation:
- `go test ./...`
- Manual dry run in the TUI.

## Milestone 3: Remote-ready design

Status: [ ]

Goal: prepare the provider boundary for a future dedicated runner host over SSH.

Tasks:
- [ ] Keep host and transport in every runner record.
- [ ] Isolate local Windows, WSL, GitHub, and lifecycle behavior behind provider/controller functions.
- [ ] Document SSH host configuration shape without requiring a daemon.

Definition of done:
- Future SSH provider can be added without changing the TUI command model.

Validation:
- `go test ./...`

## Assumptions

- The TUI runs on the operator machine.
- Future runner machine control uses SSH rather than a daemon.
- Runner folder migration into `Runners` is a separate phase.
- OneDev support is a future provider, not part of v1.
