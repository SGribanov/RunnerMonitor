<!-- Issue: SGribanov/RunnerMonitor#2 -->
# Runner Audit Status

## Done

- Created issue #2.
- Added read-only `--audit` command.
- Current audit reports candidate removals for local-only or GitHub-unknown manual/inactive runners.
- Added `reports/runner-audit-2026-06-03.md` with read-only cleanup proposal and exact commands.
- Removed approved runner `legion-ubuntu-wsl-x64` at `/home/gsv777/actions-runner-linux-x64`.
- Created backup `/home/gsv777/runner-backups/actions-runner-linux-x64-legion-ubuntu-wsl-x64-2026-06-03.tar.gz`.
- Removed approved runner `legion-windows-x64` at `C:\actions-runner-win-x64`.
- Deleted GitHub runner registration id `21` from `SGribanov/DeltaG`.
- Created backup `C:\Runners-backup\actions-runner-win-x64-legion-windows-x64-2026-06-03.zip`.
- Removed approved runner `newgenosengine-windows-local` at `C:\actions-runner-newgenosengine`.
- Created backup `C:\Runners-backup\actions-runner-newgenosengine-windows-local-2026-06-03.zip`.
- Removed approved runner `newgen-wsl-linux` at `/home/gsv777/newgen-runner`.
- Created backup `/home/gsv777/runner-backups/newgen-runner-newgen-wsl-linux-2026-06-03.tar.gz`.
- Could not remove its WSL systemd unit without sudo password; unit cleanup remains elevated.

## Current audit snapshot

- Candidate remove: `windows-local`, `mycloneosengine-linux`.
- Investigate: `backtester-runner`, DeltaG runners with queued jobs, `mycloneosengine-windows-local`.
- Keep: `ideabox-runner`.

## Next

- Operator decides which remaining `candidate-remove` runners are actually safe to remove.
- Elevated cleanup remains for stale WSL unit `actions.runner.SGribanov-NewGenOsEngine.newgen-wsl-linux.service`.

## Blockers

- Destructive cleanup requires explicit approval.
