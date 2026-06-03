<!-- Issue: SGribanov/RunnerMonitor#2 -->
# Runner Audit Test Plan

## Automated

- Audit keeps busy runners.
- Audit investigates queued repositories.
- Audit marks local-only manual runners as candidate removal.

## Manual

- Run `go run ./cmd/runner-monitor --audit`.
- Confirm no delete/stop commands are executed.
- Review every candidate before removal.

