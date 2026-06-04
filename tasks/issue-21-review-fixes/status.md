<!-- Issue: SGribanov/RunnerMonitor#21 -->
# Review Fixes Status

## 2026-06-04

Current phase: release publishing.

Done:
- Full review completed and issue #21 reopened for fixes.
- Hardened WSL/Linux and Windows safe-root checks.
- Moved runner token passing through env references where practical and redacted token text from command errors.
- Added 30-second GitHub status cache for automatic TUI refreshes; manual refresh remains fresh.
- Quoted remote TUI command paths.
- Reworked runner sorting to precompute keys.
- Updated README, README_RU, and technology insights.

- Bumped minor release metadata to `v0.3.0`.
- Built `dist/RunnerMonitor-v0.3.0-windows-x64.zip` and SHA256.

In progress:
- Commit, push, merge, publish GitHub release, and final cleanup.

Validation notes:
- `go test ./...`, `go vet ./...`, and `govulncheck ./...` passed before implementation.
- `go test -race ./...` is blocked locally by MinGW assembler failure with the available GCC toolchain.
- Post-implementation `go test ./...`, `go vet ./...`, `govulncheck ./...`, `scripts/build.ps1`, and `git diff --check` passed.
- Release validation on 2026-06-04: `go test ./...`, `go vet ./...`, `govulncheck ./...`, and `git diff --check` passed after the `v0.3.0` version bump.
