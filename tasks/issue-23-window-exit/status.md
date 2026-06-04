<!-- Issue: SGribanov/RunnerMonitor#23 -->
# Window Exit Status

## 2026-06-04

Current phase: validation complete.

Done:
- Found a lingering PowerShell subprocess with the inherited `runner-monitor.exe` window title.
- Identified the command as Windows service discovery via `Get-CimInstance Win32_Service`.
- Added a 5-second timeout wrapper for Windows discovery PowerShell calls.
- Added a test that forces the timeout path and verifies a warning is returned.

Validation:
- `go test ./...` passed.
- `go vet ./...` passed.
- `govulncheck ./...` passed.
- `git diff --check` passed.
- `scripts/build.ps1` passed.
- `runner-monitor.ps1 --once` passed.
- `go test -race ./...` remains blocked by local cgo/C compiler toolchain; with `CGO_ENABLED=1`, this environment reports `gcc not found in %PATH%`.
