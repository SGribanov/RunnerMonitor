<!-- Issue: SGribanov/RunnerMonitor#10 -->
# GitHub-Ready Documentation Test Plan

## Static Checks

- Confirm `LICENSE` exists at repository root.
- Confirm `README.md` is English and links to `README_RU.MD`.
- Confirm `README_RU.MD` links back to `README.md`.
- Confirm `.gitignore` excludes local config, build outputs, logs, and temporary
  files.
- Confirm community files are present:
  - `CONTRIBUTING.md`
  - `SECURITY.md`
  - `CODE_OF_CONDUCT.md`
  - `.github/ISSUE_TEMPLATE/*`
  - `.github/PULL_REQUEST_TEMPLATE.md`

## Validation Commands

- `go test ./...`
- `scripts\build.ps1`
- `runner-monitor.ps1 --show-config`
- `gh repo view SGribanov/RunnerMonitor --json description,repositoryTopics,licenseInfo`

## Acceptance Gates

- GitHub can detect an open license after push.
- README gives a new user enough context to build, configure, and run the tool.
- Russian README mirrors the important operator guidance.
- No real local secrets are included.
