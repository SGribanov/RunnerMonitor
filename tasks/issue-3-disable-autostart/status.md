<!-- Issue: SGribanov/RunnerMonitor#3 -->
# Disable Runner Autostart Status

## Done

- Created issue #3.
- Added project-scoped lifecycle commands.
- Added `--disable-autostart`.
- Added `--start-current` so Codex can start runners from any project root.
- Added `runner-monitor.ps1` wrapper and `scripts/build.ps1`.
- Updated parent `C:\Repos\AGENTS.md` with the Runner Startup Policy.
- Confirmed current non-elevated session cannot change startup policy.

## Current blocker

- Windows service startup changes require elevated rights: `Access is denied`.
- WSL systemd disable requires authentication: `Interactive authentication required`.

## Next

- Run the required elevated commands from the plan.
- Use this command from a project root before push/CI work:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --start-current
```

- Recheck with:

```powershell
Get-CimInstance Win32_Service | Where-Object { $_.Name -like 'actions.runner.*' } | Select-Object Name,State,StartMode
wsl.exe systemctl list-unit-files 'actions.runner.*.service' --no-pager
```
