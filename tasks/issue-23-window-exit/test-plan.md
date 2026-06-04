<!-- Issue: SGribanov/RunnerMonitor#23 -->
# Window Exit Test Plan

- Unit:
  - Force Windows discovery PowerShell timeout and verify a timeout warning.
- Static:
  - `go test ./...`
  - `go vet ./...`
  - `govulncheck ./...`
  - `git diff --check`
- Smoke:
  - `scripts/build.ps1`
  - `runner-monitor.ps1 --once`
- Known blocker:
  - `go test -race ./...` requires a working Windows cgo compiler toolchain.
