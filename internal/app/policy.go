package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type AuditPolicy struct {
	Keep []AuditPolicyRule `json:"keep"`
}

type AuditPolicyRule struct {
	Repo   string `json:"repo"`
	Runner string `json:"runner"`
	Reason string `json:"reason"`
}

func LoadAuditPolicy() AuditPolicy {
	for _, path := range auditPolicyPaths() {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var policy AuditPolicy
		if err := json.Unmarshal(data, &policy); err == nil {
			return policy
		}
	}
	return AuditPolicy{}
}

func auditPolicyPaths() []string {
	var paths []string
	if explicit := os.Getenv("RUNNER_MONITOR_POLICY"); explicit != "" {
		paths = append(paths, explicit)
	}
	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(cwd, "runner-policy.json"))
	}
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		paths = append(paths, filepath.Join(exeDir, "runner-policy.json"))
		paths = append(paths, filepath.Join(filepath.Dir(exeDir), "runner-policy.json"))
	}
	paths = append(paths, `C:\Repos\RunnerMonitor\runner-policy.json`)
	return paths
}

func (p AuditPolicy) KeepReason(r Runner) string {
	for _, rule := range p.Keep {
		if strings.EqualFold(rule.Repo, r.Repo) && strings.EqualFold(rule.Runner, r.Name) {
			return rule.Reason
		}
	}
	return ""
}
