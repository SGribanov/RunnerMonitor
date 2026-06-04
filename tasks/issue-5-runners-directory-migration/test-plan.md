<!-- Issue: SGribanov/RunnerMonitor#5 -->
# Runner Directory Migration Test Plan

## Static validation

```powershell
go test ./...
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --audit
```

## Per-runner validation

After each approved move:

```powershell
gh api repos/<owner>/<repo>/actions/runners
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --audit
```

From the corresponding project root:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --start-current
```

Expected:

- moved runner is discovered at the new `Runners` path;
- GitHub status is `online`;
- local state is `running` or `active`;
- `--start-current` returns `already running` for an already active runner;
- no unrelated runner was stopped or moved.
