<!-- Issue: SGribanov/RunnerMonitor#7 -->
# Runner Remove And Reprovision Test Plan

## Unit Tests

- Project folder resolution:
  - existing `C:\Repos\<Project>` resolves to a GitHub repo;
  - missing folder returns a clear error;
  - non-GitHub origin returns a clear error.
- Removal safety:
  - busy runner is refused without force;
  - unknown runner name reports not found;
  - folder deletion is refused outside known runner roots;
  - dry-run returns planned steps without executing commands.
- Command construction:
  - Windows runner uses `config.cmd remove`;
  - WSL/Linux runner uses `config.sh remove`;
  - registration/remove token API endpoints are correct;
  - PowerShell quoting remains safe for paths and runner names.

## CLI Smoke

- `runner-monitor.ps1 --help`
- `runner-monitor.ps1 --remove-runner definitely-not-a-runner`
- `runner-monitor.ps1 --remove-runner ideabox-runner --repo SGribanov/IdeaBox`
- `runner-monitor.ps1 --project RunnerMonitor --add-runner test-runner --runner-folder <prepared-runner-folder>`

## Manual Safety Checks

- Before executing removal, run `runner-monitor.ps1 --audit` and verify
  `busy=false` for the selected runner.
- For Windows service runners, confirm UAC elevation is expected before service
  uninstall or folder deletion.
- For WSL runners, confirm sudo password file is readable and never printed.

## Acceptance Gates

- No command deletes a runner folder unless an explicit delete option is present.
- No command unregisters a busy runner unless an explicit force option is present.
- No command unregisters or configures a runner unless `--confirm` is present.
- Dry-run output includes repo, runner name, runner path, service name, and
  exact high-level steps.
- Tests pass with `go test ./...`.
