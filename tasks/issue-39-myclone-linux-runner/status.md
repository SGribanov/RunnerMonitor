<!-- Issue: SGribanov/RunnerMonitor#39 -->
# MyClone Linux Runner Reattach Status

## Done

- Created tracking issue `SGribanov/RunnerMonitor#39`.
- Recreated WSL runner folder:
  `/home/gsv777/Runners/SGribanov-MyCloneOsEngine/mycloneosengine-linux`.
- Installed GitHub Actions runner `2.335.1`.
- Registered `mycloneosengine-linux` to `SGribanov/MyCloneOsEngine` with labels:
  `self-hosted`, `X64`, `Linux`, `mycloneosengine`, `local`.
- Installed and started systemd unit:
  `actions.runner.SGribanov-MyCloneOsEngine.mycloneosengine-linux.service`.
- Verified GitHub runner id `25` is `online` and `busy=false`.
- Verified RunnerMonitor discovers it as `wsl:Ubuntu`, `active`, `online`.

## Notes

- `MyCloneOsEngine` workflow Linux jobs still use `ubuntu-latest`.
- Do not route the Docker job to this WSL runner until Docker availability on
  the runner host is verified.

## Validation

```powershell
gh api repos/SGribanov/MyCloneOsEngine/actions/runners
go run ./cmd/runner-monitor --once
wsl.exe -d Ubuntu -- bash -lc 'systemctl is-active actions.runner.SGribanov-MyCloneOsEngine.mycloneosengine-linux.service'
```
