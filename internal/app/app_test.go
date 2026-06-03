package app

import (
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
