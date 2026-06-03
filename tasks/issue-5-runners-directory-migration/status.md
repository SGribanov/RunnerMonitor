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
- Migrated `SGribanov/MyCloneOsEngine mycloneosengine-linux` from
  `/home/gsv777/myclone-runner-linux` to
  `/home/gsv777/Runners/SGribanov-MyCloneOsEngine/mycloneosengine-linux`.
- Created backup:
  `/home/gsv777/runner-backups/myclone-runner-linux-mycloneosengine-linux-move-2026-06-03.tar.gz`.
- Reinstalled WSL systemd service from the new path:
  `actions.runner.SGribanov-MyCloneOsEngine.mycloneosengine-linux.service`.
- Fixed post-move runner symlinks:
  - `bin` -> `/home/gsv777/Runners/SGribanov-MyCloneOsEngine/mycloneosengine-linux/bin.2.334.0`
  - `externals` -> `/home/gsv777/Runners/SGribanov-MyCloneOsEngine/mycloneosengine-linux/externals.2.334.0`
- Validation passed:
  - service is `active`;
  - process path is under `/home/gsv777/Runners/SGribanov-MyCloneOsEngine/mycloneosengine-linux`;
  - GitHub reports `mycloneosengine-linux` online and `busy=false`;
  - `runner-monitor --audit` shows MyClone Linux as `keep`;
  - from `C:\Repos\MyCloneOsEngine`, `--start-current` returns both MyClone
    runners already running.
- Candidate first moves, after explicit approval:
  - `SGribanov/IdeaBox ideabox-runner` if admin rights are available

## Next

- Continue with the next non-busy runner, preferably
  `SGribanov/IdeaBox ideabox-runner` if admin rights are available.
- Include junction/symlink retargeting after every moved runner folder.

## Blockers

- Runner moves are intentionally blocked until explicit runner-by-runner
  approval.
- Windows service moves may require elevated PowerShell.
- WSL systemd moves require sudo.
