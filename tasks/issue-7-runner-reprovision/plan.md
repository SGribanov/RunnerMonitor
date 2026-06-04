<!-- Issue: SGribanov/RunnerMonitor#7 -->
# Runner Remove And Reprovision Plan

## Scope

Add a safe runner lifecycle workflow for full runner removal and project-folder
based reprovisioning. A project is specified by folder name under `C:\Repos`,
then RunnerMonitor resolves `C:\Repos\<ProjectName>` to the GitHub `owner/repo`
remote.

## Milestone 1 -- Safe Removal Core `[x]`

Goal: unregister and locally retire a selected runner without touching busy
work or unrelated folders.

Tasks:
- [x] Resolve a runner by TUI number or CLI runner name.
- [x] Refuse removal when GitHub reports `busy=true` unless an explicit force
  flag is used.
- [x] Stop the selected runner when needed.
- [x] Run official `config.cmd remove` or `config.sh remove` with a GitHub
  remove token.
- [x] Uninstall Windows service or WSL systemd service when the runner is
  service-managed.
- [x] Delete or quarantine the runner folder only with an explicit
  `--delete-folder` / `delete` confirmation path and only inside known runner
  roots.

Definition of done:
- CLI removal can dry-run and execute a selected idle runner.
- TUI exposes a clear command form for selected-runner removal.
- Tests cover busy refusal, project/path safety, and command construction.

Validation:
- `go test ./...`
- `runner-monitor.ps1 --help`
- dry-run remove command against a known runner

Risks:
- Windows services require elevation; non-elevated TUI must launch the existing
  UAC helper style or return a precise instruction.
- WSL service removal may require sudo fallback.

## Milestone 2 -- Project Folder Resolution `[x]`

Goal: let the operator specify a project by local folder name.

Tasks:
- [x] Resolve `PROJECT` to `C:\Repos\PROJECT`.
- [x] Read the GitHub repo from that folder's `origin`.
- [x] Use the resolved repo for token generation and runner naming.
- [x] Return precise errors for missing folder, missing git remote, or non-GitHub
  remote.

Definition of done:
- CLI commands accept `--project PROJECT`.
- Unit tests cover project folder resolution logic without depending on real
  repos.

Validation:
- `go test ./...`

## Milestone 3 -- Reconfigure Existing Runner Folder `[x]`

Goal: add a new runner by configuring an existing runner distribution folder,
without trying to clone binary trees in v1.

Tasks:
- [x] Accept runner folder path, project folder name, runner name, and labels.
- [x] Require the folder to contain `config.cmd` or `config.sh`.
- [x] Refuse to overwrite an already configured folder unless `--replace` is
  explicit.
- [x] Fetch a GitHub registration token with `gh api`.
- [x] Run official unattended config command for Windows or WSL/Linux.
- [x] Optionally install/start service when requested.

Definition of done:
- CLI can print a dry-run plan and execute the config command with explicit
  confirmation.
- TUI can show the intended command path or delegate to CLI text for now.

Validation:
- `go test ./...`
- dry-run add command for a real project folder

Risks:
- Cloning an existing runner folder can break junctions/symlinks after runner
  auto-update; v1 should prefer configuring an already prepared distribution
  folder.

## Milestone 4 -- Clone Template Folder `[ ]`

Goal: later support cloning an existing runner folder into the standard
`C:\Runners\<Owner-Repo>\<RunnerName>` or WSL equivalent.

Tasks:
- [ ] Preserve runner distribution files but exclude registration and work
  state: `.runner`, `.credentials*`, `.service`, `_work`, `_diag`.
- [ ] Handle Windows junctions and WSL symlinks deliberately.
- [ ] Validate the cloned folder can run `config --check`.

Definition of done:
- Clone path is tested against current local runner folder structure.
- A failed clone leaves either no destination folder or a clearly quarantined
  partial folder.

Validation:
- `go test ./...`
- manual dry run on a disposable destination

## Stop-And-Fix Rules

- Stop before deleting any folder outside known runner roots.
- Stop before unregistering a busy runner unless the user explicitly forces it.
- Stop if GitHub token generation fails.
- Stop if service uninstall fails; do not delete the folder afterward.
