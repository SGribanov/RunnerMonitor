# Changelog

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

[0.3.0]: https://github.com/SGribanov/RunnerMonitor/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/SGribanov/RunnerMonitor/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/SGribanov/RunnerMonitor/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/SGribanov/RunnerMonitor/releases/tag/v0.1.0
