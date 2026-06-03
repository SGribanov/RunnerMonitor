<!-- Issue: SGribanov/RunnerMonitor#2 -->
# Runner Audit And Cleanup Plan

## Milestone 1: Read-only audit

Status: [x]

Goal: classify discovered runners before any deletion.

Tasks:
- [x] Add `--audit`.
- [x] Classify runners as `keep`, `investigate`, or `candidate-remove`.
- [x] Include queue impact and local/GitHub state as evidence.

Validation:
- `go test ./...`
- `go run ./cmd/runner-monitor --audit`

## Milestone 2: Approved removal plan

Status: [ ]

Goal: remove only explicitly approved old runner services, GitHub registrations, and folders.

Tasks:
- [ ] Review audit output with the operator.
- [ ] Produce exact removal commands per approved runner.
- [ ] Remove service/unit first, then GitHub registration, then folder backup/delete.

Safety rule:
- No deletion without explicit runner-by-runner approval.

