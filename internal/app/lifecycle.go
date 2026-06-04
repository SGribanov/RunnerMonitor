package app

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func RunLifecycle(action string, runner Runner) string {
	force := false
	switch action {
	case "force-stop":
		action = "stop"
		force = true
	case "force-restart":
		action = "restart"
		force = true
	}
	if runner.Busy && !force && (action == "stop" || action == "restart") {
		return fmt.Sprintf("%s is busy; use force-%s to override", runner.Name, action)
	}
	if action == "start" && isAlreadyRunning(runner.LocalState) {
		return fmt.Sprintf("%s already running", runner.Name)
	}
	if runner.ControlMode == "manual" && runner.Transport == "windows" {
		if runtime.GOOS != "windows" {
			return "manual Windows runner control is only available on Windows"
		}
		if runner.Path == "" {
			return fmt.Sprintf("%s path is unknown; cannot %s", runner.Name, action)
		}
		if err := runWindowsManualRunner(action, runner.Path); err != nil {
			return fmt.Sprintf("%s %s failed: %v", action, runner.Name, err)
		}
		return fmt.Sprintf("%s %s requested", action, runner.Name)
	}
	if runner.ServiceName == "" {
		return fmt.Sprintf("%s is not service-managed; cannot %s", runner.Name, action)
	}

	var err error
	switch runner.ControlMode {
	case "windows-service":
		if runtime.GOOS != "windows" {
			return "Windows service control is only available on Windows"
		}
		err = runWindowsService(action, runner.ServiceName)
	case "wsl-systemd":
		err = runWSLService(action, runner.ServiceName)
	case "systemd":
		err = runCommandWithOutput("systemctl", action, runner.ServiceName)
	default:
		return fmt.Sprintf("%s has unsupported control mode %q", runner.Name, runner.ControlMode)
	}
	if err != nil {
		return fmt.Sprintf("%s %s failed: %v", action, runner.Name, err)
	}
	return fmt.Sprintf("%s %s requested", action, runner.Name)
}

func OpenLogs(runner Runner) string {
	if runner.Path == "" {
		return "runner path is unknown"
	}
	return fmt.Sprintf("logs: %s/_diag", runner.Path)
}

func RunRepoLifecycle(action string, repo string, inventory Inventory) string {
	var b strings.Builder
	count := 0
	for _, runner := range inventory.Runners {
		if !strings.EqualFold(runner.Repo, repo) {
			continue
		}
		if runner.ServiceName == "" && !(runner.ControlMode == "manual" && runner.Transport == "windows") {
			fmt.Fprintf(&b, "skip %s: not controllable\n", runner.Name)
			continue
		}
		count++
		fmt.Fprintf(&b, "%s\n", RunLifecycle(action, runner))
	}
	if count == 0 {
		fmt.Fprintf(&b, "no controllable runners found for %s\n", repo)
	}
	return b.String()
}

func DisableAutostart(inventory Inventory) string {
	var b strings.Builder
	count := 0
	for _, runner := range inventory.Runners {
		if runner.ServiceName == "" {
			continue
		}
		count++
		var err error
		switch runner.ControlMode {
		case "windows-service":
			err = disableWindowsServiceAutostart(runner.ServiceName)
		case "wsl-systemd":
			err = disableWSLServiceAutostart(runner.ServiceName)
		case "systemd":
			err = runCommandWithOutput("systemctl", "disable", runner.ServiceName)
		default:
			fmt.Fprintf(&b, "skip %s: unsupported control mode %q\n", runner.Name, runner.ControlMode)
			continue
		}
		if err != nil {
			fmt.Fprintf(&b, "failed %s (%s): %v\n", runner.Name, runner.ServiceName, err)
			continue
		}
		fmt.Fprintf(&b, "disabled autostart: %s (%s)\n", runner.Name, runner.ServiceName)
	}
	if count == 0 {
		b.WriteString("no service-managed runners discovered\n")
	}
	return b.String()
}

func runCommandWithOutput(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func isAlreadyRunning(state string) bool {
	state = strings.ToLower(strings.TrimSpace(state))
	return state == "running" || state == "active"
}
