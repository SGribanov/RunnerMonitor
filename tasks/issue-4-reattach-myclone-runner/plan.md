<!-- Issue: SGribanov/RunnerMonitor#4 -->
# Reattach MyClone WSL Runner Plan

## Milestone 1: Preserve and reattach

Status: [x]

Goal: keep AU runner and reattach MyClone Linux runner.

Tasks:
- [x] Mark `SGribanov/AU windows-local` as keep.
- [x] Back up existing MyClone WSL runner config.
- [x] Reconfigure `/home/gsv777/myclone-runner-linux` with GitHub registration token.
- [x] Install and start WSL systemd service.
- [x] Verify GitHub API reports `mycloneosengine-linux` online.

Validation:
- `go test ./...`
- `runner-monitor --audit`
- `gh api repos/SGribanov/MyCloneOsEngine/actions/runners`

