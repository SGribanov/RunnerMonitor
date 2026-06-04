# TUI Help Test Plan

- Run `go test ./...`.
- Run `powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build.ps1`.
- Manual TUI check: press `h`, `?`, type `help`, and press `Esc`.
- Confirm normal table navigation still works when help is closed.
