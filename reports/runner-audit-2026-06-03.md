# Runner Audit -- 2026-06-03

| Field | Value |
|---|---|
| Project | RunnerMonitor |
| Source | `runner-monitor --audit` |
| Status | read-only proposal |

## Summary

No runner was removed during this audit. The table below separates runners into
`keep`, `investigate`, and `candidate-remove`. Removal commands are provided
only for candidate runners and must not be executed without explicit approval.

## Keep

| Repo | Runner | Host | Evidence |
|---|---|---|---|
| `SGribanov/IdeaBox` | `ideabox-runner` | local Windows | Service-managed and GitHub online. |

## Investigate Before Any Cleanup

| Repo | Runner | Host | Evidence | Suggested action |
|---|---|---|---|---|
| `SGribanov/BackTester` | `backtester-runner` | local Windows | GitHub online, but not service-managed locally. | Decide whether to install as service or remove/recreate later. |
| `SGribanov/DeltaG` | `deltag-win` | local Windows | Running, GitHub online, repo has `1/1 stale` queued jobs. | Keep for now; investigate queued job label/routing first. |
| `SGribanov/DeltaG` | `legion-windows-x64` | local Windows | GitHub online, manual local directory, repo has `1/1 stale` queued jobs. | Decide whether this is the desired Windows DeltaG runner or legacy. |
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-windows-local` | local Windows | GitHub online, but not service-managed locally. | Decide whether it should become service-managed. |
| `SGribanov/DeltaG` | `deltag-linux-wsl` | WSL Ubuntu | Active, GitHub online, repo has `1/1 stale` queued jobs. | Keep for now; investigate queued job label/routing first. |
| `SGribanov/DeltaG` | `legion-ubuntu-wsl-x64` | WSL Ubuntu | Local configured runner not visible in GitHub API, but DeltaG has queued jobs. | Investigate labels before removal. |

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

### `SGribanov/NewGenOsEngine` -- `newgenosengine-windows-local`

Path: `C:\actions-runner-newgenosengine`

```powershell
New-Item -ItemType Directory -Force -Path C:\Runners-backup | Out-Null
Compress-Archive -LiteralPath C:\actions-runner-newgenosengine -DestinationPath C:\Runners-backup\actions-runner-newgenosengine-windows-local.zip -Force
Get-Service | Where-Object { $_.Name -like '*NewGenOsEngine*' -or $_.Name -like '*newgenosengine-windows-local*' }
Remove-Item -LiteralPath C:\actions-runner-newgenosengine -Recurse -Force
```

### `SGribanov/MyCloneOsEngine` -- `mycloneosengine-linux`

Path: `/home/gsv777/myclone-runner-linux`

```powershell
wsl.exe sh -lc 'mkdir -p ~/runner-backups && tar -czf ~/runner-backups/myclone-runner-linux.tar.gz -C /home/gsv777 myclone-runner-linux'
wsl.exe sh -lc 'systemctl list-unit-files "actions.runner.*.service" --no-pager | grep -i myclone || true'
wsl.exe sh -lc 'rm -rf /home/gsv777/myclone-runner-linux'
```

### `SGribanov/NewGenOsEngine` -- `newgen-wsl-linux`

Path: `/home/gsv777/newgen-runner`

```powershell
# Service exists and is inactive; disable/remove service only after approval.
wsl.exe sudo systemctl disable actions.runner.SGribanov-NewGenOsEngine.newgen-wsl-linux.service
wsl.exe sudo rm -f /etc/systemd/system/actions.runner.SGribanov-NewGenOsEngine.newgen-wsl-linux.service
wsl.exe sudo systemctl daemon-reload

wsl.exe sh -lc 'mkdir -p ~/runner-backups && tar -czf ~/runner-backups/newgen-runner.tar.gz -C /home/gsv777 newgen-runner'
wsl.exe sh -lc 'rm -rf /home/gsv777/newgen-runner'
```

## Approval Needed

Deletion requires explicit approval naming each runner:

- `SGribanov/AU windows-local C:\actions-runner`
- `SGribanov/NewGenOsEngine newgenosengine-windows-local C:\actions-runner-newgenosengine`
- `SGribanov/MyCloneOsEngine mycloneosengine-linux /home/gsv777/myclone-runner-linux`
- `SGribanov/NewGenOsEngine newgen-wsl-linux /home/gsv777/newgen-runner`
