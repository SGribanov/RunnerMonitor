package app

import (
	"os"
	"path/filepath"
	"strings"
)

func discoverLinuxRunnerDirs(roots []string) ([]Runner, error) {
	files, err := findRunnerFiles(roots, 4)
	if err != nil {
		return nil, err
	}

	var runners []Runner
	host, _ := os.Hostname()
	if host == "" {
		host = "local"
	}
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
		runner := runnerFromConfig(config, dir, host, "linux")
		if service := readTextFile(filepath.Join(dir, ".service")); service != "" {
			runner.ServiceName = strings.TrimSpace(service)
			runner.ControlMode = "systemd"
			runner.LocalState = "configured"
		}
		runners = append(runners, runner)
	}
	return runners, nil
}

func readTextFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
