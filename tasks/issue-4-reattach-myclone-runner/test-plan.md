<!-- Issue: SGribanov/RunnerMonitor#4 -->
# Reattach MyClone WSL Runner Test Plan

## Automated

- `go test ./...`

## Manual

- Verify GitHub API lists `mycloneosengine-linux` as `online`.
- Verify WSL systemd unit is `active`.
- Verify audit reports AU and MyClone Linux as `keep`.

