# Security Policy

RunnerMonitor interacts with CI runner registrations, local services, WSL sudo,
and remote hosts. Treat configuration and runner tokens as sensitive.

## Supported Versions

RunnerMonitor is currently maintained from the `main` branch and active
feature branches in this repository.

## Reporting A Vulnerability

Do not publish secrets, runner tokens, sudo passwords, private hostnames, or
private infrastructure details in a public issue.

For this repository, report sensitive findings directly to the maintainer
through a private GitHub channel or another trusted private contact path.

## Secret Handling

- Keep `runner-monitor.json` beside the executable and out of git.
- Keep `wslSudoPassword` only in the app-local config file.
- Never paste GitHub runner registration/remove tokens into issues, PRs, or
  committed files.
- Redact local hostnames and paths when they reveal private infrastructure.

## Safety Expectations

- Destructive operations must remain dry-run by default.
- Folder deletion must stay limited to configured runner roots.
- Busy runner protection must remain enabled unless an explicit force command is
  used by the operator.
