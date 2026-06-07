package app

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	startLocalTimeout  = 15 * time.Second
	startGitHubTimeout = 60 * time.Second
	startPollInterval  = 2 * time.Second
	stopLocalTimeout   = 15 * time.Second
	stopGitHubTimeout  = 60 * time.Second
)

var (
	runServiceAction       = defaultRunServiceAction
	serviceState           = defaultServiceState
	loadRunnerGitHubStatus = loadRunnerGitHubStatusFromAPI
	sleepForLifecyclePoll  = time.Sleep
)

func RunLifecycle(action string, runner Runner) string {
	if runner.IsGitHubHosted() {
		return fmt.Sprintf("%s is GitHub-hosted and read-only; cannot %s", runner.Name, action)
	}
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
	if runner.ControlMode == "manual" && runner.Transport == "windows" {
		if action == "start" && isAlreadyRunning(runner.LocalState) {
			return fmt.Sprintf("%s already running", runner.Name)
		}
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
	if needsElevatedWindowsLifecycle(runner) && !isElevatedWindowsProcess() {
		if err := launchElevatedLifecycleRunner(action, runner); err != nil {
			return fmt.Sprintf("%s %s requires elevated PowerShell, but UAC launch failed: %v", action, runner.Name, err)
		}
		return fmt.Sprintf("%s %s requested in elevated PowerShell", action, runner.Name)
	}

	if action == "start" {
		if err := startServiceManagedRunner(runner); err != nil {
			return fmt.Sprintf("start %s failed: %v", runner.Name, err)
		}
		return fmt.Sprintf("%s start requested; service active; GitHub online", runner.Name)
	}
	if action == "stop" {
		if err := stopServiceManagedRunner(runner); err != nil {
			return fmt.Sprintf("stop %s failed: %v", runner.Name, err)
		}
		return fmt.Sprintf("%s stop requested; service stopped; GitHub offline", runner.Name)
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

func startServiceManagedRunner(runner Runner) error {
	if supportsServiceEnable(runner.ControlMode) {
		if err := runServiceAction(runner.ControlMode, "enable", runner.ServiceName); err != nil {
			return fmt.Errorf("enable service %s: %w", runner.ServiceName, err)
		}
	}
	if err := runServiceAction(runner.ControlMode, "start", runner.ServiceName); err != nil {
		return fmt.Errorf("start service %s: %w", runner.ServiceName, err)
	}
	if err := waitForServiceActive(runner.ControlMode, runner.ServiceName, startLocalTimeout, startPollInterval); err != nil {
		return err
	}
	if err := waitForGitHubRunnerOnline(runner.Repo, runner.Name, startGitHubTimeout, startPollInterval); err != nil {
		return err
	}
	return nil
}

func stopServiceManagedRunner(runner Runner) error {
	if err := runServiceAction(runner.ControlMode, "stop", runner.ServiceName); err != nil {
		return fmt.Errorf("stop service %s: %w", runner.ServiceName, err)
	}
	if err := waitForServiceStopped(runner.ControlMode, runner.ServiceName, stopLocalTimeout, startPollInterval); err != nil {
		return err
	}
	if err := waitForGitHubRunnerOffline(runner.Repo, runner.Name, stopGitHubTimeout, startPollInterval); err != nil {
		return err
	}
	return nil
}

func supportsServiceEnable(controlMode string) bool {
	return controlMode == "wsl-systemd" || controlMode == "systemd"
}

func waitForServiceActive(controlMode, serviceName string, timeout, interval time.Duration) error {
	attempts := lifecyclePollAttempts(timeout, interval)
	lastState := ""
	for i := 0; i < attempts; i++ {
		state, err := serviceState(controlMode, serviceName)
		if err == nil && isAlreadyRunning(state) {
			return nil
		}
		if strings.TrimSpace(state) != "" {
			lastState = strings.TrimSpace(state)
		}
		if i < attempts-1 {
			sleepForLifecyclePoll(interval)
		}
	}
	if lastState == "" {
		lastState = "unknown"
	}
	return fmt.Errorf("service %s did not become active; last state: %s", serviceName, lastState)
}

func waitForServiceStopped(controlMode, serviceName string, timeout, interval time.Duration) error {
	attempts := lifecyclePollAttempts(timeout, interval)
	lastState := ""
	for i := 0; i < attempts; i++ {
		state, err := serviceState(controlMode, serviceName)
		trimmed := strings.TrimSpace(state)
		if err == nil && !isAlreadyRunning(trimmed) {
			return nil
		}
		if trimmed != "" {
			lastState = trimmed
		}
		if i < attempts-1 {
			sleepForLifecyclePoll(interval)
		}
	}
	if lastState == "" {
		lastState = "unknown"
	}
	return fmt.Errorf("service %s did not stop; last state: %s", serviceName, lastState)
}

func waitForGitHubRunnerOnline(repo, name string, timeout, interval time.Duration) error {
	if strings.TrimSpace(repo) == "" || strings.TrimSpace(name) == "" {
		return fmt.Errorf("cannot verify GitHub online status without repo and runner name")
	}
	attempts := lifecyclePollAttempts(timeout, interval)
	lastStatus := ""
	var lastErr error
	for i := 0; i < attempts; i++ {
		status, err := loadRunnerGitHubStatus(repo, name)
		if err == nil && strings.EqualFold(status, "online") {
			return nil
		}
		if err != nil {
			lastErr = err
		}
		if strings.TrimSpace(status) != "" {
			lastStatus = strings.TrimSpace(status)
		}
		if i < attempts-1 {
			sleepForLifecyclePoll(interval)
		}
	}
	if lastStatus == "" {
		lastStatus = "unknown"
	}
	if lastErr != nil {
		return fmt.Errorf("GitHub runner %s/%s did not become online; last status: %s; last error: %v", repo, name, lastStatus, lastErr)
	}
	return fmt.Errorf("GitHub runner %s/%s did not become online; last status: %s", repo, name, lastStatus)
}

func waitForGitHubRunnerOffline(repo, name string, timeout, interval time.Duration) error {
	if strings.TrimSpace(repo) == "" || strings.TrimSpace(name) == "" {
		return fmt.Errorf("cannot verify GitHub offline status without repo and runner name")
	}
	attempts := lifecyclePollAttempts(timeout, interval)
	lastStatus := ""
	var lastErr error
	for i := 0; i < attempts; i++ {
		status, err := loadRunnerGitHubStatus(repo, name)
		if err == nil && !strings.EqualFold(status, "online") {
			return nil
		}
		if err != nil {
			lastErr = err
		}
		if strings.TrimSpace(status) != "" {
			lastStatus = strings.TrimSpace(status)
		}
		if i < attempts-1 {
			sleepForLifecyclePoll(interval)
		}
	}
	if lastStatus == "" {
		lastStatus = "unknown"
	}
	if lastErr != nil {
		return fmt.Errorf("GitHub runner %s/%s did not become offline; last status: %s; last error: %v", repo, name, lastStatus, lastErr)
	}
	return fmt.Errorf("GitHub runner %s/%s did not become offline; last status: %s", repo, name, lastStatus)
}

func lifecyclePollAttempts(timeout, interval time.Duration) int {
	if interval <= 0 || timeout <= 0 {
		return 1
	}
	attempts := int(timeout / interval)
	if timeout%interval != 0 {
		attempts++
	}
	if attempts < 1 {
		return 1
	}
	return attempts
}

func defaultRunServiceAction(controlMode, action, serviceName string) error {
	switch controlMode {
	case "windows-service":
		if runtime.GOOS != "windows" {
			return fmt.Errorf("Windows service control is only available on Windows")
		}
		return runWindowsService(action, serviceName)
	case "wsl-systemd":
		return runWSLService(action, serviceName)
	case "systemd":
		return runCommandWithOutput("systemctl", action, serviceName)
	default:
		return fmt.Errorf("unsupported control mode %q", controlMode)
	}
}

func defaultServiceState(controlMode, serviceName string) (string, error) {
	switch controlMode {
	case "windows-service":
		return windowsServiceState(serviceName)
	case "wsl-systemd":
		return wslServiceState(serviceName), nil
	case "systemd":
		out, err := exec.Command("systemctl", "is-active", serviceName).Output()
		return strings.TrimSpace(string(out)), err
	default:
		return "", fmt.Errorf("unsupported control mode %q", controlMode)
	}
}

func loadRunnerGitHubStatusFromAPI(repo, name string) (string, error) {
	data, err := ghAPI(fmt.Sprintf("repos/%s/actions/runners", repo))
	if err != nil {
		return "", err
	}
	response, err := parseRunnersResponse(data)
	if err != nil {
		return "", err
	}
	for _, runner := range response.Runners {
		if strings.EqualFold(runner.Name, name) {
			return runner.Status, nil
		}
	}
	return "missing", nil
}

func OpenLogs(runner Runner) string {
	if runner.IsGitHubHosted() {
		return fmt.Sprintf("%s is GitHub-hosted; open workflow logs in GitHub: %s", runner.Name, emptyAsDash(runner.Path))
	}
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
		if runner.IsGitHubHosted() {
			fmt.Fprintf(&b, "skip %s: GitHub-hosted read-only\n", runner.Name)
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

func RunNamedLifecycle(action string, name string, repo string, inventory Inventory) string {
	for _, runner := range inventory.Runners {
		if !strings.EqualFold(runner.Name, name) {
			continue
		}
		if strings.TrimSpace(repo) != "" && !strings.EqualFold(runner.Repo, repo) {
			continue
		}
		return RunLifecycle(action, runner) + "\n"
	}
	if strings.TrimSpace(repo) != "" {
		return fmt.Sprintf("no runner found named %s for %s\n", name, repo)
	}
	return fmt.Sprintf("no runner found named %s\n", name)
}

func needsElevatedWindowsLifecycle(runner Runner) bool {
	return runtime.GOOS == "windows" && runner.ControlMode == "windows-service"
}

func launchElevatedLifecycleRunner(action string, runner Runner) error {
	target, args := elevatedLifecycleRunnerTarget(action, runner)
	ps := fmt.Sprintf(
		"$argsList = @(%s); Start-Process -FilePath %s -ArgumentList $argsList -Verb RunAs",
		powerShellArray(args),
		powerShellQuote(target),
	)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", ps)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func elevatedLifecycleRunnerTarget(action string, runner Runner) (string, []string) {
	args := []string{
		"--" + action + "-runner",
		runner.Name,
		"--repo",
		runner.Repo,
	}
	if script := runnerMonitorScriptPath(); script != "" {
		return "powershell.exe", append([]string{
			"-NoExit",
			"-NoProfile",
			"-ExecutionPolicy",
			"Bypass",
			"-File",
			script,
		}, args...)
	}
	return os.Args[0], args
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
