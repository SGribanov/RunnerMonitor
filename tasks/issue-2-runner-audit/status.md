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
- Investigate: DeltaG still has queued jobs.
- Keep: `windows-local`, `backtester-runner`, `ideabox-runner`, `mycloneosengine-windows-local`, `mycloneosengine-linux`.

## DeltaG queue note

- Remaining queued run: `26447257991`.
- Workflow/title: `ci` / `research(margin): diagnose vertical freshness drift`.
- Branch/event: `codex/604-vertical-freshness-diagnostic` / `pull_request`.
- Created/updated: `2026-05-26T10:36:43Z` / `2026-05-26T10:52:39Z`.
- Jobs: `0`.
- Related PRs `#709` and `#710` are closed.
- At diagnosis time, DeltaG runners `deltag-win` and `deltag-linux-wsl` were online and not busy.
- Attempted GitHub cancel API: `POST /repos/SGribanov/DeltaG/actions/runs/26447257991/cancel` returned HTTP `409` with `Cannot cancel a workflow re-run that has not yet queued`; run still reports `queued`.

## Next

- Teach audit to ignore or separately flag closed-PR/no-job stale runs, because the GitHub cancel API cannot cancel this one.

## Blockers

- Destructive cleanup requires explicit approval.
