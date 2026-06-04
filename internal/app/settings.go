package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Settings struct {
	ProjectsRoot              string   `json:"projectsRoot"`
	WindowsRunnerRoots        []string `json:"windowsRunnerRoots"`
	WSLRunnerRoots            []string `json:"wslRunnerRoots"`
	LinuxRunnerRoots          []string `json:"linuxRunnerRoots"`
	TUIRefreshIntervalSeconds int      `json:"tuiRefreshIntervalSeconds"`
	WSLSudoPassword           string   `json:"wslSudoPassword"`
}

func DefaultSettings() Settings {
	return Settings{
		ProjectsRoot:              `C:\Repos`,
		WindowsRunnerRoots:        []string{`C:\Runners`},
		WSLRunnerRoots:            []string{`/home/gsv777/Runners`},
		LinuxRunnerRoots:          []string{"/opt/Runners", "/srv/Runners"},
		TUIRefreshIntervalSeconds: 5,
	}
}

func SettingsPath() (string, error) {
	if path := strings.TrimSpace(os.Getenv("RUNNER_MONITOR_CONFIG")); path != "" {
		return path, nil
	}
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exe), "runner-monitor.json"), nil
}

func LoadSettings() (Settings, error) {
	path, err := SettingsPath()
	if err != nil {
		return DefaultSettings(), err
	}
	return LoadSettingsAt(path)
}

func LoadSettingsAt(path string) (Settings, error) {
	settings := DefaultSettings()
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return settings, nil
	}
	if err != nil {
		return settings, err
	}
	text := strings.TrimPrefix(string(data), "\uFEFF")
	if err := json.Unmarshal([]byte(text), &settings); err != nil {
		return settings, err
	}
	return normalizeSettings(settings), nil
}

func effectiveSettings() Settings {
	settings, err := LoadSettings()
	if err != nil {
		return DefaultSettings()
	}
	return settings
}

func InitSettings(overwrite bool) (string, bool, error) {
	path, err := SettingsPath()
	if err != nil {
		return "", false, err
	}
	if _, err := os.Stat(path); err == nil && !overwrite {
		return path, false, nil
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return path, false, err
	}
	if err := SaveSettingsAt(path, DefaultSettings()); err != nil {
		return path, false, err
	}
	return path, true, nil
}

func SaveSettingsAt(path string, settings Settings) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(normalizeSettings(settings), "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o600)
}

func RenderSettings(settings Settings) string {
	masked := normalizeSettings(settings)
	if strings.TrimSpace(masked.WSLSudoPassword) != "" {
		masked.WSLSudoPassword = "<set>"
	} else {
		masked.WSLSudoPassword = "<empty>"
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(masked); err != nil {
		return fmt.Sprintf("settings render failed: %v\n", err)
	}
	return buf.String()
}

func normalizeSettings(settings Settings) Settings {
	defaults := DefaultSettings()
	if strings.TrimSpace(settings.ProjectsRoot) == "" {
		settings.ProjectsRoot = defaults.ProjectsRoot
	}
	if len(nonEmptyStrings(settings.WindowsRunnerRoots)) == 0 {
		settings.WindowsRunnerRoots = defaults.WindowsRunnerRoots
	} else {
		settings.WindowsRunnerRoots = nonEmptyStrings(settings.WindowsRunnerRoots)
	}
	if len(nonEmptyStrings(settings.WSLRunnerRoots)) == 0 {
		settings.WSLRunnerRoots = defaults.WSLRunnerRoots
	} else {
		settings.WSLRunnerRoots = nonEmptyStrings(settings.WSLRunnerRoots)
	}
	if len(nonEmptyStrings(settings.LinuxRunnerRoots)) == 0 {
		settings.LinuxRunnerRoots = defaults.LinuxRunnerRoots
	} else {
		settings.LinuxRunnerRoots = nonEmptyStrings(settings.LinuxRunnerRoots)
	}
	if settings.TUIRefreshIntervalSeconds <= 0 {
		settings.TUIRefreshIntervalSeconds = defaults.TUIRefreshIntervalSeconds
	}
	return settings
}

func (settings Settings) TUIRefreshInterval() time.Duration {
	normalized := normalizeSettings(settings)
	return time.Duration(normalized.TUIRefreshIntervalSeconds) * time.Second
}

func nonEmptyStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}
