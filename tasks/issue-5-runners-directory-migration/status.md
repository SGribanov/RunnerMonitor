<!-- Issue: SGribanov/RunnerMonitor#5 -->
# Runner Directory Migration Status

## Current state

- Issue #5 created.
- Fresh audit shows all current runners are `keep`.
- DeltaG Windows and WSL runners are currently busy, so they are not migration
  candidates until idle.
- Candidate first moves, after explicit approval:
  - `SGribanov/BackTester backtester-runner`
  - `SGribanov/MyCloneOsEngine mycloneosengine-windows-local`
  - `SGribanov/MyCloneOsEngine mycloneosengine-linux`
  - `SGribanov/IdeaBox ideabox-runner` if admin rights are available

## Next

- Choose the first runner to migrate.
- Produce runner-specific backup, stop, move, reconfigure, start, validate, and
  rollback commands.

## Blockers

- Runner moves are intentionally blocked until explicit runner-by-runner
  approval.
- Windows service moves may require elevated PowerShell.
- WSL systemd moves require sudo.
