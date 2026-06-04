<!-- Issue: SGribanov/RunnerMonitor#21 -->
# Review Fixes Test Plan

- Unit tests:
  - WSL/Linux safe-root path traversal is rejected.
  - Config command construction does not expose tokens in parent process argv where avoidable.
  - Secret masking removes token text from command errors.
  - Remote TUI command quotes configured paths.
  - Auto-refresh uses cached GitHub status while manual refresh stays fresh.
- Static checks:
  - `go test ./...`
  - `go vet ./...`
  - `govulncheck ./...`
  - `git diff --check`
- Race:
  - Try `go test -race ./...` with local MinGW path; record blocker if toolchain still fails.
