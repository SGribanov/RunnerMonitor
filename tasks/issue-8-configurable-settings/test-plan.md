<!-- Issue: SGribanov/RunnerMonitor#8 -->
# Configurable Settings Test Plan

## Unit Tests

- Defaults:
  - missing config returns current workstation defaults;
  - default Windows runner root includes `C:\Runners`;
  - default project root is `C:\Repos`;
  - default WSL runner root includes `/home/gsv777/Runners`.
- Config load:
  - custom roots override defaults;
  - empty lists fall back to defaults;
  - `wslSudoPassword` loads as a value.
- Config render:
  - `--show-config` masks non-empty `wslSudoPassword`;
  - empty password renders as `<empty>`.
- Integration helpers:
  - project folder resolution uses configured root;
  - safe delete checks configured runner roots;
  - WSL sudo fallback uses password value and never mentions a password file
    unless legacy environment fallback is used.

## CLI Smoke

- `runner-monitor.ps1 --show-config`
- `runner-monitor.ps1 --init-config`
- `runner-monitor.ps1 --audit`
- `runner-monitor.ps1 --remove-runner ideabox-runner --repo SGribanov/IdeaBox`

## Manual Safety Checks

- Inspect generated config path next to `runner-monitor.exe`.
- Confirm generated config either has an empty password or a locally provided
  value, and no repo file contains the real password.
- Confirm `--show-config` never prints the real password.

## Acceptance Gates

- Tests pass with `go test ./...`.
- Existing runner discovery still finds current Windows/WSL runners.
- No real sudo password is staged or committed.
