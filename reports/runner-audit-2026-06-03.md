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
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-linux` | `/home/gsv777/Runners/SGribanov-MyCloneOsEngine/mycloneosengine-linux` | Reconfigured with GitHub runner id `24`; later moved from `/home/gsv777/myclone-runner-linux`; backup `/home/gsv777/runner-backups/myclone-runner-linux-mycloneosengine-linux-move-2026-06-03.tar.gz`; fixed `bin`/`externals` symlinks; systemd unit active/enabled from new path. |
| `SGribanov/BackTester` | `backtester-runner` | `C:\Runners\SGribanov-BackTester\backtester-runner` | Moved from `C:\actions-runner-backtester`; backup `C:\Runners-backup\actions-runner-backtester-backtester-runner-2026-06-03.zip`; fixed `bin`/`externals` junctions; GitHub online. |
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-windows-local` | `C:\Runners\SGribanov-MyCloneOsEngine\mycloneosengine-windows-local` | Moved from `C:\actions-runner-mycloneosengine`; backup `C:\Runners-backup\actions-runner-mycloneosengine-mycloneosengine-windows-local-2026-06-03.zip`; fixed `bin`/`externals` junctions; GitHub online. |
| `SGribanov/DeltaG` | `deltag-linux-wsl` | `/home/gsv777/Runners/SGribanov-DeltaG/deltag-linux-wsl` | Moved from `/home/gsv777/actions-runner-deltag`; backup `/home/gsv777/runner-backups/actions-runner-deltag-deltag-linux-wsl-move-2026-06-03.tar.gz`; fixed `bin`/`externals` symlinks; systemd unit active/enabled from new path; GitHub online. |
| `SGribanov/AU` | `windows-local` | `C:\Runners\SGribanov-AU\windows-local` | Moved from `C:\actions-runner`; backup `C:\Runners-backup\actions-runner-windows-local-move-2026-06-03.zip`; fixed `bin`/`externals` junctions; reattached GitHub runner binding; GitHub online. |

## Resolved Elevated Cleanup

The `newgen-wsl-linux` runner folder was removed first. Its systemd unit was
then removed manually from WSL with sudo. Current WSL unit list no longer shows
`actions.runner.SGribanov-NewGenOsEngine.newgen-wsl-linux.service`.

## Keep

| Repo | Runner | Host | Evidence |
|---|---|---|---|
| `SGribanov/AU` | `windows-local` | local Windows | Manual `run.cmd` runner is running from `C:\Runners\SGribanov-AU\windows-local`, GitHub online, and controllable by RunnerMonitor. |
| `SGribanov/BackTester` | `backtester-runner` | local Windows | Manual `run.cmd` runner is running from `C:\Runners\SGribanov-BackTester\backtester-runner`, GitHub online, and controllable by RunnerMonitor. |
| `SGribanov/IdeaBox` | `ideabox-runner` | local Windows | Service-managed and GitHub online. |
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-windows-local` | local Windows | Manual `run.cmd` runner is running from `C:\Runners\SGribanov-MyCloneOsEngine\mycloneosengine-windows-local`, GitHub online, and controllable by RunnerMonitor. |
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-linux` | WSL Ubuntu | Service-managed from `/home/gsv777/Runners/SGribanov-MyCloneOsEngine/mycloneosengine-linux`, systemd active, GitHub online. |
| `SGribanov/DeltaG` | `deltag-linux-wsl` | WSL Ubuntu | Service-managed from `/home/gsv777/Runners/SGribanov-DeltaG/deltag-linux-wsl`, systemd active, GitHub online, and controllable by RunnerMonitor. |

## Investigate Before Any Cleanup

| Repo | Runner | Host | Evidence | Suggested action |
|---|---|---|---|---|
| `SGribanov/DeltaG` | `deltag-win` | local Windows | Running, GitHub online, repo has `1/1 stale` queued jobs. | Keep for now; investigate queued job label/routing first. |
| `SGribanov/DeltaG` | `deltag-linux-wsl` | WSL Ubuntu | Active from the new `Runners` path, GitHub online, repo has `1/1 stale` queued jobs. | Keep; investigate queued job label/routing separately from folder migration. |

### DeltaG Stale Queue Detail

The remaining queued DeltaG run is `26447257991`, workflow `ci`, display title
`research(margin): diagnose vertical freshness drift`, event `pull_request`,
branch `codex/604-vertical-freshness-diagnostic`, created
`2026-05-26T10:36:43Z`, updated `2026-05-26T10:52:39Z`.
GitHub reports zero jobs for the run. Related PRs `#709` and `#710` for the same
branch are both closed. At diagnosis time, DeltaG self-hosted runners
`deltag-win` and `deltag-linux-wsl` were online and not busy, so this looks like
a stale GitHub Actions run rather than missing local runner capacity.

An explicit cancel attempt was made with
`POST /repos/SGribanov/DeltaG/actions/runs/26447257991/cancel`. GitHub returned
HTTP `409` with `Cannot cancel a workflow re-run that has not yet queued`; a
follow-up `gh run view` still reports the run as `queued`.

## Candidate Remove

None in the latest audit.

## Approval Needed

No deletion approval is pending.

## Elevated Windows Service Blocker

The remaining unmoved runner folders are Windows service-managed:

- `SGribanov/IdeaBox ideabox-runner` at `C:\actions-runner-ideabox`
- `SGribanov/DeltaG deltag-win` at `C:\github-runners\deltag`

Both services are currently running with `StartMode=Auto`, and GitHub reports
both runners online/busy=false. The current shell is not elevated
(`Administrator=False`). Attempting to disable autostart with `sc.exe config
... start= demand` returned `OpenService FAILED 5: Access is denied` for both
services. Moving these folders safely requires elevated PowerShell so the
service can be stopped, reconfigured to the new `C:\Runners` path, switched to
manual startup, and restarted.
