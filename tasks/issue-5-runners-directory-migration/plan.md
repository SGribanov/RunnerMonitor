<!-- Issue: SGribanov/RunnerMonitor#5 -->
# Runner Directory Migration Plan

## Goal

Move active GitHub Actions runner folders into common `Runners` directories so
the current workstation and future dedicated runner host have predictable
layout and easier discovery.

## Target layout

Windows:

```text
C:\Runners\<owner>-<repo>\<runner-name>
```

WSL/Linux on the current workstation:

```text
/home/gsv777/Runners/<owner>-<repo>/<runner-name>
```

Future dedicated Linux host:

```text
/opt/Runners/<owner>-<repo>/<runner-name>
```

## Current migration table

| Repo | Runner | Current path | Target path | Control mode | Move now? |
|---|---|---|---|---|---|
| `SGribanov/AU` | `windows-local` | `C:\actions-runner` | `C:\Runners\SGribanov-AU\windows-local` | manual Windows | later, when AU work resumes |
| `SGribanov/BackTester` | `backtester-runner` | `C:\actions-runner-backtester` | `C:\Runners\SGribanov-BackTester\backtester-runner` | manual Windows | yes, if idle |
| `SGribanov/DeltaG` | `deltag-win` | `C:\github-runners\deltag` | `C:\Runners\SGribanov-DeltaG\deltag-win` | Windows service | no, currently busy |
| `SGribanov/IdeaBox` | `ideabox-runner` | `C:\actions-runner-ideabox` | `C:\Runners\SGribanov-IdeaBox\ideabox-runner` | Windows service | yes, if idle/admin available |
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-windows-local` | `C:\actions-runner-mycloneosengine` | `C:\Runners\SGribanov-MyCloneOsEngine\mycloneosengine-windows-local` | manual Windows | yes, if idle |
| `SGribanov/DeltaG` | `deltag-linux-wsl` | `/home/gsv777/actions-runner-deltag` | `/home/gsv777/Runners/SGribanov-DeltaG/deltag-linux-wsl` | WSL systemd | no, currently busy |
| `SGribanov/MyCloneOsEngine` | `mycloneosengine-linux` | `/home/gsv777/myclone-runner-linux` | `/home/gsv777/Runners/SGribanov-MyCloneOsEngine/mycloneosengine-linux` | WSL systemd | yes, if idle/sudo available |

## Per-runner migration sequence

1. Confirm explicit approval naming repo, runner, current path, and target path.
2. Confirm GitHub `busy=false` for the runner unless the user explicitly
   approves force handling.
3. Back up the runner folder to `C:\Runners-backup` or
   `/home/gsv777/runner-backups`.
4. Stop only the runner being moved.
5. Move the folder to the target path.
6. Reconfigure path-bound control:
   - manual Windows: update RunnerMonitor audit docs; no service update needed.
   - Windows service: uninstall/reinstall service or re-register with
     `config.cmd --replace` if needed.
   - WSL systemd: uninstall/reinstall `svc.sh` from the new path or re-register
     with `config.sh --replace` if needed.
7. Start the runner from the new path.
8. Validate:
   - `runner-monitor --audit`
   - GitHub runner status is `online`
   - project-scoped `--start-current` returns `already running` from the
     corresponding project root.

## Safety rules

- Do not move busy runners.
- Do not move more than one runner per approval.
- Do not delete old folders until the moved runner is online from the new path
  and a rollback backup exists.
- Do not modify GitHub registrations unless the specific runner migration plan
  requires it.
