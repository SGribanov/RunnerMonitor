# Changelog

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

[0.1.0]: https://github.com/SGribanov/RunnerMonitor/releases/tag/v0.1.0
