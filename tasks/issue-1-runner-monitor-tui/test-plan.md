<!-- Issue: SGribanov/RunnerMonitor#1 -->
# RunnerMonitor Test Plan

## Automated tests

- Parse `.runner` JSON, including UTF-8 BOM.
- Parse GitHub repository URLs from `.runner` data.
- Merge local runner records with GitHub API runner status.
- Attach queued and stale queued workflow counts by repository.
- Parse command input for lifecycle commands and invalid runner ids.

## Manual smoke tests

- Run `go run ./cmd/runner-monitor`.
- Confirm current Windows and WSL runner folders appear.
- Confirm GitHub online/offline/busy status appears where API access exists.
- Confirm `SGribanov/DeltaG` shows queued jobs.
- Try `logs N` on a known runner.
- Use lifecycle commands only on a non-critical service-managed runner.

## Release gates

- `go test ./...` passes.
- No credential files such as `.credentials` are read or displayed.
- Unsupported manual runners are clearly marked and not controlled as services.

