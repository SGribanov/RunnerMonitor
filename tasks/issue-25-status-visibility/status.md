<!-- Issue: SGribanov/RunnerMonitor#25 -->
# Status Visibility Status

## 2026-06-04

Current phase: ready to merge.

Findings:
- `ideabox-runner` backend data is present: `Local=running`, `GitHub=online`, `Busy=true`.
- `--audit` showed `Windows runner discovery timed out after 5s`, so the 5-second bound from #23 is too aggressive for this machine.
- ANSI styling inside the `Busy=true` table cell can interfere with table width accounting and selected-row visibility.

Changes:
- Raised Windows discovery timeout to 30 seconds.
- Changed `Busy` table cell values back to plain `true`/`false` text.
- Updated README and README_RU.

Validation:
- `go test ./...` passed.
- `go vet ./...` passed.
- `go test -race ./...` passed.
- `git diff --check` passed.
- `runner-monitor.ps1 --once` shows `IdeaBox ideabox-runner running online true` without a Windows discovery timeout warning.
- Rebuilt local `bin\runner-monitor.exe`.
