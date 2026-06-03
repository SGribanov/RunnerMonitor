# Runner Audit -- 2026-06-03

| Field | Value |
|---|---|
| Project | RunnerMonitor |
| Source | `runner-monitor --audit` |
| Status | active |

## Summary

Approved removals and one reconfiguration have been completed. The latest audit
has no `candidate-remove` rows. Do not delete any additional runner without a
new explicit approval naming repo, runner, host, and path.

## Completed Removals

| Repo | Runner | Path | Backup |
|---|---|---|---|
| `SGribanov/DeltaG` | `legion-ubuntu-wsl-x64` | `/home/gsv777/actions-runner-linux-x64` | `/home/gsv777/runner-backups/actions-runner-linux-x64-legion-ubuntu-wsl-x64-2026-06-03.tar.gz` |
| `SGribanov/DeltaG` | `legion-windows-x64` | `C:\actions-runner-win-x64` | `C:\Runners-backup\actions-runner-win-x64-legion-windows-x64-2026-06-03.zip` |
| `SGribanov/NewGenOsEngine` | `newgenosengine-windows-local` | `C:\actions-runner-newgenosengine` | `C:\Runners-backup\actions-runner-newgenosengine-windows-local-2026-06-03.zip` |
| `SGribanov/NewGenOsEngine` | `newgen-wsl-linux` | `/home/gsv777/newgen-runner` | `/home/gsv777/runner-backups/newgen-runner-newgen-wsl-linux-2026-06-03.tar.gz` |

## Completed Reconfiguration

| Repo | Runner | Path | Result |
|---|---|---|---|
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-linux` | `/home/gsv777/myclone-runner-linux` | Reconfigured with GitHub runner id `24`; systemd unit `actions.runner.SGribanov-MyCloneOsEngine.mycloneosengine-linux.service` is active/enabled. |

## Resolved Elevated Cleanup

The `newgen-wsl-linux` runner folder was removed first. Its systemd unit was
then removed manually from WSL with sudo. Current WSL unit list no longer shows
`actions.runner.SGribanov-NewGenOsEngine.newgen-wsl-linux.service`.

## Keep

| Repo | Runner | Host | Evidence |
|---|---|---|---|
| `SGribanov/AU` | `windows-local` | local Windows | Policy keep: AU project will continue. |
| `SGribanov/IdeaBox` | `ideabox-runner` | local Windows | Service-managed and GitHub online. |
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-linux` | WSL Ubuntu | Service-managed, systemd active, GitHub online. |

## Investigate Before Any Cleanup

| Repo | Runner | Host | Evidence | Suggested action |
|---|---|---|---|---|
| `SGribanov/BackTester` | `backtester-runner` | local Windows | Process running and GitHub online, but not service-managed locally. | Decide whether to install as service or keep manual/background. |
| `SGribanov/DeltaG` | `deltag-win` | local Windows | Running, GitHub online, repo has `1/1 stale` queued jobs. | Keep for now; investigate queued job label/routing first. |
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-windows-local` | local Windows | Process running and GitHub online, but not service-managed locally. | Decide whether it should become service-managed. |
| `SGribanov/DeltaG` | `deltag-linux-wsl` | WSL Ubuntu | Active, GitHub online, repo has `1/1 stale` queued jobs. | Keep for now; investigate queued job label/routing first. |

## Candidate Remove

None in the latest audit.

## Approval Needed

No deletion approval is pending.
