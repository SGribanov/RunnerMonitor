<!-- Issue: SGribanov/RunnerMonitor#25 -->
# Status Visibility Test Plan

- `go test ./...`
- `go vet ./...`
- `go test -race ./...`
- `git diff --check`
- `runner-monitor.ps1 --once`
- Rebuild local `bin\runner-monitor.exe`.
