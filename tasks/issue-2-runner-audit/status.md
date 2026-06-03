<!-- Issue: SGribanov/RunnerMonitor#2 -->
# Runner Audit Status

## Done

- Created issue #2.
- Added read-only `--audit` command.
- Current audit reports candidate removals for local-only or GitHub-unknown manual/inactive runners.
- Added `reports/runner-audit-2026-06-03.md` with read-only cleanup proposal and exact commands.

## Current audit snapshot

- Candidate remove: `windows-local`, `newgenosengine-windows-local`, `mycloneosengine-linux`, `newgen-wsl-linux`.
- Investigate: `backtester-runner`, DeltaG runners with queued jobs, `mycloneosengine-windows-local`, `legion-ubuntu-wsl-x64`.
- Keep: `ideabox-runner`.

## Next

- Operator decides which `candidate-remove` runners are actually safe to remove.

## Blockers

- Destructive cleanup requires explicit approval.
