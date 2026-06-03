package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	input     textinput.Model
	inventory Inventory
	message   string
}

func NewModel(inventory Inventory) Model {
	input := textinput.New()
	input.Placeholder = "refresh | start 1 | stop 1 | restart 1 | logs 1 | q"
	input.Focus()
	input.CharLimit = 120
	input.Width = 80

	return Model{
		input:     input,
		inventory: inventory,
		message:   "ready",
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			command := strings.TrimSpace(m.input.Value())
			m.input.SetValue("")
			return m.runCommand(command)
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var b strings.Builder
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")).Render("RunnerMonitor")
	b.WriteString(title)
	b.WriteString("\n")
	b.WriteString("Commands: refresh | start N | stop N | restart N | force-stop N | logs N | q\n\n")
	b.WriteString(renderTable(m.inventory.Runners))
	b.WriteString("\n")
	if m.message != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(m.message))
		b.WriteString("\n")
	}
	b.WriteString("> ")
	b.WriteString(m.input.View())
	b.WriteString("\n")
	return b.String()
}

func (m Model) runCommand(command string) (tea.Model, tea.Cmd) {
	if command == "" {
		return m, nil
	}
	if command == "q" || command == "quit" || command == "exit" {
		return m, tea.Quit
	}
	if command == "refresh" {
		inventory, err := Refresh()
		m.inventory = inventory
		if err != nil {
			m.message = fmt.Sprintf("refresh completed with warnings: %v", err)
		} else {
			m.message = "refreshed"
		}
		return m, nil
	}

	parts := strings.Fields(command)
	if len(parts) != 2 {
		m.message = "unknown command"
		return m, nil
	}

	index, err := strconv.Atoi(parts[1])
	if err != nil || index < 1 || index > len(m.inventory.Runners) {
		m.message = "invalid runner number"
		return m, nil
	}

	runner := m.inventory.Runners[index-1]
	switch parts[0] {
	case "start", "stop", "restart", "force-stop", "force-restart":
		m.message = RunLifecycle(parts[0], runner)
	case "logs":
		m.message = OpenLogs(runner)
	default:
		m.message = "unknown command"
	}
	return m, nil
}
