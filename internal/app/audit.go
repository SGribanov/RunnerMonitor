package app

import "strings"

func AuditRunner(r Runner) (decision string, evidence string) {
	return AuditRunnerWithPolicy(r, AuditPolicy{})
}

func AuditRunnerWithPolicy(r Runner, policy AuditPolicy) (decision string, evidence string) {
	if reason := policy.KeepReason(r); reason != "" {
		return "keep", "policy keep: " + reason
	}
	if r.Busy {
		return "keep", "runner is currently busy"
	}
	if r.Path == "(unit only)" {
		return "candidate-remove", "orphan service unit without runner directory"
	}
	if r.QueueCount > 0 {
		return "investigate", "repo has queued jobs; check labels/routes before removal"
	}
	if r.GitHubStatus == "online" && r.ControlMode == "manual" {
		if r.Transport == "windows" && r.LocalState == "running" {
			return "keep", "manual Windows runner is controllable by RunnerMonitor"
		}
		return "investigate", "online in GitHub but not service-managed locally"
	}
	if r.GitHubStatus == "online" && (r.LocalState == "running" || r.LocalState == "active") {
		return "keep", "service-managed and online"
	}
	if r.GitHubStatus == "unknown" && isInactiveOrManual(r.LocalState) {
		return "candidate-remove", "local configured runner not visible in GitHub API"
	}
	if r.GitHubStatus == "offline" && isInactiveOrManual(r.LocalState) {
		return "candidate-remove", "offline/manual runner"
	}
	if r.LocalState == "inactive" {
		return "candidate-remove", "inactive service"
	}
	return "investigate", "state needs manual review"
}

func isInactiveOrManual(state string) bool {
	state = strings.ToLower(state)
	return state == "manual" || state == "inactive" || state == "configured" || state == "unknown" || state == ""
}
