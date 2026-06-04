# Admin Window Rendering Status

## 2026-06-04

- Created issue #15 for the elevated-window rendering regression.
- Root cause hypothesis: the TUI was running in the normal terminal buffer and also forced small terminal heights to a 12-line layout. In an elevated PowerShell window this can leave only the bottom of the redraw visible.
- Implemented `tea.WithAltScreen()` for interactive TUI startup.
- Added compact rendering for terminal heights of 8 lines or less.
