package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	input         textinput.Model
	table         table.Model
	inventory     Inventory
	message       string
	loading       bool
	spinnerFrame  int
	autoClearIdle bool
	width         int
	height        int
}

type refreshResultMsg struct {
	inventory Inventory
	err       error
}

type spinnerTickMsg time.Time

type remoteConnectDoneMsg struct {
	name string
	err  error
}

type clearResultMsg struct {
	message string
}

var hourglassFrames = []string{"⌛", "⏳"}

func NewModel(inventory Inventory) Model {
	input := textinput.New()
	input.Placeholder = "refresh | start 1 | connect remote runnerbox | q"
	input.Focus()
	input.CharLimit = 120
	input.Width = 80

	model := Model{
		input:     input,
		inventory: inventory,
		message:   "ready",
		width:     120,
		height:    30,
	}
	model.table = newRunnerTable(inventory.Runners, model.width, tableHeight(model.height))
	model.resize(model.width, model.height)
	return model
}

func NewLoadingModel() Model {
	model := NewModel(Inventory{})
	model.loading = true
	model.message = "Ожидайте, идет опрос раннеров..."
	return model
}

func (m Model) Init() tea.Cmd {
	if m.loading {
		return tea.Batch(textinput.Blink, refreshInventoryCmd(), spinnerTickCmd())
	}
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil
	case refreshResultMsg:
		m.inventory = msg.inventory
		m.syncTable()
		m.loading = false
		if msg.err != nil {
			m.message = fmt.Sprintf("refresh completed with warnings: %v", msg.err)
		} else {
			m.message = "ready"
		}
		if m.autoClearIdle {
			m.message = "auto-clear idle runners..."
			return m, clearIdleRunnersCmd(m.inventory)
		}
		return m, nil
	case spinnerTickMsg:
		if !m.loading {
			return m, nil
		}
		m.spinnerFrame = (m.spinnerFrame + 1) % len(hourglassFrames)
		return m, spinnerTickCmd()
	case remoteConnectDoneMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("remote %s connection failed: %v", msg.name, msg.err)
		} else {
			m.message = fmt.Sprintf("remote %s connection closed", msg.name)
		}
		return m, nil
	case clearResultMsg:
		m.message = msg.message
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.loading {
				return m, nil
			}
			command := strings.TrimSpace(m.input.Value())
			m.input.SetValue("")
			return m.runCommand(command)
		}
		if m.input.Value() == "" && isTableNavigationKey(msg) {
			var cmd tea.Cmd
			m.table, cmd = m.table.Update(msg)
			return m, cmd
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
	if m.loading {
		frame := hourglassFrames[m.spinnerFrame%len(hourglassFrames)]
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(frame + " Ожидайте, идет опрос раннеров..."))
		b.WriteString("\n")
		return b.String()
	}
	b.WriteString(commandHelp(m.width))
	b.WriteString("\n\n")
	b.WriteString(m.table.View())
	b.WriteString("\n")
	if m.message != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(trunc(m.message, max(20, m.width))))
		b.WriteString("\n")
	}
	b.WriteString(m.input.View())
	b.WriteString("\n")
	return b.String()
}

func (m *Model) resize(width, height int) {
	m.width = max(width, 60)
	m.height = max(height, 12)
	m.input.Width = max(20, m.width-2)
	m.syncTable()
}

func (m *Model) syncTable() {
	if len(m.table.Columns()) == 0 {
		m.table = newRunnerTable(m.inventory.Runners, m.width, tableHeight(m.height))
		return
	}
	m.table.SetColumns(runnerTableColumns(m.width))
	m.table.SetRows(runnerTableRows(m.inventory.Runners))
	m.table.SetWidth(m.width)
	m.table.SetHeight(tableHeight(m.height))
}

func (m Model) runCommand(command string) (tea.Model, tea.Cmd) {
	if command == "" {
		return m, nil
	}
	if command == "q" || command == "quit" || command == "exit" {
		return m, tea.Quit
	}
	if command == "refresh" {
		m.loading = true
		m.message = "Ожидайте, идет опрос раннеров..."
		return m, tea.Batch(refreshInventoryCmd(), spinnerTickCmd())
	}
	if command == "clear idle" {
		m.message = "clearing idle runners..."
		return m, clearIdleRunnersCmd(m.inventory)
	}
	if command == "auto-clear on" {
		m.autoClearIdle = true
		m.message = "auto-clear enabled; clearing idle runners..."
		return m, clearIdleRunnersCmd(m.inventory)
	}
	if command == "auto-clear off" {
		m.autoClearIdle = false
		m.message = "auto-clear disabled"
		return m, nil
	}
	if command == "connect remote" || strings.HasPrefix(command, "connect remote ") {
		name := strings.TrimSpace(strings.TrimPrefix(command, "connect remote"))
		if name == "" {
			name = "runnerbox"
		}
		cmd, err := RemoteTUIProcess(name)
		if err != nil {
			m.message = err.Error()
			return m, nil
		}
		m.message = fmt.Sprintf("connecting remote %s", name)
		return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
			return remoteConnectDoneMsg{name: name, err: err}
		})
	}

	parts := strings.Fields(command)
	if len(parts) < 1 || len(parts) > 3 {
		m.message = "unknown command"
		return m, nil
	}

	switch parts[0] {
	case "start", "stop", "restart", "force-stop", "force-restart":
		runner, ok := m.commandRunner(parts)
		if !ok {
			return m, nil
		}
		m.message = RunLifecycle(parts[0], runner)
	case "clear":
		runner, ok := m.commandRunner(parts)
		if !ok {
			return m, nil
		}
		m.message = fmt.Sprintf("clearing %s...", runner.Name)
		return m, clearRunnerCmd(runner)
	case "remove":
		runner, ok := m.commandRunner(parts)
		if !ok {
			return m, nil
		}
		confirm := commandHasConfirm(parts)
		m.message = fmt.Sprintf("removing %s...", runner.Name)
		return m, removeRunnerCmd(runner, RemoveRunnerOptions{Confirm: confirm})
	case "delete":
		runner, ok := m.commandRunner(parts)
		if !ok {
			return m, nil
		}
		if !commandHasConfirm(parts) {
			m.message = "delete requires: delete [N] confirm"
			return m, nil
		}
		m.message = fmt.Sprintf("removing %s and deleting folder...", runner.Name)
		return m, removeRunnerCmd(runner, RemoveRunnerOptions{Confirm: true, DeleteFolder: true})
	case "logs":
		runner, ok := m.commandRunner(parts)
		if !ok {
			return m, nil
		}
		m.message = OpenLogs(runner)
	default:
		m.message = "unknown command"
	}
	return m, nil
}

func (m *Model) commandRunner(parts []string) (Runner, bool) {
	index, err := commandRunnerIndex(parts, m.table.Cursor(), len(m.inventory.Runners))
	if err != nil {
		m.message = err.Error()
		return Runner{}, false
	}
	return m.inventory.Runners[index], true
}

func commandRunnerIndex(parts []string, selected int, runnerCount int) (int, error) {
	if runnerCount == 0 {
		return 0, fmt.Errorf("no runners available")
	}
	if len(parts) >= 2 && parts[1] != "confirm" {
		index, err := strconv.Atoi(parts[1])
		if err != nil || index < 1 || index > runnerCount {
			return 0, fmt.Errorf("invalid runner number")
		}
		return index - 1, nil
	}
	if selected >= 0 && selected < runnerCount {
		return selected, nil
	}
	return 0, fmt.Errorf("invalid runner number")
}

func commandHasConfirm(parts []string) bool {
	return len(parts) >= 2 && parts[len(parts)-1] == "confirm"
}

func isTableNavigationKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "up", "down", "k", "j", "pgup", "pgdown", "home", "end", "g", "G", " ", "b", "f", "u", "d", "ctrl+u", "ctrl+d":
		return true
	default:
		return false
	}
}

func clearRunnerCmd(runner Runner) tea.Cmd {
	return func() tea.Msg {
		return clearResultMsg{message: ClearRunner(runner)}
	}
}

func clearIdleRunnersCmd(inventory Inventory) tea.Cmd {
	return func() tea.Msg {
		return clearResultMsg{message: strings.TrimSpace(ClearIdleRunners(inventory))}
	}
}

func removeRunnerCmd(runner Runner, options RemoveRunnerOptions) tea.Cmd {
	return func() tea.Msg {
		return clearResultMsg{message: RemoveRunner(runner, options)}
	}
}

func refreshInventoryCmd() tea.Cmd {
	return func() tea.Msg {
		inventory, err := Refresh()
		return refreshResultMsg{inventory: inventory, err: err}
	}
}

func spinnerTickCmd() tea.Cmd {
	return tea.Tick(140*time.Millisecond, func(t time.Time) tea.Msg {
		return spinnerTickMsg(t)
	})
}
