package app

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
)

func discoverWSLRunners() ([]Runner, error) {
	cmd := exec.Command("wsl.exe", "sh", "-lc", "find /home /opt /srv -name .runner -type f 2>/dev/null | head -200")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var runners []Runner
	for _, file := range strings.Fields(string(out)) {
		data, err := exec.Command("wsl.exe", "cat", file).Output()
		if err != nil {
			continue
		}
		config, err := parseRunnerConfig(data)
		if err != nil {
			continue
		}
		dir := path.Dir(file)
		runner := runnerFromConfig(config, dir, "wsl:Ubuntu", "wsl")
		if service := wslCat(path.Join(dir, ".service")); service != "" {
			runner.ServiceName = strings.TrimSpace(service)
			runner.ControlMode = "wsl-systemd"
			runner.LocalState = wslServiceState(runner.ServiceName)
		}
		runners = append(runners, runner)
	}
	return runners, nil
}

func wslCat(file string) string {
	out, err := exec.Command("wsl.exe", "cat", file).Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func wslServiceState(serviceName string) string {
	out, err := exec.Command("wsl.exe", "systemctl", "is-active", serviceName).Output()
	if err != nil {
		text := strings.TrimSpace(string(out))
		if text == "" {
			return "inactive"
		}
		return text
	}
	return strings.TrimSpace(string(out))
}

func runWSLService(action string, serviceName string) error {
	cmd := exec.Command("wsl.exe", "systemctl", action, serviceName)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func disableWSLServiceAutostart(serviceName string) error {
	cmd := exec.Command("wsl.exe", "systemctl", "disable", serviceName)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
