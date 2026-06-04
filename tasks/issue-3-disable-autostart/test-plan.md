<!-- Issue: SGribanov/RunnerMonitor#3 -->
# Disable Runner Autostart Test Plan

## Automated

- `go test ./...`
- `go run ./cmd/runner-monitor --start-repo SGribanov/NoSuchRepo`
- `powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --start-current`
- `powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\scripts\install-prepush-hook.ps1 -RepoPath <test-repo>`

## Manual

- Verify Windows services show `StartMode: Manual`.
- Verify WSL units show `disabled`.
- Verify existing running runners are not stopped by disabling autostart.
- Verify `runner-monitor --start-repo SGribanov/DeltaG` starts service-managed DeltaG runners in the background.
- Verify `--start-current` is invoked from the target project root, not from `RunnerMonitor`.
- Verify pre-push hook installation only modifies the selected repository.
