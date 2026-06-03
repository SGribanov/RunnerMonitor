package app

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type windowsService struct {
	Name     string `json:"Name"`
	State    string `json:"State"`
	PathName string `json:"PathName"`
}

func discoverWindowsRunnerDirs(services map[string]windowsService) ([]Runner, error) {
	roots := []string{
		`C:\Runners`,
		`C:\actions-runner*`,
		`C:\github-runners`,
	}
	files, err := findRunnerFiles(roots, 3)
	if err != nil {
		return nil, err
	}

	var runners []Runner
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		config, err := parseRunnerConfig(data)
		if err != nil {
			continue
		}
		dir := filepath.Dir(file)
		runner := runnerFromConfig(config, dir, "local", "windows")
		if service, ok := services[strings.ToLower(dir)]; ok {
			runner.ServiceName = service.Name
			runner.LocalState = strings.ToLower(service.State)
			runner.ControlMode = "windows-service"
		}
		runners = append(runners, runner)
	}
	return runners, nil
}

func discoverWindowsServices() (map[string]windowsService, error) {
	script := `Get-CimInstance Win32_Service | Where-Object { $_.Name -like 'actions.runner.*' -or $_.DisplayName -like '*GitHub Actions Runner*' } | Select-Object Name,State,PathName | ConvertTo-Json -Depth 3`
	out, err := exec.Command("powershell", "-NoProfile", "-Command", script).Output()
	if err != nil {
		return map[string]windowsService{}, err
	}
	out = []byte(strings.TrimSpace(string(out)))
	if len(out) == 0 {
		return map[string]windowsService{}, nil
	}

	var services []windowsService
	if out[0] == '{' {
		var service windowsService
		if err := json.Unmarshal(out, &service); err != nil {
			return nil, err
		}
		services = []windowsService{service}
	} else if err := json.Unmarshal(out, &services); err != nil {
		return nil, err
	}

	byDir := map[string]windowsService{}
	for _, service := range services {
		dir := runnerDirFromServicePath(service.PathName)
		if dir == "" {
			continue
		}
		byDir[strings.ToLower(dir)] = service
	}
	return byDir, nil
}

func runnerDirFromServicePath(pathName string) string {
	pathName = strings.Trim(pathName, `"`)
	if pathName == "" {
		return ""
	}
	dir := filepath.Dir(filepath.Dir(pathName))
	return filepath.Clean(dir)
}

func runWindowsService(action string, serviceName string) error {
	var verb string
	switch action {
	case "start":
		verb = "Start-Service"
	case "stop":
		verb = "Stop-Service"
	case "restart":
		verb = "Restart-Service"
	default:
		return fmt.Errorf("unsupported action %q", action)
	}
	cmd := exec.Command("powershell", "-NoProfile", "-Command", fmt.Sprintf("%s -Name %q", verb, serviceName))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func disableWindowsServiceAutostart(serviceName string) error {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", fmt.Sprintf("Set-Service -Name %q -StartupType Manual", serviceName))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
