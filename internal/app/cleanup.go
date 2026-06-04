package app

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func ClearRunner(runner Runner) string {
	if runner.Busy {
		return fmt.Sprintf("%s is busy; cleanup skipped", runner.Name)
	}
	if runner.Path == "" || runner.Path == "(unit only)" {
		return fmt.Sprintf("%s path is unknown; cleanup skipped", runner.Name)
	}

	wasRunning := shouldStopForCleanup(runner)
	if wasRunning {
		if err := stopRunnerForCleanup(runner); err != nil {
			return fmt.Sprintf("clear %s failed while stopping: %v", runner.Name, err)
		}
	}

	if err := clearRunnerFolder(runner); err != nil {
		if wasRunning {
			_ = startRunnerForCleanup(runner)
		}
		return fmt.Sprintf("clear %s failed: %v", runner.Name, err)
	}

	if wasRunning {
		if err := startRunnerForCleanup(runner); err != nil {
			return fmt.Sprintf("cleared %s, but restart failed: %v", runner.Name, err)
		}
	}
	return fmt.Sprintf("cleared %s", runner.Name)
}

func shouldStopForCleanup(runner Runner) bool {
	if isAlreadyRunning(runner.LocalState) {
		return true
	}
	return runner.ControlMode == "manual" &&
		runner.Transport == "windows" &&
		strings.EqualFold(runner.GitHubStatus, "online")
}

func ClearIdleRunners(inventory Inventory) string {
	var b strings.Builder
	count := 0
	for _, runner := range inventory.Runners {
		if runner.Busy {
			continue
		}
		count++
		fmt.Fprintf(&b, "%s\n", ClearRunner(runner))
	}
	if count == 0 {
		return "no idle runners found\n"
	}
	return b.String()
}

func ClearRepoRunners(repo string, inventory Inventory) string {
	var b strings.Builder
	count := 0
	for _, runner := range inventory.Runners {
		if !strings.EqualFold(runner.Repo, repo) {
			continue
		}
		count++
		fmt.Fprintf(&b, "%s\n", ClearRunner(runner))
	}
	if count == 0 {
		return fmt.Sprintf("no runners found for %s\n", repo)
	}
	return b.String()
}

func stopRunnerForCleanup(runner Runner) error {
	return controlRunnerForCleanup("stop", runner)
}

func startRunnerForCleanup(runner Runner) error {
	return controlRunnerForCleanup("start", runner)
}

func controlRunnerForCleanup(action string, runner Runner) error {
	if runner.ControlMode == "manual" && runner.Transport == "windows" {
		if runtime.GOOS != "windows" {
			return fmt.Errorf("manual Windows runner control is only available on Windows")
		}
		return runWindowsManualRunner(action, runner.Path)
	}
	if runner.ServiceName == "" {
		return fmt.Errorf("%s is not controllable", runner.Name)
	}
	switch runner.ControlMode {
	case "windows-service":
		if runtime.GOOS != "windows" {
			return fmt.Errorf("Windows service control is only available on Windows")
		}
		return runWindowsService(action, runner.ServiceName)
	case "wsl-systemd":
		return runWSLService(action, runner.ServiceName)
	case "systemd":
		return runCommandWithOutput("systemctl", action, runner.ServiceName)
	default:
		return fmt.Errorf("unsupported control mode %q", runner.ControlMode)
	}
}

func clearRunnerFolder(runner Runner) error {
	if runner.Transport == "wsl" {
		return clearWSLRunnerFolder(runner.Path)
	}
	return clearLocalRunnerFolder(runner.Path)
}

func clearLocalRunnerFolder(root string) error {
	work := filepath.Join(root, "_work")
	if err := clearDirectoryContents(work); err != nil {
		return err
	}
	patterns := []string{"actions-runner*.zip", "actions-runner*.tar.gz"}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(root, pattern))
		if err != nil {
			return err
		}
		for _, match := range matches {
			if err := os.Remove(match); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}

func clearDirectoryContents(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := os.RemoveAll(filepath.Join(dir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func clearWSLRunnerFolder(root string) error {
	encodedRoot := base64.StdEncoding.EncodeToString([]byte(root))
	script := `
import base64
import pathlib
import shutil
import sys

root = pathlib.Path(base64.b64decode(sys.argv[1]).decode("utf-8"))
if not str(root):
    raise SystemExit("runner root argument is empty")
work = root / "_work"
work.mkdir(parents=True, exist_ok=True)
for child in work.iterdir():
    if child.is_dir() and not child.is_symlink():
        shutil.rmtree(child)
    else:
        child.unlink(missing_ok=True)
for pattern in ("actions-runner*.zip", "actions-runner*.tar.gz"):
    for child in root.glob(pattern):
        child.unlink(missing_ok=True)
`
	cmd := exec.Command("wsl.exe", "--", "python3", "-c", script, encodedRoot)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
