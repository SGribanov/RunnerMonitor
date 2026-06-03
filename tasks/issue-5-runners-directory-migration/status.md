<!-- Issue: SGribanov/RunnerMonitor#5 -->
# Runner Directory Migration Status

## Current state

- Issue #5 created.
- Fresh audit shows all current runners are `keep`.
- DeltaG Windows and WSL runners are currently busy, so they are not migration
  candidates until idle.
- Migrated `SGribanov/BackTester backtester-runner` from
  `C:\actions-runner-backtester` to
  `C:\Runners\SGribanov-BackTester\backtester-runner`.
- Created backup:
  `C:\Runners-backup\actions-runner-backtester-backtester-runner-2026-06-03.zip`.
- Fixed post-move runner junctions:
  - `bin` -> `C:\Runners\SGribanov-BackTester\backtester-runner\bin.2.334.0`
  - `externals` -> `C:\Runners\SGribanov-BackTester\backtester-runner\externals.2.334.0`
- Validation passed:
  - one `Runner.Listener.exe` process remains;
  - process path is under `C:\Runners\SGribanov-BackTester\backtester-runner`;
  - GitHub reports `backtester-runner` online and `busy=false`;
  - `runner-monitor --audit` shows BackTester as `keep`;
  - from `C:\Repos\BackTester`, `--start-current` returns `backtester-runner already running`.
- Migrated `SGribanov/MyCloneOsEngine mycloneosengine-windows-local` from
  `C:\actions-runner-mycloneosengine` to
  `C:\Runners\SGribanov-MyCloneOsEngine\mycloneosengine-windows-local`.
- Created backup:
  `C:\Runners-backup\actions-runner-mycloneosengine-mycloneosengine-windows-local-2026-06-03.zip`.
- Fixed post-move runner junctions:
  - `bin` -> `C:\Runners\SGribanov-MyCloneOsEngine\mycloneosengine-windows-local\bin.2.334.0`
  - `externals` -> `C:\Runners\SGribanov-MyCloneOsEngine\mycloneosengine-windows-local\externals.2.334.0`
- Validation passed:
  - one `Runner.Listener.exe` process remains;
  - process path is under `C:\Runners\SGribanov-MyCloneOsEngine\mycloneosengine-windows-local`;
  - GitHub reports `mycloneosengine-windows-local` online and `busy=false`;
  - `runner-monitor --audit` shows MyClone Windows as `keep`;
  - from `C:\Repos\MyCloneOsEngine`, `--start-current` returns both MyClone
    runners already running.
- Candidate first moves, after explicit approval:
  - `SGribanov/MyCloneOsEngine mycloneosengine-linux`
  - `SGribanov/IdeaBox ideabox-runner` if admin rights are available

## Next

- Continue with the next non-busy runner, preferably
  `SGribanov/MyCloneOsEngine mycloneosengine-linux`.
- Include junction retargeting after every Windows runner folder move.

## Blockers

- Runner moves are intentionally blocked until explicit runner-by-runner
  approval.
- Windows service moves may require elevated PowerShell.
- WSL systemd moves require sudo.
