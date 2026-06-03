package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	input        textinput.Model
	inventory    Inventory
	message      string
	loading      bool
	spinnerFrame int
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

var hourglassFrames = []string{"⌛", "⏳"}

func NewModel(inventory Inventory) Model {
	input := textinput.New()
	input.Placeholder = "refresh | start 1 | connect remote runnerbox | q"
	input.Focus()
	input.CharLimit = 120
	input.Width = 80

	return Model{
		input:     input,
		inventory: inventory,
		message:   "ready",
	}
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
	case refreshResultMsg:
		m.inventory = msg.inventory
		m.loading = false
		if msg.err != nil {
			m.message = fmt.Sprintf("refresh completed with warnings: %v", msg.err)
		} else {
			m.message = "ready"
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
	b.WriteString("Commands: refresh | start N | stop N | restart N | force-stop N | logs N | connect remote NAME | q\n\n")
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
		m.loading = true
		m.message = "Ожидайте, идет опрос раннеров..."
		return m, tea.Batch(refreshInventoryCmd(), spinnerTickCmd())
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
