# Startup Update Check Status

## 2026-06-04

- Added a best-effort startup update check for the TUI.
- The check uses `gh release view --repo SGribanov/RunnerMonitor --json tagName,url` with a short timeout.
- Offline, unauthenticated, or failed checks are silent so runner monitoring still starts normally.
- A newer release is shown as a separate notice line above the command input.
- Scope update: the TUI Busy column now highlights only `true`; `false` remains plain text.
