package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RemoteConfig struct {
	Hosts map[string]RemoteHost `json:"hosts"`
}

type RemoteHost struct {
	Name               string `json:"name"`
	SSHHost            string `json:"sshHost"`
	OS                 string `json:"os"`
	RunnerMonitorPath  string `json:"runnerMonitorPath"`
	DefaultProjectPath string `json:"defaultProjectPath,omitempty"`
}

func RemoteConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "RunnerMonitor", "remote-hosts.json"), nil
}

func ConfigureRemoteHost(name string, in io.Reader, out io.Writer) error {
	path, err := RemoteConfigPath()
	if err != nil {
		return err
	}
	config, err := loadRemoteConfig(path)
	if err != nil {
		return err
	}
	host, err := promptRemoteHost(name, config.Hosts[name], in, out)
	if err != nil {
		return err
	}
	config.Hosts[host.Name] = host
	if err := saveRemoteConfig(path, config); err != nil {
		return err
	}
	fmt.Fprintf(out, "saved remote host %q to %s\n", host.Name, path)
	return nil
}

func ConnectRemoteHost(name string, in io.Reader, out io.Writer, errOut io.Writer) error {
	path, err := RemoteConfigPath()
	if err != nil {
		return err
	}
	config, err := loadRemoteConfig(path)
	if err != nil {
		return err
	}
	host, ok := config.Hosts[name]
	if !ok {
		fmt.Fprintf(out, "remote host %q is not configured yet\n", name)
		host, err = promptRemoteHost(name, RemoteHost{}, in, out)
		if err != nil {
			return err
		}
		config.Hosts[host.Name] = host
		if err := saveRemoteConfig(path, config); err != nil {
			return err
		}
		fmt.Fprintf(out, "saved remote host %q to %s\n", host.Name, path)
	}
	args := remoteTUISSHArgs(host)
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = errOut
	return cmd.Run()
}

func RemoteTUIProcess(name string) (*exec.Cmd, error) {
	path, err := RemoteConfigPath()
	if err != nil {
		return nil, err
	}
	config, err := loadRemoteConfig(path)
	if err != nil {
		return nil, err
	}
	host, ok := config.Hosts[name]
	if !ok {
		return nil, fmt.Errorf("remote host %q is not configured; run runner-monitor --configure-remote %s", name, name)
	}
	args := remoteTUISSHArgs(host)
	return exec.Command("ssh", args...), nil
}

func loadRemoteConfig(path string) (RemoteConfig, error) {
	config := RemoteConfig{Hosts: map[string]RemoteHost{}}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return config, nil
	}
	if err != nil {
		return config, err
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return config, err
	}
	if config.Hosts == nil {
		config.Hosts = map[string]RemoteHost{}
	}
	return config, nil
}

func saveRemoteConfig(path string, config RemoteConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o600)
}

func promptRemoteHost(name string, existing RemoteHost, in io.Reader, out io.Writer) (RemoteHost, error) {
	reader := bufio.NewReader(in)
	defaultName := firstNonEmpty(existing.Name, name, "runnerbox")
	host := RemoteHost{
		Name:               defaultName,
		SSHHost:            firstNonEmpty(existing.SSHHost, defaultName),
		OS:                 firstNonEmpty(existing.OS, "windows"),
		RunnerMonitorPath:  existing.RunnerMonitorPath,
		DefaultProjectPath: existing.DefaultProjectPath,
	}

	var err error
	host.Name, err = promptValue(reader, out, "Remote name", host.Name)
	if err != nil {
		return host, err
	}
	host.SSHHost, err = promptValue(reader, out, "SSH host", firstNonEmpty(host.SSHHost, host.Name))
	if err != nil {
		return host, err
	}
	host.OS, err = promptValue(reader, out, "Host OS (windows/linux)", host.OS)
	if err != nil {
		return host, err
	}
	host.OS = strings.ToLower(strings.TrimSpace(host.OS))
	if host.OS != "windows" && host.OS != "linux" {
		return host, fmt.Errorf("unsupported host OS %q", host.OS)
	}
	host.RunnerMonitorPath, err = promptValue(reader, out, "RunnerMonitor path", firstNonEmpty(host.RunnerMonitorPath, defaultRunnerMonitorPath(host.OS)))
	if err != nil {
		return host, err
	}
	host.DefaultProjectPath, err = promptValue(reader, out, "Default project path", firstNonEmpty(host.DefaultProjectPath, defaultProjectPath(host.OS)))
	if err != nil {
		return host, err
	}
	return host, nil
}

func promptValue(reader *bufio.Reader, out io.Writer, label string, fallback string) (string, error) {
	fmt.Fprintf(out, "%s [%s]: ", label, fallback)
	value, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback, nil
	}
	return value, nil
}

func remoteTUISSHArgs(host RemoteHost) []string {
	return []string{"-t", host.SSHHost, remoteTUICommand(host)}
}

func remoteTUICommand(host RemoteHost) string {
	switch strings.ToLower(host.OS) {
	case "linux":
		return host.RunnerMonitorPath
	default:
		return fmt.Sprintf("powershell -NoProfile -ExecutionPolicy Bypass -File %s", host.RunnerMonitorPath)
	}
}

func defaultRunnerMonitorPath(osName string) string {
	if strings.EqualFold(osName, "linux") {
		return "/opt/RunnerMonitor/runner-monitor"
	}
	return "C:/Repos/RunnerMonitor/runner-monitor.ps1"
}

func defaultProjectPath(osName string) string {
	if strings.EqualFold(osName, "linux") {
		return "/srv/DeltaG"
	}
	return "C:/Repos/DeltaG"
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
