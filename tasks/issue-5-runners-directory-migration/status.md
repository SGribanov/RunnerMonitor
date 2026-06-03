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

## Next

- Continue with the next non-busy runner.
- Windows service moves (`IdeaBox`, `DeltaG Windows`) require elevated
  PowerShell; current shell is not elevated.
- `SGribanov/AU windows-local` has no GitHub runner binding at the moment, so
  it needs re-registration or a fresh GitHub runner attach before normal move
  validation can pass.
- Include junction/symlink retargeting after every moved runner folder.

## Blockers

- Windows service moves may require elevated PowerShell.
- WSL systemd moves require sudo.
