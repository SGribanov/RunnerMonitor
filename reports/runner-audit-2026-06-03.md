# Runner Audit -- 2026-06-03

| Field | Value |
|---|---|
| Project | RunnerMonitor |
| Source | `runner-monitor --audit` |
| Status | partially executed |

## Summary

One runner was removed after explicit approval: `legion-ubuntu-wsl-x64`.
Remaining runners are separated into `keep`, `investigate`, and
`candidate-remove`. Removal commands are provided only for candidate runners and
must not be executed without explicit approval.

## Completed Removals

| Repo | Runner | Path | Backup |
|---|---|---|---|
| `SGribanov/DeltaG` | `legion-ubuntu-wsl-x64` | `/home/gsv777/actions-runner-linux-x64` | `/home/gsv777/runner-backups/actions-runner-linux-x64-legion-ubuntu-wsl-x64-2026-06-03.tar.gz` |
| `SGribanov/DeltaG` | `legion-windows-x64` | `C:\actions-runner-win-x64` | `C:\Runners-backup\actions-runner-win-x64-legion-windows-x64-2026-06-03.zip` |
| `SGribanov/NewGenOsEngine` | `newgenosengine-windows-local` | `C:\actions-runner-newgenosengine` | `C:\Runners-backup\actions-runner-newgenosengine-windows-local-2026-06-03.zip` |
| `SGribanov/NewGenOsEngine` | `newgen-wsl-linux` | `/home/gsv777/newgen-runner` | `/home/gsv777/runner-backups/newgen-runner-newgen-wsl-linux-2026-06-03.tar.gz` |

## Remaining Elevated Cleanup

The `newgen-wsl-linux` runner folder was removed, but its systemd unit remains
enabled because `sudo -n` reported that a password is required.

```powershell
wsl.exe sudo systemctl disable actions.runner.SGribanov-NewGenOsEngine.newgen-wsl-linux.service
wsl.exe sudo rm -f /etc/systemd/system/actions.runner.SGribanov-NewGenOsEngine.newgen-wsl-linux.service
wsl.exe sudo systemctl daemon-reload
```

## Keep

| Repo | Runner | Host | Evidence |
|---|---|---|---|
| `SGribanov/IdeaBox` | `ideabox-runner` | local Windows | Service-managed and GitHub online. |

## Investigate Before Any Cleanup

| Repo | Runner | Host | Evidence | Suggested action |
|---|---|---|---|---|
| `SGribanov/BackTester` | `backtester-runner` | local Windows | GitHub online, but not service-managed locally. | Decide whether to install as service or remove/recreate later. |
| `SGribanov/DeltaG` | `deltag-win` | local Windows | Running, GitHub online, repo has `1/1 stale` queued jobs. | Keep for now; investigate queued job label/routing first. |
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-windows-local` | local Windows | GitHub online, but not service-managed locally. | Decide whether it should become service-managed. |
| `SGribanov/DeltaG` | `deltag-linux-wsl` | WSL Ubuntu | Active, GitHub online, currently busy, repo has `2/1 stale` queued jobs. | Keep for now; investigate queued job label/routing first. |

## Candidate Remove

These runners are local configured/manual or inactive and are not visible in the
current GitHub API status. The commands below are intentionally separated into
backup, service/GitHub registration check, and deletion steps.

### `SGribanov/AU` -- `windows-local`

Path: `C:\actions-runner`

```powershell
# Backup first
New-Item -ItemType Directory -Force -Path C:\Runners-backup | Out-Null
Compress-Archive -LiteralPath C:\actions-runner -DestinationPath C:\Runners-backup\actions-runner-au-windows-local.zip -Force

# Confirm no service exists
Get-Service | Where-Object { $_.Name -like 'actions.runner.*AU*' -or $_.Name -like '*windows-local*' }

# Delete only after approval
Remove-Item -LiteralPath C:\actions-runner -Recurse -Force
```

### `SGribanov/MyCloneOsEngine` -- `mycloneosengine-linux`

Path: `/home/gsv777/myclone-runner-linux`

```powershell
wsl.exe sh -lc 'mkdir -p ~/runner-backups && tar -czf ~/runner-backups/myclone-runner-linux.tar.gz -C /home/gsv777 myclone-runner-linux'
wsl.exe sh -lc 'systemctl list-unit-files "actions.runner.*.service" --no-pager | grep -i myclone || true'
wsl.exe sh -lc 'rm -rf /home/gsv777/myclone-runner-linux'
```

## Approval Needed

Deletion requires explicit approval naming each runner:

- `SGribanov/AU windows-local C:\actions-runner`
- `SGribanov/MyCloneOsEngine mycloneosengine-linux /home/gsv777/myclone-runner-linux`
