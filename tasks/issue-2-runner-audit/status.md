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

## Current audit snapshot

- Candidate remove: `windows-local`, `newgenosengine-windows-local`, `mycloneosengine-linux`, `newgen-wsl-linux`.
- Investigate: `backtester-runner`, DeltaG runners with queued jobs, `mycloneosengine-windows-local`.
- Keep: `ideabox-runner`.

## Next

- Operator decides which remaining `candidate-remove` runners are actually safe to remove.

## Blockers

- Destructive cleanup requires explicit approval.
