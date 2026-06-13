package app

import (
	"fmt"
	"io"
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
	updateNotice  string
	loading       bool
	refreshing    bool
	refreshEvery  time.Duration
	spinnerFrame  int
	autoClearIdle bool
	showHelp      bool
	width         int
	height        int
}

type refreshSource string

const (
	refreshStartup refreshSource = "startup"
	refreshManual  refreshSource = "manual"
	refreshAuto    refreshSource = "auto"
	refreshCommand refreshSource = "command"
)

type refreshResultMsg struct {
	inventory   Inventory
	err         error
	source      refreshSource
	completedAt time.Time
}

type spinnerTickMsg time.Time
type autoRefreshTickMsg time.Time

type remoteConnectDoneMsg struct {
	name string
	err  error
}

type commandResultMsg struct {
	message      string
	refreshAfter bool
}

type updateCheckDoneMsg struct {
	notice string
}

var hourglassFrames = []string{"⌛", "⏳"}

const terminalTextModeReset = "\x1b[0m\x1b(B"

func NewModel(inventory Inventory) Model {
	input := textinput.New()
	input.Placeholder = "refresh | start 1 | connect remote runnerbox | q"
	input.Focus()
	input.CharLimit = 120
	input.Width = 80

	model := Model{
		input:        input,
		inventory:    inventory,
		message:      "ready",
		refreshEvery: effectiveSettings().TUIRefreshInterval(),
		width:        120,
		height:       30,
	}
	model.table = newRunnerTable(inventory.Runners, model.width, tableHeight(model.height))
	model.resize(model.width, model.height)
	return model
}

func NewLoadingModel() Model {
	model := NewModel(Inventory{})
	model.loading = true
	model.refreshing = true
	model.message = "Ожидайте, идет опрос раннеров..."
	return model
}

func (m Model) Init() tea.Cmd {
	if m.loading {
		return tea.Batch(textinput.Blink, refreshInventoryCmd(refreshStartup), updateCheckCmd(), spinnerTickCmd(), m.autoRefreshTickCmd())
	}
	return tea.Batch(textinput.Blink, updateCheckCmd(), m.autoRefreshTickCmd())
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
		m.refreshing = false
		if msg.err != nil {
			if msg.source == refreshAuto {
				m.message = fmt.Sprintf("auto-refresh warning: %v", msg.err)
			} else {
				m.message = fmt.Sprintf("refresh completed with warnings: %v", msg.err)
			}
		} else if msg.source == refreshAuto {
			m.message = fmt.Sprintf("auto-refreshed at %s", msg.completedAt.Format("15:04:05"))
		} else if msg.source == refreshCommand {
			if m.message == "" {
				m.message = "status refreshed after command"
			}
		} else {
			m.message = "ready"
		}
		if m.autoClearIdle {
			m.message = "auto-clear idle runners..."
			return m, terminalManagedCommand(func() string {
				return strings.TrimSpace(ClearIdleRunners(m.inventory))
			}, false)
		}
		return m, nil
	case spinnerTickMsg:
		if !m.loading {
			return m, nil
		}
		m.spinnerFrame = (m.spinnerFrame + 1) % len(hourglassFrames)
		return m, spinnerTickCmd()
	case autoRefreshTickMsg:
		if m.refreshing {
			return m, m.autoRefreshTickCmd()
		}
		m.refreshing = true
		return m, tea.Batch(refreshInventoryCmd(refreshAuto), m.autoRefreshTickCmd())
	case remoteConnectDoneMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("remote %s connection failed: %v", msg.name, msg.err)
		} else {
			m.message = fmt.Sprintf("remote %s connection closed", msg.name)
		}
		return m, nil
	case commandResultMsg:
		m.message = msg.message
		if msg.refreshAfter && !m.refreshing {
			m.refreshing = true
			return m, refreshInventoryCmd(refreshCommand)
		}
		return m, nil
	case updateCheckDoneMsg:
		m.updateNotice = msg.notice
		m.syncTable()
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}
			return m, tea.Quit
		case tea.KeyEnter:
			if m.loading {
				return m, nil
			}
			command := strings.TrimSpace(m.input.Value())
			m.input.SetValue("")
			return m.runCommand(command)
		}
		if m.input.Value() == "" && (msg.String() == "h" || msg.String() == "?") {
			m.showHelp = !m.showHelp
			return m, nil
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

func (m Model) View() (view string) {
	var b strings.Builder
	b.WriteString(terminalTextModeReset)
	defer func() {
		view = b.String() + terminalTextModeReset
	}()

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")).Render("RunnerMonitor " + CurrentVersion)
	updateNotice := renderUpdateNotice(m.updateNotice, max(20, m.width))
	b.WriteString(title)
	b.WriteString("\n")
	if m.loading {
		frame := hourglassFrames[m.spinnerFrame%len(hourglassFrames)]
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(frame + " Ожидайте, идет опрос раннеров..."))
		b.WriteString("\n")
		return b.String()
	}
	if m.height <= 8 {
		if m.updateNotice != "" {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(updateNotice))
			b.WriteString("\n")
		}
		if m.message != "" {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(trunc(m.message, max(20, m.width))))
			b.WriteString("\n")
		}
		b.WriteString(m.input.View())
		b.WriteString("\n")
		return b.String()
	}
	if m.showHelp {
		b.WriteString(tuiHelp(m.width, m.height))
		b.WriteString("\n")
		if m.message != "" {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(trunc(m.message, max(20, m.width))))
			b.WriteString("\n")
		}
		if m.updateNotice != "" {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(updateNotice))
			b.WriteString("\n")
		}
		b.WriteString(m.input.View())
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
	if m.updateNotice != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(updateNotice))
		b.WriteString("\n")
	}
	b.WriteString(m.input.View())
	b.WriteString("\n")
	return b.String()
}

func renderUpdateNotice(notice string, width int) string {
	open := strings.LastIndex(notice, "(")
	close := strings.LastIndex(notice, ")")
	if open < 0 || close <= open {
		return trunc(notice, width)
	}
	url := strings.TrimSpace(notice[open+1 : close])
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		return trunc(notice, width)
	}
	prefix := notice[:open+1]
	suffix := notice[close:]
	labelWidth := max(1, width-len(prefix)-len(suffix))
	return prefix + terminalHyperlink(url, trunc(url, labelWidth)) + suffix
}

func terminalHyperlink(url string, label string) string {
	return "\x1b]8;;" + url + "\x1b\\" + label + "\x1b]8;;\x1b\\"
}

func (m *Model) resize(width, height int) {
	m.width = max(width, 60)
	m.height = max(height, 1)
	m.input.Width = max(20, m.width-2)
	m.syncTable()
}

func (m *Model) syncTable() {
	height := tableHeight(m.height)
	if m.updateNotice != "" {
		height = max(1, height-1)
	}
	if len(m.table.Columns()) == 0 {
		m.table = newRunnerTable(m.inventory.Runners, m.width, height)
		return
	}
	m.table.SetColumns(runnerTableColumns(m.width))
	m.table.SetRows(runnerTableRows(m.inventory.Runners))
	m.table.SetWidth(m.width)
	m.table.SetHeight(height)
}

func (m Model) runCommand(command string) (tea.Model, tea.Cmd) {
	if command == "" {
		return m, nil
	}
	if command == "q" || command == "quit" || command == "exit" {
		return m, tea.Quit
	}
	if command == "h" || command == "help" || command == "?" {
		m.showHelp = !m.showHelp
		return m, nil
	}
	if command == "refresh" {
		if m.refreshing {
			m.message = "refresh already in progress"
			return m, nil
		}
		m.loading = len(m.inventory.Runners) == 0
		m.refreshing = true
		if m.loading {
			m.message = "Ожидайте, идет опрос раннеров..."
			return m, tea.Batch(refreshInventoryCmd(refreshManual), spinnerTickCmd())
		}
		m.message = "refreshing..."
		return m, refreshInventoryCmd(refreshManual)
	}
	if command == "clear idle" {
		m.message = "clearing idle runners..."
		return m, terminalManagedCommand(func() string {
			return strings.TrimSpace(ClearIdleRunners(m.inventory))
		}, false)
	}
	if command == "auto-clear on" {
		m.autoClearIdle = true
		m.message = "auto-clear enabled; clearing idle runners..."
		return m, terminalManagedCommand(func() string {
			return strings.TrimSpace(ClearIdleRunners(m.inventory))
		}, false)
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
		if len(parts) == 2 && parts[1] == "all" && (parts[0] == "start" || parts[0] == "stop") {
			action := parts[0]
			m.message = fmt.Sprintf("%s all runners...", action)
			return m, terminalManagedCommand(func() string {
				return strings.TrimSpace(RunAllLifecycle(action, m.inventory))
			}, true)
		}
		runner, ok := m.commandRunner(parts)
		if !ok {
			return m, nil
		}
		action := parts[0]
		m.message = fmt.Sprintf("%s %s...", action, runner.Name)
		return m, terminalManagedCommand(func() string {
			return RunLifecycle(action, runner)
		}, true)
	case "clear":
		runner, ok := m.commandRunner(parts)
		if !ok {
			return m, nil
		}
		m.message = fmt.Sprintf("clearing %s...", runner.Name)
		return m, terminalManagedCommand(func() string {
			return ClearRunner(runner)
		}, false)
	case "remove":
		runner, ok := m.commandRunner(parts)
		if !ok {
			return m, nil
		}
		confirm := commandHasConfirm(parts)
		m.message = fmt.Sprintf("removing %s...", runner.Name)
		return m, terminalManagedCommand(func() string {
			return RemoveRunner(runner, RemoveRunnerOptions{Confirm: confirm})
		}, false)
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
		return m, terminalManagedCommand(func() string {
			return RemoveRunner(runner, RemoveRunnerOptions{Confirm: true, DeleteFolder: true})
		}, false)
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
	case "up", "down", "k", "j", "pgup", "pgdown", "home", "end", "g", "G", " ", "b", "f", "u", "ctrl+u", "ctrl+d":
		return true
	default:
		return false
	}
}

type terminalManagedAction struct {
	run     func() string
	message string
}

func (action *terminalManagedAction) Run() error {
	action.message = action.run()
	return nil
}

func (action *terminalManagedAction) SetStdin(io.Reader)  {}
func (action *terminalManagedAction) SetStdout(io.Writer) {}
func (action *terminalManagedAction) SetStderr(io.Writer) {}

func terminalManagedCommand(run func() string, refreshAfter bool) tea.Cmd {
	action := &terminalManagedAction{run: run}
	return tea.Exec(action, func(err error) tea.Msg {
		message := strings.TrimSpace(action.message)
		if err != nil {
			if message != "" {
				message += "; "
			}
			message += fmt.Sprintf("terminal restore failed: %v", err)
		}
		return commandResultMsg{message: message, refreshAfter: refreshAfter}
	})
}

func refreshInventoryCmd(source refreshSource) tea.Cmd {
	return func() tea.Msg {
		var inventory Inventory
		var err error
		if source == refreshAuto {
			inventory, err = RefreshWithGitHubCache(autoRefreshGitHubCacheTTL)
		} else {
			inventory, err = Refresh()
		}
		return refreshResultMsg{inventory: inventory, err: err, source: source, completedAt: time.Now()}
	}
}

func updateCheckCmd() tea.Cmd {
	return func() tea.Msg {
		return updateCheckDoneMsg{notice: CheckForUpdate(CurrentVersion)}
	}
}

func spinnerTickCmd() tea.Cmd {
	return tea.Tick(140*time.Millisecond, func(t time.Time) tea.Msg {
		return spinnerTickMsg(t)
	})
}

func (m Model) autoRefreshTickCmd() tea.Cmd {
	interval := m.refreshEvery
	if interval <= 0 {
		interval = DefaultSettings().TUIRefreshInterval()
	}
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return autoRefreshTickMsg(t)
	})
}
