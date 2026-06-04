<!-- Issue: SGribanov/RunnerMonitor#7 -->
# Runner Remove And Reprovision Status

## Current Phase

Milestones 1-3 implemented; Milestone 4 template cloning remains deferred.

## Done

- Created GitHub issue #7 for runner removal and reprovisioning.
- Created branch `codex/7-runner-reprovision`.
- Confirmed current implementation already has discovery, GitHub status merge,
  service/manual lifecycle control, cleanup safety, WSL sudo fallback, and
  Windows UAC helper patterns to reuse.
- Refreshed GitHub Actions runner token/config command behavior against
  official GitHub documentation:
  - repository registration/remove tokens are created with REST API endpoints;
  - tokens expire after one hour;
  - labels can be assigned during initial runner configuration.
- Added CLI flags:
  - `--remove-runner NAME`
  - `--add-runner NAME`
  - `--project PROJECT`
  - `--repo OWNER/REPO`
  - `--runner-folder PATH`
  - `--confirm`
  - `--force`
  - `--delete-folder`
  - `--replace`
  - `--service`
- Added TUI commands:
  - `remove N`
  - `remove N confirm`
  - `delete N confirm`
- Added dry-run default for remove/add operations.
- Added project-folder resolution from `C:\Repos\<Project>` to GitHub origin.
- Added safety gates for busy runners, explicit folder deletion, known runner
  roots, configured-folder replacement, and path-like project input.

## In Progress

- Final validation and handoff.

## Next

- Run final `go test ./...`.
- Rebuild `bin\runner-monitor.exe`.
- Run CLI smoke commands.
- Push branch and publish issue handoff.

## Decisions

- V1 should not blindly clone existing runner folders because current runner
  folders may contain Windows junctions or WSL symlinks created by runner
  auto-update.
- V1 will support adding by configuring an existing prepared runner distribution
  folder. Template cloning is a later milestone after symlink/junction handling
  is tested.
- Project input means folder name under `C:\Repos`, not a GitHub repo string.

## Assumptions

- `gh` is authenticated and has repository administration/actions-runner access.
- Standard runner folders contain `config.cmd` on Windows or `config.sh` on WSL
  and Linux.
- Windows service operations may need elevation; WSL service operations may need
  sudo through the existing password-file fallback.

## Validation Log

- Passed: `go test ./...`
- Passed: `runner-monitor.ps1 --help`
- Passed: `runner-monitor.ps1 --remove-runner ideabox-runner --repo SGribanov/IdeaBox`
- Passed: `runner-monitor.ps1 --add-runner runner-monitor-test --project RunnerMonitor --runner-folder <temp prepared folder> --labels "self-hosted,Windows,X64"`
- Passed: configured runner folder is rejected without `--replace`.
- Passed: rebuilt `bin\runner-monitor.exe` with `scripts\build.ps1`.
