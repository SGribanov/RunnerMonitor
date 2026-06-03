# RunnerMonitor -- Strategy Insights

| Field | Value |
|---|---|
| Project | RunnerMonitor |
| Type | strategy-research |
| Last updated | 2026-06-03 |
| Status | active |
| Tags | ci-runners, github-actions, tui, operations |

## 2026-06-03 -- Monitoring goal

RunnerMonitor should not duplicate GitHub's UI. Its useful product boundary is
combining local control state, GitHub runner status, and queue health in one
keyboard-first view. The key operator workflow is short numbered commands such
as `start 1` and `stop 1`, backed by automatic discovery so new runners appear
without per-runner configuration.

## 2026-06-03 -- Remote topology

The accepted operating model is local TUI over SSH. The TUI runs on the
operator machine and controls the future dedicated runner host remotely. This
keeps v1 lightweight and avoids introducing a daemon before OneDev migration
requirements are known.

## 2026-06-03 -- Runner cleanup policy

Runner cleanup must be evidence-driven and explicit. The app should classify
runners as `keep`, `investigate`, or `candidate-remove`, but actual deletion of
services, GitHub registrations, or folders requires separate approval for each
runner.

## 2026-06-03 -- First cleanup candidates

The first audit identifies four candidate removals: AU `windows-local`, NewGen
Windows `newgenosengine-windows-local`, MyClone Linux `mycloneosengine-linux`,
and NewGen WSL `newgen-wsl-linux`. DeltaG runners stay under `investigate`
because the repository currently has a stale queued workflow, so label/routing
must be checked before removing anything.

## 2026-06-03 -- First approved removal

`legion-ubuntu-wsl-x64` was explicitly approved for deletion and removed from
`/home/gsv777/actions-runner-linux-x64` after creating backup
`/home/gsv777/runner-backups/actions-runner-linux-x64-legion-ubuntu-wsl-x64-2026-06-03.tar.gz`.
Post-removal audit no longer lists that runner.

## 2026-06-03 -- Second approved removal

`legion-windows-x64` was explicitly approved for deletion. The GitHub runner
registration id `21` was deleted from `SGribanov/DeltaG`; local folder
`C:\actions-runner-win-x64` was backed up to
`C:\Runners-backup\actions-runner-win-x64-legion-windows-x64-2026-06-03.zip`
and removed. Post-removal audit no longer lists that runner.

## 2026-06-03 -- Third approved removal

`newgenosengine-windows-local` was explicitly approved for deletion and removed
from `C:\actions-runner-newgenosengine` after backup
`C:\Runners-backup\actions-runner-newgenosengine-windows-local-2026-06-03.zip`.
GitHub API already showed zero active runners for `SGribanov/NewGenOsEngine`.
Post-removal audit no longer lists that runner.

## 2026-06-03 -- Fourth approved removal

`newgen-wsl-linux` was explicitly approved for deletion and removed from
`/home/gsv777/newgen-runner` after backup
`/home/gsv777/runner-backups/newgen-runner-newgen-wsl-linux-2026-06-03.tar.gz`.
GitHub API already showed zero active runners for `SGribanov/NewGenOsEngine`.
Post-removal audit no longer lists that runner. The systemd unit remains enabled
until an interactive sudo cleanup is performed.

## 2026-06-03 -- Manual startup policy

Runners should not auto-start at boot. The desired workflow is for Codex or the
operator to start only the required project runner with a short command such as
`runner-monitor --start-repo SGribanov/DeltaG`.

## 2026-06-03 -- Current-project startup command

The preferred Codex workflow is even simpler than passing an explicit repo:
from any project root, run `runner-monitor --start-current`. The command
derives `owner/repo` from `git remote get-url origin`, which avoids hard-coding
repo names in per-project instructions.
