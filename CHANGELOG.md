# Changelog

## [0.5.0] - 2026-06-07

### Added

- Self-hosted runners registered in monitored GitHub repositories now appear as
  read-only `github`/`remote` rows when no matching local runner folder is
  discovered.

### Changed

- Mutating commands now explicitly skip GitHub-hosted and GitHub-remote rows
  because they cannot be started, stopped, cleaned, removed, or reprovisioned
  by RunnerMonitor.

## [0.4.0] - 2026-06-07

### Added

- GitHub-hosted Actions jobs are now monitored as read-only rows beside
  self-hosted runners, using queued and in-progress workflow job data from
  `gh api`.
- New `githubHostedRepos` config field lets hosted-only repositories be
  monitored even when there is no local runner folder for that repository.

### Changed

- Hosted rows use `github` as the host, `hosted` as local state, and link to
  the GitHub Actions job URL in the path column.
- Mutating commands now explicitly skip GitHub-hosted rows because GitHub-hosted
  runners cannot be started, stopped, cleaned, removed, or reprovisioned by
  RunnerMonitor.

## [0.3.4] - 2026-06-05

### Fixed

- TUI command input now accepts `d` as the first typed character, so
  `delete N confirm` can be entered from an empty command line.
- TUI lifecycle commands now trigger an immediate fresh status refresh after
  `start`, `stop`, `restart`, and force variants, so local/GitHub state catches
  up without waiting for the next auto-refresh tick.
- TUI lifecycle command messages are preserved after the post-command refresh,
  so errors and elevated PowerShell handoff messages remain visible.
- Service-managed runner `stop` now waits for local service state to stop and
  for GitHub to report the runner offline before reporting success.
- Non-elevated Windows TUI lifecycle commands for Windows service runners now
  launch an elevated PowerShell helper instead of silently failing behind the
  refresh.

### Added

- CLI commands `--start-runner`, `--stop-runner`, and `--restart-runner` allow
  targeting one named runner, optionally disambiguated with `--repo`.

## [0.3.3] - 2026-06-05

### Fixed

- TUI update notices now render release URLs as OSC-8 terminal hyperlinks while
  keeping the visible URL readable in terminals that do not support clickable
  links.
- Service-managed runner `start` now waits for real readiness: systemd-backed
  runners are enabled before start, local service state must become active, and
  GitHub must report the runner as `online` before success is reported.

## [0.3.2] - 2026-06-04

### Added

- GitHub Actions CI now runs Go tests, module tidy verification, builds, and
  PowerShell syntax checks on push and pull requests.

### Fixed

- Runner status stays visible in the TUI table for both local and GitHub
  runner state.

## [0.3.1] - 2026-06-04

### Fixed

- Windows runner discovery PowerShell calls now have a timeout, preventing a
  stuck CIM query from keeping the TUI launch window alive after exit.

## [0.3.0] - 2026-06-04

### Added

- Automatic TUI refresh briefly reuses recent GitHub status to avoid repeated
  `gh` process fan-out while preserving fresh manual refresh behavior.

### Changed

- Runner registration and removal token handling now uses environment-based
  references where the upstream runner config scripts allow it.
- Remote TUI SSH commands quote configured paths for Windows and Linux hosts.
- Runner sorting precomputes comparison keys during refresh.

### Fixed

- Runner folder deletion now rejects configured root folders themselves and
  normalized Windows/WSL path traversal outside configured safe roots.
- Runner config command errors redact registration/remove tokens before
  surfacing failures.

## [0.2.1] - 2026-06-04

### Added

- TUI checks once on startup whether a newer GitHub release is available and
  shows a concise update notice when one exists.
- TUI has a concise in-app help panel opened with `h`, `?`, or `help`.

### Changed

- `Busy=true` is highlighted in the TUI table while `false` remains plain text.
- README and README_RU document in-app help and startup update notices.

## [0.2.0] - 2026-06-04

### Added

- TUI auto-refreshes local and GitHub runner state every 5 seconds.
- TUI auto-refresh interval is configurable through
  `tuiRefreshIntervalSeconds` in `runner-monitor.json`.
- Manual refresh keeps the existing table visible while new data is loading.

## [0.1.0] - 2026-06-04

### Added

- First public GitHub release of RunnerMonitor.
- Resize-aware Bubble Tea TUI for runner monitoring and commands.
- GitHub Actions runner discovery across Windows runner roots, WSL runner
  roots, and Linux runner roots.
- GitHub status, busy state, queue count, stale queue count, labels, project,
  local state, and runner path display.
- Project-scoped lifecycle commands: start, stop, restart, and current-project
  detection from Git remotes.
- Safe cleanup commands for idle runner work folders.
- Guarded runner removal and reprovisioning commands with dry-run defaults.
- App-local `runner-monitor.json` settings file beside the executable.
- SSH remote host profile workflow for future dedicated runner hosts.
- English and Russian documentation.
- MIT license and GitHub community files.

### Safety

- Busy runners are protected by default.
- Destructive operations require explicit confirmation.
- Folder deletion is limited to configured runner roots.
- `--show-config` masks `wslSudoPassword`.
- Release ZIP includes only a sanitized default config with an empty
  `wslSudoPassword`.

[0.5.0]: https://github.com/SGribanov/RunnerMonitor/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/SGribanov/RunnerMonitor/compare/v0.3.4...v0.4.0
[0.3.4]: https://github.com/SGribanov/RunnerMonitor/compare/v0.3.3...v0.3.4
[0.3.3]: https://github.com/SGribanov/RunnerMonitor/compare/v0.3.2...v0.3.3
[0.3.2]: https://github.com/SGribanov/RunnerMonitor/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/SGribanov/RunnerMonitor/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/SGribanov/RunnerMonitor/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/SGribanov/RunnerMonitor/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/SGribanov/RunnerMonitor/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/SGribanov/RunnerMonitor/releases/tag/v0.1.0
