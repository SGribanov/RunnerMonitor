package app

import (
	"fmt"
	"strings"
)

func renderTable(runners []Runner) string {
	if len(runners) == 0 {
		return "No runners discovered. Check GitHub CLI auth and scan roots.\n"
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%-3s %-10s %-24s %-28s %-9s %-10s %-7s %-11s %-20s %s\n",
		"#", "Host", "Repo", "Runner", "Local", "GitHub", "Busy", "Queue", "Labels", "Path"))
	b.WriteString(strings.Repeat("-", 150))
	b.WriteString("\n")

	for i, r := range runners {
		busy := "false"
		if r.Busy {
			busy = "true"
		}
		queue := fmt.Sprintf("%d", r.QueueCount)
		if r.StaleQueueCount > 0 {
			queue = fmt.Sprintf("%d/%d stale", r.QueueCount, r.StaleQueueCount)
		}
		b.WriteString(fmt.Sprintf("%-3d %-10s %-24s %-28s %-9s %-10s %-7s %-11s %-20s %s\n",
			i+1,
			trunc(r.Host, 10),
			trunc(r.Repo, 24),
			trunc(r.Name, 28),
			trunc(r.LocalState, 9),
			trunc(r.GitHubStatus, 10),
			busy,
			trunc(queue, 11),
			trunc(strings.Join(r.Labels, ","), 20),
			r.Path,
		))
	}
	return b.String()
}

func RenderInventory(inventory Inventory) string {
	return renderTable(inventory.Runners)
}

func RenderAudit(inventory Inventory) string {
	if len(inventory.Runners) == 0 {
		return "No runners discovered.\n"
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%-3s %-17s %-24s %-28s %-9s %-10s %-7s %-11s %s\n",
		"#", "Decision", "Repo", "Runner", "Local", "GitHub", "Busy", "Queue", "Evidence"))
	b.WriteString(strings.Repeat("-", 150))
	b.WriteString("\n")

	for i, r := range inventory.Runners {
		decision, evidence := AuditRunner(r)
		busy := "false"
		if r.Busy {
			busy = "true"
		}
		queue := fmt.Sprintf("%d", r.QueueCount)
		if r.StaleQueueCount > 0 {
			queue = fmt.Sprintf("%d/%d stale", r.QueueCount, r.StaleQueueCount)
		}
		b.WriteString(fmt.Sprintf("%-3d %-17s %-24s %-28s %-9s %-10s %-7s %-11s %s\n",
			i+1,
			trunc(decision, 17),
			trunc(r.Repo, 24),
			trunc(r.Name, 28),
			trunc(r.LocalState, 9),
			trunc(r.GitHubStatus, 10),
			busy,
			trunc(queue, 11),
			evidence,
		))
	}
	return b.String()
}

func trunc(value string, width int) string {
	if len(value) <= width {
		return value
	}
	if width <= 1 {
		return value[:width]
	}
	return value[:width-1] + "…"
}
