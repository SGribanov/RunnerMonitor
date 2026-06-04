<!-- Issue: SGribanov/RunnerMonitor#23 -->
# Window Exit Plan

Goal: prevent Windows discovery PowerShell subprocesses from keeping the TUI launch window alive after exit.

- [x] Create and link GitHub issue #23.
- [x] Add a bounded timeout around Windows discovery PowerShell calls.
- [x] Add coverage for discovery timeout warnings.
- [x] Run validation checks.
- [ ] Publish GitHub handoff after merge.
