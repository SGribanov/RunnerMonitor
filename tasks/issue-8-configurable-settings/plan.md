<!-- Issue: SGribanov/RunnerMonitor#8 -->
# Configurable Settings Plan

## Scope

Make RunnerMonitor portable by moving host-specific paths and WSL sudo
credentials into a user-level settings file. The repository must not contain a
real sudo password.

## Milestone 1 -- Settings Model `[x]`

Goal: add a loaded settings object with safe current-workstation defaults.

Tasks:
- [x] Add `runner-monitor.json` next to the executable.
- [x] Add defaults for `projectsRoot`, `windowsRunnerRoots`, `wslRunnerRoots`,
  `linuxRunnerRoots`, and `wslSudoPassword`.
- [x] Add `--init-config` to create the local config without overwriting by
  default.
- [x] Add `--show-config` with password masking.

Validation:
- `go test ./...`
- `runner-monitor.ps1 --show-config`

## Milestone 2 -- Replace Hard-Coded Roots `[x]`

Goal: use settings in discovery, project resolution, delete safety, and WSL sudo
fallback.

Tasks:
- [x] Windows runner discovery reads `windowsRunnerRoots`.
- [x] WSL runner discovery reads `wslRunnerRoots`.
- [x] Linux discovery reads `linuxRunnerRoots`.
- [x] Project folder resolution reads `projectsRoot`.
- [x] Safe folder deletion checks configured runner roots.
- [x] WSL sudo fallback reads `wslSudoPassword` directly.

Validation:
- `go test ./...`
- `runner-monitor.ps1 --audit`
- dry-run remove/add smoke commands

## Milestone 3 -- Documentation And Handoff `[x]`

Goal: document config location, fields, and security behavior.

Tasks:
- [x] Update README.
- [x] Update research insights in repo and IdeaBox vault.
- [x] Publish issue handoff.

Validation:
- `go test ./...`
- `scripts\build.ps1`

## Safety Rules

- Never commit a real sudo password.
- `--show-config` must never print the real password.
- If config is missing, current workstation defaults must still work.
- If sudo password is missing, WSL sudo fallback must fail with a clear message
  that does not reveal secrets.
