<!-- Issue: SGribanov/RunnerMonitor#27 -->
# Plans

## Source
- Task: Make update notice release URL clickable.
- Canonical input: GitHub issue #27.
- Repo context: TUI update notice rendering in `internal/app/model.go`.
- Last updated: 2026-06-05

## Assumptions
- Terminals that support OSC-8 should expose the release URL as a clickable link.
- Terminals that do not support OSC-8 should still display a readable URL label.

## Milestone Order
| ID | Title | Depends on | Status |
| --- | --- | --- | --- |
| M1 | Render clickable update URL | - | [x] |
| M2 | Regression coverage and research note | M1 | [x] |

## M1. Render clickable update URL `[x]`
### Goal
- Update notices preserve the visible URL while wrapping it in a terminal hyperlink sequence.

### Tasks
- [x] Detect release URL in the notice suffix.
- [x] Truncate only the visible URL label before adding OSC-8 escape sequences.
- [x] Keep non-URL notices on the existing truncation path.

### Definition of Done
- Update notice URL remains readable and clickable where supported.

### Validation
```sh
go test ./...
```

### Known Risks
- Some terminals ignore OSC-8, so readability of the visible URL remains important.

### Stop-and-Fix Rule
- If tests fail, fix before considering the issue ready.

## M2. Regression coverage and research note `[x]`
### Goal
- Capture the behavior so future truncation changes do not break terminal links.

### Tasks
- [x] Add unit test for OSC-8 update URL rendering.
- [x] Update technology insights.

### Definition of Done
- Regression test proves the escape sequence and visible URL are present.

### Validation
```sh
go test ./...
```

### Known Risks
- Visual click behavior is terminal-dependent and not fully exercised by unit tests.

### Stop-and-Fix Rule
- If test coverage regresses, restore the clickable and readable URL behavior.
