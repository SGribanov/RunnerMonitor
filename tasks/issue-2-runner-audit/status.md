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
- WSL systemd unit cleanup was completed manually; `actions.runner.SGribanov-NewGenOsEngine.newgen-wsl-linux.service` is no longer present.
- Marked `AU/windows-local` as keep via `runner-policy.json`.
- Reattached `mycloneosengine-linux`; it is now service-managed, active, and GitHub online.

## Current audit snapshot

- Candidate remove: none.
- Investigate: `backtester-runner`, DeltaG runners with queued jobs, `mycloneosengine-windows-local`.
- Keep: `windows-local`, `ideabox-runner`, `mycloneosengine-linux`.

## Next

- Investigate remaining manual/queued runners.

## Blockers

- Destructive cleanup requires explicit approval.
