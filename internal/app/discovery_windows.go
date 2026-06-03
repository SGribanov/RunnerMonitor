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

type windowsRunnerProcess struct {
	ProcessID      int    `json:"ProcessId"`
	Name           string `json:"Name"`
	ExecutablePath string `json:"ExecutablePath"`
	CommandLine    string `json:"CommandLine"`
}

func discoverWindowsRunnerDirs(services map[string]windowsService, processes map[string]windowsRunnerProcess) ([]Runner, error) {
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
		} else if _, ok := processes[strings.ToLower(dir)]; ok {
			runner.LocalState = "running"
			runner.ControlMode = "manual"
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

func discoverWindowsRunnerProcesses() (map[string]windowsRunnerProcess, error) {
	script := `Get-CimInstance Win32_Process | Where-Object { $_.Name -eq 'Runner.Listener.exe' -and $_.ExecutablePath } | Select-Object ProcessId,Name,ExecutablePath,CommandLine | ConvertTo-Json -Depth 3`
	out, err := exec.Command("powershell", "-NoProfile", "-Command", script).Output()
	if err != nil {
		return map[string]windowsRunnerProcess{}, err
	}
	out = []byte(strings.TrimSpace(string(out)))
	if len(out) == 0 {
		return map[string]windowsRunnerProcess{}, nil
	}

	var processes []windowsRunnerProcess
	if out[0] == '{' {
		var process windowsRunnerProcess
		if err := json.Unmarshal(out, &process); err != nil {
			return nil, err
		}
		processes = []windowsRunnerProcess{process}
	} else if err := json.Unmarshal(out, &processes); err != nil {
		return nil, err
	}

	byDir := map[string]windowsRunnerProcess{}
	for _, process := range processes {
		dir := runnerDirFromProcessPath(process.ExecutablePath)
		if dir == "" {
			continue
		}
		byDir[strings.ToLower(dir)] = process
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

func runnerDirFromProcessPath(pathName string) string {
	pathName = strings.Trim(pathName, `"`)
	if pathName == "" {
		return ""
	}
	return filepath.Clean(filepath.Dir(filepath.Dir(pathName)))
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

func runWindowsManualRunner(action string, runnerPath string) error {
	switch action {
	case "start":
		return startWindowsManualRunner(runnerPath)
	case "stop":
		return stopWindowsManualRunner(runnerPath)
	case "restart":
		if err := stopWindowsManualRunner(runnerPath); err != nil {
			return err
		}
		return startWindowsManualRunner(runnerPath)
	default:
		return fmt.Errorf("unsupported action %q", action)
	}
}

func startWindowsManualRunner(runnerPath string) error {
	script := `
$RunnerPath = $env:RUNNER_MONITOR_RUNNER_PATH
$run = Join-Path $RunnerPath 'run.cmd'
if (!(Test-Path -LiteralPath $run)) {
    throw "run.cmd not found at $run"
}
Start-Process -FilePath $run -WorkingDirectory $RunnerPath -WindowStyle Hidden
`
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	cmd.Env = append(os.Environ(), "RUNNER_MONITOR_RUNNER_PATH="+runnerPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func stopWindowsManualRunner(runnerPath string) error {
	script := `
$RunnerPath = $env:RUNNER_MONITOR_RUNNER_PATH
$resolved = (Resolve-Path -LiteralPath $RunnerPath).Path.TrimEnd('\')
$prefix = $resolved + '\'
$names = @('Runner.Listener.exe', 'Runner.Worker.exe', 'Runner.PluginHost.exe')
$procs = Get-CimInstance Win32_Process | Where-Object {
    $_.ExecutablePath -and
    $_.ExecutablePath.StartsWith($prefix, [System.StringComparison]::OrdinalIgnoreCase) -and
    $names -contains $_.Name
}
foreach ($proc in $procs) {
    Stop-Process -Id $proc.ProcessId -Force
}
`
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	cmd.Env = append(os.Environ(), "RUNNER_MONITOR_RUNNER_PATH="+runnerPath)
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
