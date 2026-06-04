<!-- Issue: SGribanov/RunnerMonitor#6 -->
# Safe Runner Cleanup Plan

## Goal

Add a safe cleanup command for runner bloat without deleting registration,
credentials, versioned runner binaries, or service configuration.

## Scope

- TUI command `clear N` for one runner.
- TUI command `clear idle` for all idle runners.
- TUI toggle `auto-clear on/off` that clears idle runners after refresh.
- CLI commands:
  - `--clear-current`
  - `--clear-repo owner/repo`
  - `--clear-idle`

## Safety Rules

- Refuse cleanup when GitHub reports `busy=true`.
- If the runner is running/active, stop only that runner before cleanup.
- Restart the runner after cleanup if it was running/active before cleanup.
- Clear only `_work` contents and runner installer archives.
- Preserve `.runner`, `.credentials`, `.credentials_rsaparams`, `bin.*`,
  `externals.*`, service files, and runner registration.

## Automatic Cleanup

`auto-clear on` is intentionally opt-in. It clears idle runners after a TUI
refresh instead of running as an always-on background delete loop, reducing the
risk of racing a newly assigned job.
