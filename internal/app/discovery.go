package app

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func DiscoverLocal() ([]Runner, error) {
	var warnings []error
	var runners []Runner
	settings := effectiveSettings()

	if runtime.GOOS == "windows" {
		services, err := discoverWindowsServices()
		if err != nil {
			warnings = append(warnings, err)
		}
		processes, err := discoverWindowsRunnerProcesses()
		if err != nil {
			warnings = append(warnings, err)
		}
		windows, err := discoverWindowsRunnerDirs(services, processes, settings.WindowsRunnerRoots)
		if err != nil {
			warnings = append(warnings, err)
		}
		runners = append(runners, windows...)

		wsl, err := discoverWSLRunners(settings.WSLRunnerRoots)
		if err != nil {
			warnings = append(warnings, err)
		}
		runners = append(runners, wsl...)
	}

	if runtime.GOOS != "windows" {
		linux, err := discoverLinuxRunnerDirs(settings.LinuxRunnerRoots)
		if err != nil {
			warnings = append(warnings, err)
		}
		runners = append(runners, linux...)
	}

	return dedupeRunners(runners), errors.Join(warnings...)
}

func parseRunnerConfig(data []byte) (runnerConfig, error) {
	text := strings.TrimPrefix(string(data), "\uFEFF")
	var config runnerConfig
	err := json.Unmarshal([]byte(text), &config)
	return config, err
}

func runnerFromConfig(config runnerConfig, path, host, transport string) Runner {
	return Runner{
		Name:        config.AgentName,
		Repo:        repoFromGitHubURL(config.GitHubURL),
		Path:        path,
		Host:        host,
		Transport:   transport,
		LocalState:  "manual",
		ControlMode: "manual",
	}
}

func findRunnerFiles(roots []string, maxDepth int) ([]string, error) {
	var files []string
	for _, root := range roots {
		if root == "" {
			continue
		}
		matches, _ := filepath.Glob(root)
		if len(matches) == 0 {
			matches = []string{root}
		}
		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil || !info.IsDir() {
				continue
			}
			baseDepth := strings.Count(filepath.Clean(match), string(os.PathSeparator))
			err = filepath.WalkDir(match, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if d.IsDir() {
					depth := strings.Count(filepath.Clean(path), string(os.PathSeparator)) - baseDepth
					if depth > maxDepth {
						return filepath.SkipDir
					}
					if strings.EqualFold(d.Name(), "_work") || strings.EqualFold(d.Name(), "_diag") {
						return filepath.SkipDir
					}
					return nil
				}
				if d.Name() == ".runner" {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				return files, err
			}
		}
	}
	return files, nil
}

func dedupeRunners(runners []Runner) []Runner {
	seen := map[string]bool{}
	var result []Runner
	for _, runner := range runners {
		key := strings.ToLower(runner.Host + "|" + runner.Path + "|" + runner.Name)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, runner)
	}
	return result
}
