package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func newRunnerTable(runners []Runner, width int, height int) table.Model {
	styles := table.DefaultStyles()
	styles.Header = styles.Header.Bold(true).Foreground(lipgloss.Color("39"))
	styles.Selected = styles.Selected.Bold(true).Foreground(lipgloss.Color("15")).Background(lipgloss.Color("62"))

	t := table.New(
		table.WithColumns(runnerTableColumns(width)),
		table.WithRows(runnerTableRows(runners)),
		table.WithWidth(width),
		table.WithHeight(height),
		table.WithFocused(true),
		table.WithStyles(styles),
	)
	return t
}

func runnerTableRows(runners []Runner) []table.Row {
	rows := make([]table.Row, 0, len(runners))
	for i, runner := range runners {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i+1),
			runner.Host,
			projectName(runner.Repo),
			runner.Name,
			runner.LocalState,
			runner.GitHubStatus,
			boolText(runner.Busy),
			queueText(runner),
			strings.Join(runner.Labels, ","),
			runner.Path,
		})
	}
	return rows
}

func runnerTableColumns(width int) []table.Column {
	available := max(width, 60)

	columns := []table.Column{
		{Title: "#", Width: 3},
		{Title: "Host", Width: 8},
		{Title: "Project", Width: 12},
		{Title: "Runner", Width: 16},
		{Title: "Local", Width: 8},
		{Title: "GitHub", Width: 8},
		{Title: "Busy", Width: 5},
		{Title: "Queue", Width: 8},
		{Title: "Labels", Width: 10},
		{Title: "Path", Width: 16},
	}

	hideColumnsToFit(&columns, available)
	base := tableRenderWidth(columns)
	if available >= base {
		extra := available - base
		addWidth(&columns, "Project", min(extra/5, 10))
		addWidth(&columns, "Runner", min(extra/3, 18))
		addWidth(&columns, "Labels", min(extra/6, 14))
		base = tableRenderWidth(columns)
		addWidth(&columns, "Path", max(0, available-base))
		return columns
	}

	over := tableRenderWidth(columns) - available
	shrinkWidth(&columns, "Path", &over, 8)
	shrinkWidth(&columns, "Labels", &over, 6)
	shrinkWidth(&columns, "Runner", &over, 10)
	shrinkWidth(&columns, "Project", &over, 8)
	shrinkWidth(&columns, "Host", &over, 5)
	over = tableRenderWidth(columns) - available
	shrinkWidth(&columns, "Queue", &over, 5)
	shrinkWidth(&columns, "GitHub", &over, 6)
	shrinkWidth(&columns, "Local", &over, 6)
	shrinkWidth(&columns, "Runner", &over, 8)
	shrinkWidth(&columns, "Project", &over, 6)
	hideColumnsToFit(&columns, available)
	return columns
}

func hideColumnsToFit(columns *[]table.Column, available int) {
	for _, title := range []string{"Path", "Labels", "Host"} {
		if tableRenderWidth(*columns) <= available {
			return
		}
		setColumnWidth(columns, title, 0)
	}
}

func tableHeight(windowHeight int) int {
	return max(1, windowHeight-8)
}

func commandHelp(width int) string {
	long := "Commands: h/? help | refresh | start [N] | stop [N] | restart [N] | force-stop [N] | clear [N] | remove [N] [confirm] | delete [N] confirm | clear idle | auto-clear on/off | logs [N] | connect remote NAME | q"
	short := "Commands: h/? help | refresh | start/stop/restart/clear/logs [N] | remove/delete [N] confirm | connect remote NAME | q"
	tiny := "Commands: h help | refresh | start/stop/clear/logs [N] | q"
	if width < 100 {
		return tiny
	}
	if width < 130 {
		return short
	}
	return long
}

func tuiHelp(width int, height int) string {
	lines := []string{
		"Help",
		"Select a runner with ↑/↓ or type a runner number as N. Commands without N use the selected row.",
		"",
		"h, ?, help        show or hide this help",
		"refresh           refresh local and GitHub runner state",
		"start/stop N      start or stop a controllable runner",
		"restart N         restart a controllable runner",
		"force-stop N      stop a busy runner when you intentionally override the busy check",
		"force-restart N   restart a busy runner when you intentionally override the busy check",
		"clear N           clear an idle runner _work directory and installer archives",
		"clear idle        clear all idle runners",
		"auto-clear on/off clear idle runners after each refresh, or disable it",
		"logs N            open the runner logs folder",
		"remove N confirm  unregister a runner; dry-run without confirm",
		"delete N confirm  unregister a runner and delete its folder when it is under a safe root",
		"connect remote NAME open a saved remote host TUI over SSH",
		"q, esc, ctrl+c    quit",
	}
	if width < 100 {
		lines = []string{
			"Help",
			"Use ↑/↓ to select; commands without N use the selected row.",
			"",
			"h/? help         toggle help",
			"refresh          refresh runner state",
			"start/stop N     control runner",
			"restart N        restart runner",
			"clear N          clear idle runner work files",
			"logs N           open runner logs",
			"remove N confirm unregister runner",
			"delete N confirm unregister and delete safe runner folder",
			"q/esc/ctrl+c     quit",
		}
	}
	limit := max(1, height-4)
	if len(lines) > limit {
		lines = lines[:limit]
	}
	for i := range lines {
		lines[i] = trunc(lines[i], max(20, width))
	}
	return strings.Join(lines, "\n")
}

func boolText(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func queueText(runner Runner) string {
	if runner.StaleQueueCount > 0 {
		return fmt.Sprintf("%d/%d stale", runner.QueueCount, runner.StaleQueueCount)
	}
	return fmt.Sprintf("%d", runner.QueueCount)
}

func columnWidthSum(columns []table.Column) int {
	total := 0
	for _, column := range columns {
		total += column.Width
	}
	return total
}

func tableRenderWidth(columns []table.Column) int {
	total := 0
	for _, column := range columns {
		if column.Width > 0 {
			total += column.Width + 2
		}
	}
	return total
}

func addWidth(columns *[]table.Column, title string, amount int) {
	if amount <= 0 {
		return
	}
	for i := range *columns {
		if (*columns)[i].Title == title && (*columns)[i].Width > 0 {
			(*columns)[i].Width += amount
			return
		}
	}
}

func setColumnWidth(columns *[]table.Column, title string, width int) {
	for i := range *columns {
		if (*columns)[i].Title == title {
			(*columns)[i].Width = width
			return
		}
	}
}

func shrinkWidth(columns *[]table.Column, title string, over *int, minimum int) {
	if *over <= 0 {
		return
	}
	for i := range *columns {
		if (*columns)[i].Title != title {
			continue
		}
		reducible := max(0, (*columns)[i].Width-minimum)
		reduceBy := min(reducible, *over)
		(*columns)[i].Width -= reduceBy
		*over -= reduceBy
		return
	}
}
