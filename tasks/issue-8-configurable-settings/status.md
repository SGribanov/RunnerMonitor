<!-- Issue: SGribanov/RunnerMonitor#8 -->
# Configurable Settings Status

## Current Phase

Milestones 1-3 implemented; validation passed and handoff is ready.

## Done

- Created GitHub issue #8.
- Created branch `codex/8-configurable-roots`.
- Added issue #8 to RunnerMonitor Project board and set status to
  `In Progress`.
- Updated issue scope after user clarification: config stores the WSL sudo
  password value, not a path to a password file.
- Updated config location after user clarification: `runner-monitor.json` lives
  next to the compiled executable.
- Added settings model, defaults, `RUNNER_MONITOR_CONFIG` override,
  `--init-config`, and `--show-config`.
- Wired configured roots into Windows/WSL/Linux discovery, project folder
  resolution, safe runner-folder deletion, and WSL sudo fallback.
- Updated build script to create `runner-monitor.json` next to
  `runner-monitor.exe` without overwriting an existing file.
- Updated README and research insights in both the repo and IdeaBox vault.
- Confirmed the app-local config masks the real WSL sudo password in
  `--show-config`.
- Confirmed no real WSL sudo password appears in repository files outside
  ignored build output.

## In Progress

- None.

## Next

- Commit and push branch `codex/8-configurable-roots`.
- Close issue #8 after publishing the final issue handoff.

## Decisions

- Real sudo password belongs only in the app-local config file, never in repo
  files.
- `--show-config` masks the password.
- Missing config should behave like the current workstation defaults.
- `RUNNER_MONITOR_CONFIG` remains available as an override for tests and special
  launches.

## Validation Log

- Passed: `go test ./...`
- Passed: `scripts\build.ps1`
- Passed: `runner-monitor.ps1 --show-config` reads
  `C:\Repos\RunnerMonitor\bin\runner-monitor.json` and masks
  `wslSudoPassword` as `<set>`.
- Passed: `runner-monitor.ps1 --audit` discovers the current Windows and WSL
  runners from configured roots.
- Passed: no-real-password repository scan excluding `bin`.
