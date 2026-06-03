package app

import "testing"

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

func TestRunRepoLifecycleSkipsManualRunner(t *testing.T) {
	got := RunRepoLifecycle("start", "SGribanov/RunnerMonitor", Inventory{Runners: []Runner{{
		Name: "manual", Repo: "SGribanov/RunnerMonitor",
	}}})
	if got != "skip manual: not service-managed\nno service-managed runners found for SGribanov/RunnerMonitor\n" {
		t.Fatalf("RunRepoLifecycle = %q", got)
	}
}
