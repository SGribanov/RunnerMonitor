package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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

func TestRunLifecycleRejectsManualRunner(t *testing.T) {
	got := RunLifecycle("start", Runner{Name: "manual-runner"})
	want := "manual-runner is not service-managed; cannot start"
	if got != want {
		t.Fatalf("RunLifecycle = %q, want %q", got, want)
	}
}

func TestRunLifecycleProtectsBusyRunner(t *testing.T) {
	got := RunLifecycle("stop", Runner{Name: "busy-runner", Busy: true, ServiceName: "svc", ControlMode: "unsupported"})
	want := "busy-runner is busy; use force-stop to override"
	if got != want {
		t.Fatalf("RunLifecycle = %q, want %q", got, want)
	}
}

func TestRunLifecycleStartAlreadyRunningDoesNotRequireServiceAccess(t *testing.T) {
	got := RunLifecycle("start", Runner{Name: "running-runner", LocalState: "running", ServiceName: "svc", ControlMode: "unsupported"})
	want := "running-runner already running"
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
	oldRoot := defaultReposRoot
	defaultReposRoot = reposRoot
	t.Cleanup(func() { defaultReposRoot = oldRoot })

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
	host := RemoteHost{SSHHost: "runnerbox", OS: "windows", RunnerMonitorPath: "C:/Repos/RunnerMonitor/runner-monitor.ps1"}
	args := remoteTUISSHArgs(host)
	wantCommand := "powershell -NoProfile -ExecutionPolicy Bypass -File C:/Repos/RunnerMonitor/runner-monitor.ps1"
	if len(args) != 3 || args[0] != "-t" || args[1] != "runnerbox" || args[2] != wantCommand {
		t.Fatalf("remoteTUISSHArgs = %#v", args)
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
