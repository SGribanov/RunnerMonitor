package app

import (
	"bytes"
	"fmt"
	"os"
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
	runners = append(runners, discoverWSLUnitOnlyRunners(runners)...)
	return runners, nil
}

func discoverWSLUnitOnlyRunners(existing []Runner) []Runner {
	knownServices := map[string]bool{}
	for _, runner := range existing {
		if runner.ServiceName != "" {
			knownServices[runner.ServiceName] = true
		}
	}

	out, err := exec.Command("wsl.exe", "sh", "-lc", "systemctl list-unit-files 'actions.runner.*.service' --no-legend --no-pager").Output()
	if err != nil {
		return nil
	}

	var runners []Runner
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		serviceName := fields[0]
		if knownServices[serviceName] {
			continue
		}
		repo, name := repoAndRunnerFromActionsService(serviceName)
		if repo == "" || name == "" {
			continue
		}
		runners = append(runners, Runner{
			Name:         name,
			Repo:         repo,
			Host:         "wsl:Ubuntu",
			Transport:    "wsl",
			Path:         "(unit only)",
			ServiceName:  serviceName,
			ControlMode:  "wsl-systemd",
			LocalState:   wslServiceState(serviceName),
			GitHubStatus: "unknown",
		})
	}
	return runners
}

func repoAndRunnerFromActionsService(serviceName string) (string, string) {
	name := strings.TrimSuffix(serviceName, ".service")
	name = strings.TrimPrefix(name, "actions.runner.")
	parts := strings.SplitN(name, ".", 2)
	if len(parts) != 2 {
		return "", ""
	}
	ownerRepo := strings.SplitN(parts[0], "-", 2)
	if len(ownerRepo) != 2 {
		return "", ""
	}
	return ownerRepo[0] + "/" + ownerRepo[1], parts[1]
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
		return runWSLServiceWithSudo(action, serviceName, err, out)
	}
	return nil
}

func runWSLServiceWithSudo(action string, serviceName string, originalErr error, originalOut []byte) error {
	passwordFile := wslSudoPasswordFileWindowsPath()
	password, readErr := os.ReadFile(passwordFile)
	if readErr != nil {
		originalText := strings.TrimSpace(string(originalOut))
		return fmt.Errorf("%w: %s; sudo fallback failed: cannot read sudo password file %s: %v", originalErr, originalText, passwordFile, readErr)
	}
	cmd := exec.Command("wsl.exe", "--", "sudo", "-S", "-p", "", "systemctl", action, serviceName)
	cmd.Stdin = bytes.NewReader(password)
	if out, err := cmd.CombinedOutput(); err != nil {
		originalText := strings.TrimSpace(string(originalOut))
		sudoText := strings.TrimSpace(string(out))
		return fmt.Errorf("%w: %s; sudo fallback failed: %v: %s", originalErr, originalText, err, sudoText)
	}
	return nil
}

func wslSudoPasswordFileWindowsPath() string {
	passwordFile := os.Getenv("RUNNER_MONITOR_WSL_SUDO_FILE")
	if passwordFile == "" {
		return `C:\Users\gsv777\Desktop\WSL_sudo.txt`
	}
	if strings.HasPrefix(passwordFile, "/mnt/") && len(passwordFile) > len("/mnt/c/") {
		drive := strings.ToUpper(passwordFile[5:6])
		rest := strings.TrimPrefix(passwordFile[7:], "/")
		return drive + `:\` + strings.ReplaceAll(rest, "/", `\`)
	}
	return passwordFile
}

func disableWSLServiceAutostart(serviceName string) error {
	cmd := exec.Command("wsl.exe", "systemctl", "disable", serviceName)
	if out, err := cmd.CombinedOutput(); err != nil {
		return runWSLServiceWithSudo("disable", serviceName, err, out)
	}
	return nil
}
