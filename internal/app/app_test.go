package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func TestParseRunnerConfigWithBOM(t *testing.T) {
	config, err := parseRunnerConfig([]byte("\uFEFF{\"agentName\":\"runner-1\",\"gitHubUrl\":\"https://github.com/SGribanov/RunnerMonitor\",\"workFolder\":\"_work\"}"))
	if err != nil {
		t.Fatalf("parseRunnerConfig returned error: %v", err)
	}
	if config.AgentName != "runner-1" {
		t.Fatalf("AgentName = %q", config.AgentName)
	}
}

func TestRepoFromGitHubURL(t *testing.T) {
	cases := map[string]string{
		"https://github.com/SGribanov/RunnerMonitor":     "SGribanov/RunnerMonitor",
		"git@github.com:SGribanov/RunnerMonitor.git":     "SGribanov/RunnerMonitor",
		"https://github.com/SGribanov/RunnerMonitor.git": "SGribanov/RunnerMonitor",
		"git@github.com:SGribanov/RunnerMonitor":         "SGribanov/RunnerMonitor",
	}
	for input, want := range cases {
		if got := repoFromGitHubURL(input); got != want {
			t.Fatalf("repoFromGitHubURL(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestUniqueRepos(t *testing.T) {
	repos := uniqueRepos([]Runner{{Repo: "b/repo"}, {Repo: "a/repo"}, {Repo: "b/repo"}, {}})
	if len(repos) != 2 || repos[0] != "a/repo" || repos[1] != "b/repo" {
		t.Fatalf("unexpected repos: %#v", repos)
	}
}

func TestLoadingModelShowsWaitMessageBeforeTable(t *testing.T) {
	view := NewLoadingModel().View()
	if !strings.Contains(view, "Ожидайте, идет опрос раннеров...") {
		t.Fatalf("loading view does not contain wait message: %q", view)
	}
	if strings.Contains(view, "Commands:") || strings.Contains(view, "No runners discovered") {
		t.Fatalf("loading view should not show table or commands: %q", view)
	}
}

func TestModelTitleShowsCurrentVersion(t *testing.T) {
	view := NewModel(Inventory{}).View()
	if !strings.Contains(view, "RunnerMonitor "+CurrentVersion) {
		t.Fatalf("view title should include current version %q: %q", CurrentVersion, view)
	}
}

func TestAutoRefreshTickStartsNonOverlappingRefresh(t *testing.T) {
	model := NewModel(Inventory{Runners: []Runner{{Name: "runner-1", Repo: "SGribanov/RunnerMonitor"}}})

	updated, cmd := model.Update(autoRefreshTickMsg{})
	refreshed, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T", updated)
	}
	if !refreshed.refreshing {
		t.Fatalf("auto-refresh should mark refresh in progress")
	}
	if cmd == nil {
		t.Fatalf("auto-refresh should return refresh and next tick commands")
	}

	updated, cmd = refreshed.Update(autoRefreshTickMsg{})
	skipped, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T", updated)
	}
	if !skipped.refreshing {
		t.Fatalf("overlapping auto-refresh should keep current refresh in progress")
	}
	if cmd == nil {
		t.Fatalf("overlapping auto-refresh should still schedule the next tick")
	}
}

func TestManualRefreshKeepsExistingInventoryVisible(t *testing.T) {
	model := NewModel(Inventory{Runners: []Runner{{Name: "runner-1", Repo: "SGribanov/RunnerMonitor"}}})

	updated, cmd := model.runCommand("refresh")
	refreshed, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T", updated)
	}
	if refreshed.loading {
		t.Fatalf("manual refresh with existing inventory should not enter loading-only view")
	}
	if !refreshed.refreshing {
		t.Fatalf("manual refresh should mark refresh in progress")
	}
	if !strings.Contains(refreshed.View(), "runner-1") {
		t.Fatalf("manual refresh should keep existing table visible: %q", refreshed.View())
	}
	if cmd == nil {
		t.Fatalf("manual refresh should return refresh command")
	}
}

func TestManualRefreshSkipsWhenRefreshAlreadyRunning(t *testing.T) {
	model := NewModel(Inventory{Runners: []Runner{{Name: "runner-1", Repo: "SGribanov/RunnerMonitor"}}})
	model.refreshing = true

	updated, cmd := model.runCommand("refresh")
	refreshed, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T", updated)
	}
	if cmd != nil {
		t.Fatalf("manual refresh while refreshing should not start another command")
	}
	if refreshed.message != "refresh already in progress" {
		t.Fatalf("message = %q", refreshed.message)
	}
}

func TestRunnerTableColumnsFitNarrowWidth(t *testing.T) {
	columns := runnerTableColumns(60)
	if len(columns) == 0 {
		t.Fatalf("runnerTableColumns returned no columns")
	}
	for _, column := range columns {
		if column.Width < 0 {
			t.Fatalf("column %q has invalid width %d", column.Title, column.Width)
		}
	}
	if columnWidth(columns, "Project") <= 0 || columnWidth(columns, "Runner") <= 0 {
		t.Fatalf("project and runner columns must remain visible: %#v", columns)
	}
	if got := tableRenderWidth(columns); got > 60 {
		t.Fatalf("columns exceed narrow width: %d", got)
	}
}

func TestRunnerTableColumnsUseExtraWidth(t *testing.T) {
	narrow := runnerTableColumns(80)
	wide := runnerTableColumns(150)
	if columnWidth(wide, "Path") <= columnWidth(narrow, "Path") {
		t.Fatalf("wide path column should grow: narrow=%d wide=%d", columnWidth(narrow, "Path"), columnWidth(wide, "Path"))
	}
	if columnWidth(wide, "Runner") <= columnWidth(narrow, "Runner") {
		t.Fatalf("wide runner column should grow: narrow=%d wide=%d", columnWidth(narrow, "Runner"), columnWidth(wide, "Runner"))
	}
}

func TestModelWindowSizeMessageResizesTableAndInput(t *testing.T) {
	model := NewModel(Inventory{Runners: []Runner{{Name: "runner-1", Repo: "SGribanov/RunnerMonitor"}}})
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 72, Height: 18})
	resized, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T", updated)
	}
	if resized.width != 72 || resized.height != 18 {
		t.Fatalf("size = %dx%d", resized.width, resized.height)
	}
	if resized.input.Width != 70 {
		t.Fatalf("input width = %d", resized.input.Width)
	}
	if resized.table.Width() != 72 || resized.table.Height() != tableHeight(18)-1 {
		t.Fatalf("table size = %dx%d", resized.table.Width(), resized.table.Height())
	}
}

func TestModelKeepsSmallWindowHeightCompact(t *testing.T) {
	model := NewModel(Inventory{Runners: []Runner{{Name: "runner-1", Repo: "SGribanov/RunnerMonitor"}}})
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 50, Height: 5})
	resized, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T", updated)
	}
	if resized.height != 5 {
		t.Fatalf("height = %d", resized.height)
	}
	view := resized.View()
	if strings.Contains(view, "Commands:") || strings.Contains(view, "runner-1") {
		t.Fatalf("compact view should not render full table/help: %q", view)
	}
	if !strings.Contains(view, "ready") || !strings.Contains(view, "RunnerMonitor") {
		t.Fatalf("compact view should keep status and title visible: %q", view)
	}
}

func TestRunnerTableRowsIncludeProjectAndQueue(t *testing.T) {
	rows := runnerTableRows([]Runner{{
		Name:            "runner-1",
		Repo:            "SGribanov/RunnerMonitor",
		Host:            "local",
		LocalState:      "running",
		GitHubStatus:    "online",
		Busy:            true,
		QueueCount:      2,
		StaleQueueCount: 1,
		Labels:          []string{"self-hosted", "Windows"},
		Path:            `C:\Runners\SGribanov-RunnerMonitor\runner-1`,
	}})
	if len(rows) != 1 {
		t.Fatalf("rows = %d", len(rows))
	}
	if rows[0][2] != "RunnerMonitor" {
		t.Fatalf("project column = %q", rows[0][2])
	}
	if !strings.Contains(rows[0][6], "true") || rows[0][7] != "2/1 stale" {
		t.Fatalf("busy/queue columns = %q/%q", rows[0][6], rows[0][7])
	}
}

func TestBusyTextIsPlainTableCellText(t *testing.T) {
	busy := busyText(true)
	if busy != "true" {
		t.Fatalf("busy true should stay plain table text, got %q", busy)
	}
	if idle := busyText(false); idle != "false" {
		t.Fatalf("busy false should stay plain, got %q", idle)
	}
}

func TestCommandHelpUsesCompactTextForNarrowWidth(t *testing.T) {
	got := commandHelp(70)
	if strings.Contains(got, "connect remote") {
		t.Fatalf("narrow help should be compact: %q", got)
	}
	if !strings.Contains(got, "h help") {
		t.Fatalf("narrow help should mention help command: %q", got)
	}
	if !strings.Contains(commandHelp(150), "connect remote NAME") {
		t.Fatalf("wide help should include remote command")
	}
}

func TestTUIHelpViewDescribesCommands(t *testing.T) {
	model := NewModel(Inventory{Runners: []Runner{{Name: "runner-1", Repo: "SGribanov/RunnerMonitor"}}})
	updated, _ := model.runCommand("help")
	helped, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T", updated)
	}
	view := helped.View()
	for _, want := range []string{"Help", "Select a runner", "refresh", "start/stop", "remove N confirm", "q, esc, ctrl+c"} {
		if !strings.Contains(view, want) {
			t.Fatalf("help view missing %q:\n%s", want, view)
		}
	}
	if strings.Contains(view, "runner-1") {
		t.Fatalf("help view should replace the table while open:\n%s", view)
	}
}

func TestTUIHelpKeyTogglesHelpAndEscClosesIt(t *testing.T) {
	model := NewModel(Inventory{Runners: []Runner{{Name: "runner-1", Repo: "SGribanov/RunnerMonitor"}}})
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	helped, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T", updated)
	}
	if !helped.showHelp {
		t.Fatalf("h should open help")
	}
	updated, cmd := helped.Update(tea.KeyMsg{Type: tea.KeyEsc})
	closed, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T", updated)
	}
	if closed.showHelp {
		t.Fatalf("esc should close help before quitting")
	}
	if cmd != nil {
		t.Fatalf("esc while help is open should not quit")
	}
}

func TestUpdateNoticeOnlyShowsNewerRelease(t *testing.T) {
	got := updateNotice("v0.2.0", "v0.3.0", "https://github.com/SGribanov/RunnerMonitor/releases/tag/v0.3.0")
	if !strings.Contains(got, "update available: v0.2.0 -> v0.3.0") {
		t.Fatalf("updateNotice did not report newer release: %q", got)
	}
	if updateNotice("v0.2.0", "v0.2.0", "") != "" {
		t.Fatalf("same version should not report an update")
	}
	if updateNotice("v0.2.0", "v0.1.9", "") != "" {
		t.Fatalf("older version should not report an update")
	}
}

func TestModelShowsUpdateNoticeWithoutReplacingStatus(t *testing.T) {
	model := NewModel(Inventory{Runners: []Runner{{Name: "runner-1", Repo: "SGribanov/RunnerMonitor"}}})
	updated, _ := model.Update(updateCheckDoneMsg{notice: "update available: v0.2.0 -> v0.3.0"})
	noticed, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T", updated)
	}
	view := noticed.View()
	if !strings.Contains(view, "update available: v0.2.0 -> v0.3.0") {
		t.Fatalf("view missing update notice:\n%s", view)
	}
	if !strings.Contains(view, "ready") {
		t.Fatalf("update notice should not replace status message:\n%s", view)
	}
}

func TestModelRendersUpdateNoticeURLAsTerminalLink(t *testing.T) {
	releaseURL := "https://github.com/SGribanov/RunnerMonitor/releases/tag/v0.3.0"
	model := NewModel(Inventory{Runners: []Runner{{Name: "runner-1", Repo: "SGribanov/RunnerMonitor"}}})
	updated, _ := model.Update(updateCheckDoneMsg{notice: "update available: v0.2.0 -> v0.3.0 (" + releaseURL + ")"})
	noticed, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T", updated)
	}
	view := noticed.View()
	if !strings.Contains(view, "\x1b]8;;"+releaseURL+"\x1b\\") {
		t.Fatalf("view missing clickable update URL:\n%s", view)
	}
	if !strings.Contains(view, "github.com/SGribanov/RunnerMonitor") {
		t.Fatalf("view should keep release URL readable:\n%s", view)
	}
}

func TestCommandRunnerIndexUsesSelectionWhenNumberMissing(t *testing.T) {
	got, err := commandRunnerIndex([]string{"start"}, 1, 3)
	if err != nil {
		t.Fatalf("commandRunnerIndex returned error: %v", err)
	}
	if got != 1 {
		t.Fatalf("selected index = %d", got)
	}
	numbered, err := commandRunnerIndex([]string{"start", "3"}, 1, 3)
	if err != nil {
		t.Fatalf("commandRunnerIndex numbered returned error: %v", err)
	}
	if numbered != 2 {
		t.Fatalf("numbered index = %d", numbered)
	}
}

func TestRunLifecycleRejectsManualRunner(t *testing.T) {
	got := RunLifecycle("start", Runner{Name: "manual-runner"})
	want := "manual-runner is not service-managed; cannot start"
	if got != want {
		t.Fatalf("RunLifecycle = %q, want %q", got, want)
	}
}

func columnWidth(columns []table.Column, title string) int {
	for _, column := range columns {
		if column.Title == title {
			return column.Width
		}
	}
	return 0
}

func TestRunLifecycleProtectsBusyRunner(t *testing.T) {
	got := RunLifecycle("stop", Runner{Name: "busy-runner", Busy: true, ServiceName: "svc", ControlMode: "unsupported"})
	want := "busy-runner is busy; use force-stop to override"
	if got != want {
		t.Fatalf("RunLifecycle = %q, want %q", got, want)
	}
}

func TestRunLifecycleManualWindowsAlreadyRunningDoesNotSpawn(t *testing.T) {
	got := RunLifecycle("start", Runner{
		Name:        "manual-windows",
		LocalState:  "running",
		ControlMode: "manual",
		Transport:   "windows",
		Path:        `C:\actions-runner-manual`,
	})
	want := "manual-windows already running"
	if got != want {
		t.Fatalf("RunLifecycle = %q, want %q", got, want)
	}
}

func TestRunLifecycleStartEnablesStartsAndWaitsForWSLRunnerOnline(t *testing.T) {
	restore := stubLifecycleControls(t)
	defer restore()

	var actions []string
	runServiceAction = func(controlMode, action, serviceName string) error {
		actions = append(actions, controlMode+" "+action+" "+serviceName)
		return nil
	}
	localStates := []string{"inactive", "active"}
	serviceState = func(controlMode, serviceName string) (string, error) {
		state := localStates[0]
		if len(localStates) > 1 {
			localStates = localStates[1:]
		}
		return state, nil
	}
	githubStatuses := []string{"offline", "online"}
	loadRunnerGitHubStatus = func(repo, name string) (string, error) {
		status := githubStatuses[0]
		if len(githubStatuses) > 1 {
			githubStatuses = githubStatuses[1:]
		}
		return status, nil
	}

	got := RunLifecycle("start", Runner{
		Name:        "runner-1",
		Repo:        "SGribanov/RunnerMonitor",
		ServiceName: "actions.runner.SGribanov-RunnerMonitor.runner-1.service",
		ControlMode: "wsl-systemd",
		LocalState:  "inactive",
	})

	if !strings.Contains(got, "service active; GitHub online") {
		t.Fatalf("RunLifecycle should confirm active and online, got %q", got)
	}
	wantActions := []string{
		"wsl-systemd enable actions.runner.SGribanov-RunnerMonitor.runner-1.service",
		"wsl-systemd start actions.runner.SGribanov-RunnerMonitor.runner-1.service",
	}
	if strings.Join(actions, "\n") != strings.Join(wantActions, "\n") {
		t.Fatalf("actions = %#v, want %#v", actions, wantActions)
	}
}

func TestRunLifecycleStartReportsWhenGitHubRunnerStaysOffline(t *testing.T) {
	restore := stubLifecycleControls(t)
	defer restore()

	loadRunnerGitHubStatus = func(repo, name string) (string, error) {
		return "offline", nil
	}

	got := RunLifecycle("start", Runner{
		Name:        "runner-1",
		Repo:        "SGribanov/RunnerMonitor",
		ServiceName: "actions.runner.SGribanov-RunnerMonitor.runner-1.service",
		ControlMode: "wsl-systemd",
		LocalState:  "inactive",
	})

	if !strings.Contains(got, "did not become online") || !strings.Contains(got, "last status: offline") {
		t.Fatalf("RunLifecycle should report offline GitHub verification, got %q", got)
	}
}

func stubLifecycleControls(t *testing.T) func() {
	t.Helper()
	oldRunServiceAction := runServiceAction
	oldServiceState := serviceState
	oldLoadRunnerGitHubStatus := loadRunnerGitHubStatus
	oldSleepForLifecyclePoll := sleepForLifecyclePoll

	runServiceAction = func(controlMode, action, serviceName string) error {
		return nil
	}
	serviceState = func(controlMode, serviceName string) (string, error) {
		return "active", nil
	}
	loadRunnerGitHubStatus = func(repo, name string) (string, error) {
		return "online", nil
	}
	sleepForLifecyclePoll = func(time.Duration) {}

	return func() {
		runServiceAction = oldRunServiceAction
		serviceState = oldServiceState
		loadRunnerGitHubStatus = oldLoadRunnerGitHubStatus
		sleepForLifecyclePoll = oldSleepForLifecyclePoll
	}
}

func TestClearRunnerRejectsBusyRunner(t *testing.T) {
	got := ClearRunner(Runner{Name: "busy", Busy: true, Path: t.TempDir()})
	if got != "busy is busy; cleanup skipped" {
		t.Fatalf("ClearRunner = %q", got)
	}
}

func TestClearRunnerRemovesWorkContentsAndInstallerArchives(t *testing.T) {
	root := t.TempDir()
	work := filepath.Join(root, "_work")
	if err := os.MkdirAll(filepath.Join(work, "repo"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(work, "repo", "file.txt"), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "actions-runner.zip"), []byte("zip"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".runner"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := ClearRunner(Runner{Name: "idle", Path: root, LocalState: "manual"})
	if got != "cleared idle" {
		t.Fatalf("ClearRunner = %q", got)
	}
	entries, err := os.ReadDir(work)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("_work should be empty, got %d entries", len(entries))
	}
	if _, err := os.Stat(filepath.Join(root, "actions-runner.zip")); !os.IsNotExist(err) {
		t.Fatalf("installer archive still exists or stat failed unexpectedly: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".runner")); err != nil {
		t.Fatalf(".runner should be preserved: %v", err)
	}
}

func TestClearRepoRunnersFiltersByRepo(t *testing.T) {
	root := t.TempDir()
	got := ClearRepoRunners("SGribanov/A", Inventory{Runners: []Runner{
		{Name: "a", Repo: "SGribanov/A", Path: root, LocalState: "manual"},
		{Name: "b", Repo: "SGribanov/B", Path: root, LocalState: "manual"},
	}})
	if got != "cleared a\n" {
		t.Fatalf("ClearRepoRunners = %q", got)
	}
}

func TestClearNamedRunnerFiltersByName(t *testing.T) {
	root := t.TempDir()
	got := ClearNamedRunner("target", Inventory{Runners: []Runner{
		{Name: "other", Path: root, LocalState: "manual"},
		{Name: "target", Path: root, LocalState: "manual"},
	}})
	if got != "cleared target\n" {
		t.Fatalf("ClearNamedRunner = %q", got)
	}
}

func TestPowerShellQuoteEscapesSingleQuotes(t *testing.T) {
	got := powerShellQuote("runner's")
	if got != "'runner''s'" {
		t.Fatalf("powerShellQuote = %q", got)
	}
}

func TestRemoveRunnerDryRunIncludesSafePlan(t *testing.T) {
	got := RemoveRunner(Runner{
		Name:        "runner-1",
		Repo:        "SGribanov/RunnerMonitor",
		Path:        `C:\Runners\SGribanov-RunnerMonitor\runner-1`,
		ServiceName: "actions.runner.SGribanov-RunnerMonitor.runner-1",
	}, RemoveRunnerOptions{})
	for _, want := range []string{
		"dry-run remove runner-1 for SGribanov/RunnerMonitor",
		"- unregister with GitHub remove token",
		"run with --confirm to execute",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RemoveRunner dry-run missing %q in:\n%s", want, got)
		}
	}
}

func TestRemoveRunnerRejectsBusyWithoutForce(t *testing.T) {
	got := RemoveRunner(Runner{Name: "busy", Busy: true, Path: t.TempDir()}, RemoveRunnerOptions{Confirm: true})
	if got != "busy is busy; removal skipped" {
		t.Fatalf("RemoveRunner = %q", got)
	}
}

func TestRemoveRunnerRefusesUnsafeFolderDelete(t *testing.T) {
	got := RemoveRunner(Runner{Name: "unsafe", Path: t.TempDir()}, RemoveRunnerOptions{Confirm: true, Force: true, DeleteFolder: true})
	if !strings.Contains(got, "folder is outside known runner roots; delete refused") {
		t.Fatalf("RemoveRunner should refuse unsafe delete, got %q", got)
	}
}

func TestAddRunnerRejectsConfiguredFolderWithoutReplace(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	reposRoot := t.TempDir()
	project := filepath.Join(reposRoot, "ProjectA")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("git", "-C", project, "init").CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v: %s", err, out)
	}
	if out, err := exec.Command("git", "-C", project, "remote", "add", "origin", "git@github.com:SGribanov/ProjectA.git").CombinedOutput(); err != nil {
		t.Fatalf("git remote add failed: %v: %s", err, out)
	}
	settingsPath := filepath.Join(t.TempDir(), "runner-monitor.json")
	if err := SaveSettingsAt(settingsPath, Settings{
		ProjectsRoot:       reposRoot,
		WindowsRunnerRoots: []string{`C:\Runners`},
		WSLRunnerRoots:     []string{"/home/gsv777/Runners"},
		LinuxRunnerRoots:   []string{"/opt/Runners"},
	}); err != nil {
		t.Fatal(err)
	}
	t.Setenv("RUNNER_MONITOR_CONFIG", settingsPath)

	runnerFolder := t.TempDir()
	if err := os.WriteFile(filepath.Join(runnerFolder, "config.cmd"), []byte("@echo off\r\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runnerFolder, ".runner"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := AddRunner(AddRunnerOptions{Project: "ProjectA", Name: "runner-a", RunnerFolder: runnerFolder})
	if !strings.Contains(got, "already configured; use --replace") {
		t.Fatalf("AddRunner should reject configured folder without replace, got %q", got)
	}
}

func TestRemoveNamedRunnerFiltersByRepo(t *testing.T) {
	got := RemoveNamedRunner(RemoveRunnerOptions{Name: "same", Repo: "SGribanov/B"}, Inventory{Runners: []Runner{
		{Name: "same", Repo: "SGribanov/A", Path: `C:\Runners\SGribanov-A\same`},
		{Name: "same", Repo: "SGribanov/B", Path: `C:\Runners\SGribanov-B\same`},
	}})
	if !strings.Contains(got, "dry-run remove same for SGribanov/B") {
		t.Fatalf("RemoveNamedRunner = %q", got)
	}
}

func TestProjectRepoFromFolderAt(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	root := t.TempDir()
	project := filepath.Join(root, "ProjectA")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("git", "-C", project, "init").CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v: %s", err, out)
	}
	if out, err := exec.Command("git", "-C", project, "remote", "add", "origin", "git@github.com:SGribanov/ProjectA.git").CombinedOutput(); err != nil {
		t.Fatalf("git remote add failed: %v: %s", err, out)
	}
	got, err := ProjectRepoFromFolderAt(root, "ProjectA")
	if err != nil {
		t.Fatalf("ProjectRepoFromFolderAt returned error: %v", err)
	}
	if got != "SGribanov/ProjectA" {
		t.Fatalf("repo = %q", got)
	}
}

func TestProjectRepoFromFolderAtRejectsPathInput(t *testing.T) {
	_, err := ProjectRepoFromFolderAt(t.TempDir(), `..\ProjectA`)
	if err == nil || !strings.Contains(err.Error(), "project must be a folder name") {
		t.Fatalf("expected folder-name error, got %v", err)
	}
}

func TestShellQuoteEscapesSingleQuotes(t *testing.T) {
	got := shellQuote("runner's")
	if got != "'runner'\"'\"'s'" {
		t.Fatalf("shellQuote = %q", got)
	}
}

func TestLoadSettingsFallsBackToDefaults(t *testing.T) {
	got, err := LoadSettingsAt(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("LoadSettingsAt returned error: %v", err)
	}
	if got.ProjectsRoot != `C:\Repos` {
		t.Fatalf("ProjectsRoot = %q", got.ProjectsRoot)
	}
	if len(got.WindowsRunnerRoots) != 1 || got.WindowsRunnerRoots[0] != `C:\Runners` {
		t.Fatalf("WindowsRunnerRoots = %#v", got.WindowsRunnerRoots)
	}
	if len(got.WSLRunnerRoots) != 1 || got.WSLRunnerRoots[0] != "/home/gsv777/Runners" {
		t.Fatalf("WSLRunnerRoots = %#v", got.WSLRunnerRoots)
	}
	if got.TUIRefreshIntervalSeconds != 5 {
		t.Fatalf("TUIRefreshIntervalSeconds = %d", got.TUIRefreshIntervalSeconds)
	}
}

func TestLoadSettingsReadsSudoPasswordValue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "runner-monitor.json")
	if err := SaveSettingsAt(path, Settings{
		ProjectsRoot:              `D:\Repos`,
		WindowsRunnerRoots:        []string{`D:\Runners`},
		WSLRunnerRoots:            []string{"/runnerbox/Runners"},
		LinuxRunnerRoots:          []string{"/opt/Runners"},
		TUIRefreshIntervalSeconds: 9,
		WSLSudoPassword:           "secret",
	}); err != nil {
		t.Fatal(err)
	}
	got, err := LoadSettingsAt(path)
	if err != nil {
		t.Fatalf("LoadSettingsAt returned error: %v", err)
	}
	if got.WSLSudoPassword != "secret" {
		t.Fatalf("WSLSudoPassword was not loaded")
	}
	if got.ProjectsRoot != `D:\Repos` || got.WindowsRunnerRoots[0] != `D:\Runners` {
		t.Fatalf("settings not loaded: %#v", got)
	}
	if got.TUIRefreshIntervalSeconds != 9 {
		t.Fatalf("TUIRefreshIntervalSeconds = %d", got.TUIRefreshIntervalSeconds)
	}
}

func TestLoadSettingsNormalizesInvalidTUIRefreshInterval(t *testing.T) {
	path := filepath.Join(t.TempDir(), "runner-monitor.json")
	if err := os.WriteFile(path, []byte(`{"tuiRefreshIntervalSeconds":0}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := LoadSettingsAt(path)
	if err != nil {
		t.Fatalf("LoadSettingsAt returned error: %v", err)
	}
	if got.TUIRefreshIntervalSeconds != 5 {
		t.Fatalf("TUIRefreshIntervalSeconds = %d", got.TUIRefreshIntervalSeconds)
	}
}

func TestRenderSettingsMasksSudoPassword(t *testing.T) {
	got := RenderSettings(Settings{WSLSudoPassword: "secret"})
	if strings.Contains(got, "secret") {
		t.Fatalf("RenderSettings leaked password: %s", got)
	}
	if !strings.Contains(got, `"wslSudoPassword": "<set>"`) {
		t.Fatalf("RenderSettings did not show masked password: %s", got)
	}
	empty := RenderSettings(Settings{})
	if !strings.Contains(empty, `"wslSudoPassword": "<empty>"`) {
		t.Fatalf("RenderSettings did not show empty password: %s", empty)
	}
}

func TestSafeRunnerRootUsesConfiguredRoots(t *testing.T) {
	settingsPath := filepath.Join(t.TempDir(), "runner-monitor.json")
	if err := SaveSettingsAt(settingsPath, Settings{
		ProjectsRoot:       `C:\Repos`,
		WindowsRunnerRoots: []string{`D:\BuildAgents`},
		WSLRunnerRoots:     []string{"/runnerbox/Runners"},
		LinuxRunnerRoots:   []string{"/opt/Runners"},
	}); err != nil {
		t.Fatal(err)
	}
	t.Setenv("RUNNER_MONITOR_CONFIG", settingsPath)
	if !isSafeRunnerRoot(Runner{Path: `D:\BuildAgents\SGribanov-RunnerMonitor\runner-1`}) {
		t.Fatalf("configured Windows runner root should be safe")
	}
	if isSafeRunnerRoot(Runner{Path: `C:\Runners\SGribanov-RunnerMonitor\runner-1`}) {
		t.Fatalf("unconfigured Windows runner root should not be safe")
	}
	if !isSafeRunnerRoot(Runner{Transport: "wsl", Path: "/runnerbox/Runners/SGribanov-RunnerMonitor/runner-1"}) {
		t.Fatalf("configured WSL runner root should be safe")
	}
	if isSafeRunnerRoot(Runner{Transport: "wsl", Path: "/runnerbox/Runners/../danger"}) {
		t.Fatalf("WSL path traversal should not be safe")
	}
	if isSafeRunnerRoot(Runner{Transport: "wsl", Path: "/runnerbox/Runners"}) {
		t.Fatalf("configured WSL root itself should not be deletable")
	}
	if isSafeRunnerRoot(Runner{Path: `D:\BuildAgents`}) {
		t.Fatalf("configured Windows root itself should not be deletable")
	}
	if isSafeRunnerRoot(Runner{Path: `D:\BuildAgents\..\danger`}) {
		t.Fatalf("Windows path traversal should not be safe")
	}
}

func TestConfigPowerShellScriptWithTokenUsesEnvironmentReference(t *testing.T) {
	got := configPowerShellScriptWithToken(`C:\Runners\runner 1\config.cmd`, []string{"remove", "--unattended", "--token", "secret-token"})
	if strings.Contains(got, "secret-token") {
		t.Fatalf("PowerShell script should not contain raw token: %q", got)
	}
	if !strings.Contains(got, "$env:RUNNER_MONITOR_RUNNER_TOKEN") {
		t.Fatalf("PowerShell script should use token env reference: %q", got)
	}
}

func TestConfigShellLineWithTokenUsesEnvironmentReference(t *testing.T) {
	got := configShellLineWithToken(`/opt/Runner Monitor/config.sh`, []string{"remove", "--unattended", "--token", "secret-token"})
	if strings.Contains(got, "secret-token") {
		t.Fatalf("shell line should not contain raw token: %q", got)
	}
	if !strings.Contains(got, "$RUNNER_MONITOR_RUNNER_TOKEN") {
		t.Fatalf("shell line should use token env reference: %q", got)
	}
}

func TestMaskSecretRedactsToken(t *testing.T) {
	got := maskSecret("config failed with secret-token", "secret-token")
	if strings.Contains(got, "secret-token") || !strings.Contains(got, "<redacted>") {
		t.Fatalf("maskSecret did not redact token: %q", got)
	}
}

func TestAppendWSLEnvPreservesExistingEntries(t *testing.T) {
	got := appendWSLEnv("PATH/l:OTHER", "RUNNER_MONITOR_RUNNER_TOKEN")
	if got != "PATH/l:OTHER:RUNNER_MONITOR_RUNNER_TOKEN" {
		t.Fatalf("appendWSLEnv = %q", got)
	}
	if duplicated := appendWSLEnv(got, "RUNNER_MONITOR_RUNNER_TOKEN"); duplicated != got {
		t.Fatalf("appendWSLEnv duplicated entry: %q", duplicated)
	}
}

func TestRefreshWithGitHubCacheUsesCachedGitHubStatus(t *testing.T) {
	previousDiscoverLocal := discoverLocal
	discoverLocal = func() ([]Runner, error) {
		return nil, nil
	}
	t.Cleanup(func() {
		discoverLocal = previousDiscoverLocal
	})

	calls := 0
	loader := func(repos []string) (map[string]GitHubRunnerStatus, map[string]QueueStatus, error) {
		calls++
		return map[string]GitHubRunnerStatus{
			runnerKey("SGribanov/RunnerMonitor", "runner-1"): {Status: "online", Busy: true},
		}, map[string]QueueStatus{"SGribanov/RunnerMonitor": {Count: 1}}, nil
	}
	inventory, err := refreshWithGitHubStatus(loader)
	if err != nil {
		t.Fatalf("refreshWithGitHubStatus returned error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("loader calls = %d", calls)
	}
	if len(inventory.Runners) > 0 {
		t.Skip("local runner inventory is environment-dependent")
	}
}

func TestAuditRunnerCandidateRemove(t *testing.T) {
	decision, _ := AuditRunner(Runner{Name: "old", LocalState: "manual", GitHubStatus: "unknown"})
	if decision != "candidate-remove" {
		t.Fatalf("decision = %q", decision)
	}
}

func TestAuditRunnerInvestigatesQueuedRepo(t *testing.T) {
	decision, _ := AuditRunner(Runner{Name: "queued", LocalState: "manual", GitHubStatus: "unknown", QueueCount: 1})
	if decision != "investigate" {
		t.Fatalf("decision = %q", decision)
	}
}

func TestAuditRunnerInvestigatesManualRunningRunner(t *testing.T) {
	decision, evidence := AuditRunner(Runner{Name: "manual-running", LocalState: "running", GitHubStatus: "online", ControlMode: "manual"})
	if decision != "investigate" || evidence != "online in GitHub but not service-managed locally" {
		t.Fatalf("decision/evidence = %q/%q", decision, evidence)
	}
}

func TestAuditRunnerKeepsControllableManualWindowsRunner(t *testing.T) {
	decision, evidence := AuditRunner(Runner{
		Name:         "manual-windows",
		LocalState:   "running",
		GitHubStatus: "online",
		ControlMode:  "manual",
		Transport:    "windows",
	})
	if decision != "keep" || evidence != "manual Windows runner is controllable by RunnerMonitor" {
		t.Fatalf("decision/evidence = %q/%q", decision, evidence)
	}
}

func TestAuditRunnerFlagsUnitOnlyRunner(t *testing.T) {
	decision, evidence := AuditRunner(Runner{Name: "unit-only", Path: "(unit only)", LocalState: "inactive"})
	if decision != "candidate-remove" || evidence != "orphan service unit without runner directory" {
		t.Fatalf("decision/evidence = %q/%q", decision, evidence)
	}
}

func TestAuditPolicyKeepsRunner(t *testing.T) {
	policy := AuditPolicy{Keep: []AuditPolicyRule{{Repo: "SGribanov/AU", Runner: "windows-local", Reason: "needed"}}}
	decision, evidence := AuditRunnerWithPolicy(Runner{Repo: "SGribanov/AU", Name: "windows-local", LocalState: "manual"}, policy)
	if decision != "keep" || evidence != "policy keep: needed" {
		t.Fatalf("decision/evidence = %q/%q", decision, evidence)
	}
}

func TestRepoAndRunnerFromActionsService(t *testing.T) {
	repo, name := repoAndRunnerFromActionsService("actions.runner.SGribanov-NewGenOsEngine.newgen-wsl-linux.service")
	if repo != "SGribanov/NewGenOsEngine" || name != "newgen-wsl-linux" {
		t.Fatalf("repo/name = %q/%q", repo, name)
	}
}

func TestRunRepoLifecycleSkipsManualRunner(t *testing.T) {
	got := RunRepoLifecycle("start", "SGribanov/RunnerMonitor", Inventory{Runners: []Runner{{
		Name: "manual", Repo: "SGribanov/RunnerMonitor",
	}}})
	if got != "skip manual: not controllable\nno controllable runners found for SGribanov/RunnerMonitor\n" {
		t.Fatalf("RunRepoLifecycle = %q", got)
	}
}

func TestRunnerDirFromProcessPath(t *testing.T) {
	got := runnerDirFromProcessPath(`C:\actions-runner-backtester\bin\Runner.Listener.exe`)
	if got != `C:\actions-runner-backtester` {
		t.Fatalf("runnerDirFromProcessPath = %q", got)
	}
}

func TestRemoteWindowsTUICommand(t *testing.T) {
	host := RemoteHost{SSHHost: "runnerbox", OS: "windows", RunnerMonitorPath: "C:/Repos/Runner Monitor/runner-monitor.ps1"}
	args := remoteTUISSHArgs(host)
	wantCommand := "powershell -NoProfile -ExecutionPolicy Bypass -File 'C:/Repos/Runner Monitor/runner-monitor.ps1'"
	if len(args) != 3 || args[0] != "-t" || args[1] != "runnerbox" || args[2] != wantCommand {
		t.Fatalf("remoteTUISSHArgs = %#v", args)
	}
}

func TestRemoteLinuxTUICommandQuotesPath(t *testing.T) {
	host := RemoteHost{SSHHost: "runnerbox", OS: "linux", RunnerMonitorPath: "/opt/Runner Monitor/runner-monitor"}
	args := remoteTUISSHArgs(host)
	wantCommand := "'/opt/Runner Monitor/runner-monitor'"
	if len(args) != 3 || args[2] != wantCommand {
		t.Fatalf("remoteTUISSHArgs = %#v", args)
	}
}

func TestWindowsDiscoveryPowerShellTimeoutReturnsWarning(t *testing.T) {
	original := windowsDiscoveryTimeout
	windowsDiscoveryTimeout = time.Nanosecond
	t.Cleanup(func() { windowsDiscoveryTimeout = original })

	_, err := runWindowsDiscoveryPowerShell("Start-Sleep -Seconds 5")
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("expected timeout warning, got %v", err)
	}
}

func TestPromptRemoteHostUsesLinuxDefaultsAfterOSSelection(t *testing.T) {
	input := strings.NewReader("runnerlinux\nrunnerlinux\nlinux\n\n\n")
	var output strings.Builder
	host, err := promptRemoteHost("", RemoteHost{}, input, &output)
	if err != nil {
		t.Fatalf("promptRemoteHost returned error: %v", err)
	}
	if host.RunnerMonitorPath != "/opt/RunnerMonitor/runner-monitor" {
		t.Fatalf("RunnerMonitorPath = %q", host.RunnerMonitorPath)
	}
	if host.DefaultProjectPath != "/srv/DeltaG" {
		t.Fatalf("DefaultProjectPath = %q", host.DefaultProjectPath)
	}
}
