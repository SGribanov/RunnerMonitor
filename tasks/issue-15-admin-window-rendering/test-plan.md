# Admin Window Rendering Test Plan

- Run `go test ./...`.
- Run `go build ./cmd/runner-monitor`.
- Manual Windows check: launch `runner-monitor.ps1` normally and from elevated PowerShell.
- Resize the elevated terminal to a very small height and confirm the title, status, and input remain visible.
