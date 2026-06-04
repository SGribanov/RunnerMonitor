<!-- Issue: SGribanov/RunnerMonitor#5 -->
# Runner Directory Migration Status

## Current state

- Issue #5 created.
- Fresh audit shows all current runners are `keep`.
- DeltaG Windows and WSL runners are online and `busy=false`; both still show
  stale queue evidence from GitHub, but local capacity is available.
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
- Migrated `SGribanov/DeltaG deltag-linux-wsl` from
  `/home/gsv777/actions-runner-deltag` to
  `/home/gsv777/Runners/SGribanov-DeltaG/deltag-linux-wsl`.
- Created backup:
  `/home/gsv777/runner-backups/actions-runner-deltag-deltag-linux-wsl-move-2026-06-03.tar.gz`.
- Reinstalled WSL systemd service from the new path:
  `actions.runner.SGribanov-DeltaG.deltag-linux-wsl.service`.
- Fixed post-move runner symlinks:
  - `bin` -> `/home/gsv777/Runners/SGribanov-DeltaG/deltag-linux-wsl/bin.2.334.0`
  - `externals` -> `/home/gsv777/Runners/SGribanov-DeltaG/deltag-linux-wsl/externals.2.334.0`
- Validation passed:
  - service is `active`;
  - process path is under `/home/gsv777/Runners/SGribanov-DeltaG/deltag-linux-wsl`;
  - GitHub reports `deltag-linux-wsl` online and `busy=false`;
  - `runner-monitor --audit` still shows DeltaG stale queue investigation;
  - from `C:\Repos\DeltaG`, `--start-current` returns both DeltaG runners
    already running.
- Migrated `SGribanov/AU windows-local` from `C:\actions-runner` to
  `C:\Runners\SGribanov-AU\windows-local`.
- Created backup:
  `C:\Runners-backup\actions-runner-windows-local-move-2026-06-03.zip`.
- Fixed post-move runner junctions:
  - `bin` -> `C:\Runners\SGribanov-AU\windows-local\bin.2.333.0`
  - `externals` -> `C:\Runners\SGribanov-AU\windows-local\externals.2.333.0`
- Reattached GitHub runner binding because GitHub had zero registered AU
  runners while the local folder still had stale registration files.
- Fixed RunnerMonitor manual Windows start path passing so `--start-current`
  can launch paths under `C:\Runners`.
- Validation passed:
  - one `Runner.Listener.exe` process is running from
    `C:\Runners\SGribanov-AU\windows-local`;
  - GitHub reports `windows-local` online and `busy=false`;
  - `runner-monitor --audit` shows AU as `keep`, `running`, and `online`;
  - from `C:\Repos\AU`, `--start-current` returns `start windows-local requested`.
- Migrated `SGribanov/IdeaBox ideabox-runner` from
  `C:\actions-runner-ideabox` to
  `C:\Runners\SGribanov-IdeaBox\ideabox-runner`.
- Migrated `SGribanov/DeltaG deltag-win` from `C:\github-runners\deltag` to
  `C:\Runners\SGribanov-DeltaG\deltag-win`.
- Reconfigured both Windows services to the new `RunnerService.exe` paths and
  switched service startup to `Manual`.
- Fixed post-move Windows service runner junctions:
  - `IdeaBox bin` -> `C:\Runners\SGribanov-IdeaBox\ideabox-runner\bin.2.334.0`
  - `IdeaBox externals` -> `C:\Runners\SGribanov-IdeaBox\ideabox-runner\externals.2.334.0`
  - `DeltaG bin` -> `C:\Runners\SGribanov-DeltaG\deltag-win\bin.2.334.0`
  - `DeltaG externals` -> `C:\Runners\SGribanov-DeltaG\deltag-win\externals.2.334.0`
- Validation passed:
  - both Windows service processes run from `C:\Runners`;
  - GitHub reports `ideabox-runner` and `deltag-win` online and `busy=false`;
  - `runner-monitor --audit` no longer shows duplicate old DeltaG runner
    records.
- Cleanup completed after validation:
  - removed `C:\github-runners\deltag`;
  - removed Windows backup archives under `C:\Runners-backup`;
  - removed WSL backup archives under `/home/gsv777/runner-backups`;
  - cleared Windows and WSL runner `_work` directories;
  - removed runner installer zip/tar artifacts from runner folders;
  - Windows `C:\Runners` size dropped from about 16.5 GB to about 2.9 GB;
  - WSL runner `_work` directories and backup directory are down to about 4 KB.
- Autostart state:
  - Windows service runners are `StartMode=Manual`;
  - WSL systemd runner units are `disabled` but currently `active`;
  - manual Windows runners are started from project commands and remain
    controllable by RunnerMonitor.

## Next

- No runner folder migration work remains.
- DeltaG still has a separately documented stale GitHub queue item; local
  runner capacity is online and not busy.

## Blockers

- None for runner folder migration.
