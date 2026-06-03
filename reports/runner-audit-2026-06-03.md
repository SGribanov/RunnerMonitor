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
| `SGribanov/BackTester` | `backtester-runner` | local Windows | Manual `run.cmd` runner is running, GitHub online, and now controllable by RunnerMonitor. |
| `SGribanov/IdeaBox` | `ideabox-runner` | local Windows | Service-managed and GitHub online. |
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-windows-local` | local Windows | Manual `run.cmd` runner is running, GitHub online, and now controllable by RunnerMonitor. |
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-linux` | WSL Ubuntu | Service-managed, systemd active, GitHub online. |

## Investigate Before Any Cleanup

| Repo | Runner | Host | Evidence | Suggested action |
|---|---|---|---|---|
| `SGribanov/DeltaG` | `deltag-win` | local Windows | Running, GitHub online, repo has `1/1 stale` queued jobs. | Keep for now; investigate queued job label/routing first. |
| `SGribanov/DeltaG` | `deltag-linux-wsl` | WSL Ubuntu | Active, GitHub online, repo has `1/1 stale` queued jobs. | Keep for now; investigate queued job label/routing first. |

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
