<!-- Issue: SGribanov/RunnerMonitor#4 -->
# Reattach MyClone WSL Runner Status

## Done

- Created config backup:
  `/home/gsv777/runner-backups/myclone-runner-linux-config-before-reattach-2026-06-03.tar.gz`.
- Reconfigured `mycloneosengine-linux` for `SGribanov/MyCloneOsEngine`.
- Installed and started systemd unit:
  `actions.runner.SGribanov-MyCloneOsEngine.mycloneosengine-linux.service`.
- GitHub runner id `24` is online.
- Audit reports `mycloneosengine-linux` as keep.
- Added `runner-policy.json` keep rule for `SGribanov/AU windows-local`.

## Validation

```powershell
go test ./...
powershell -NoProfile -ExecutionPolicy Bypass -File C:\Repos\RunnerMonitor\runner-monitor.ps1 --audit
gh api repos/SGribanov/MyCloneOsEngine/actions/runners
```

