# RunnerMonitor -- Strategy Insights

| Field | Value |
|---|---|
| Project | RunnerMonitor |
| Type | strategy-research |
| Last updated | 2026-06-03 |
| Status | active |
| Tags | ci-runners, github-actions, tui, operations |

## 2026-06-03 -- Monitoring goal

RunnerMonitor should not duplicate GitHub's UI. Its useful product boundary is
combining local control state, GitHub runner status, and queue health in one
keyboard-first view. The key operator workflow is short numbered commands such
as `start 1` and `stop 1`, backed by automatic discovery so new runners appear
without per-runner configuration.

## 2026-06-03 -- Remote topology

The accepted operating model is local TUI over SSH. The TUI runs on the
operator machine and controls the future dedicated runner host remotely. This
keeps v1 lightweight and avoids introducing a daemon before OneDev migration
requirements are known.

## 2026-06-03 -- Runner cleanup policy

Runner cleanup must be evidence-driven and explicit. The app should classify
runners as `keep`, `investigate`, or `candidate-remove`, but actual deletion of
services, GitHub registrations, or folders requires separate approval for each
runner.

## 2026-06-03 -- Manual startup policy

Runners should not auto-start at boot. The desired workflow is for Codex or the
operator to start only the required project runner with a short command such as
`runner-monitor --start-repo SGribanov/DeltaG`.
