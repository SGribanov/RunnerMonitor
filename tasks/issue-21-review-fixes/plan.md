<!-- Issue: SGribanov/RunnerMonitor#21 -->
# Review Fixes Plan

Goal: fix all security, leak, allocation, and performance findings from the full-project review.

- [x] Harden safe runner folder checks for WSL/Linux and Windows deletion.
- [x] Reduce runner token exposure and mask token values in process error output.
- [x] Throttle GitHub polling during TUI auto-refresh.
- [x] Quote remote TUI command paths.
- [x] Remove avoidable sort comparator allocations.
- [x] Update docs/research notes.
- [x] Run validation: `go test ./...`, `go vet ./...`, `govulncheck ./...`, race test if local toolchain supports it.
- [ ] Publish GitHub handoff.
