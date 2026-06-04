# Startup Update Check Test Plan

- Run `go test ./...`.
- Run `powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1`.
- Manual TUI check with internet: start TUI and confirm no error is shown when current version is latest.
- Manual future-release check: after publishing a newer version, start an older binary and confirm the update notice appears.
