package app

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var defaultReposRoot = `C:\Repos`

type RemoveRunnerOptions struct {
	Name         string
	Project      string
	Repo         string
	Confirm      bool
	Force        bool
	DeleteFolder bool
}

type AddRunnerOptions struct {
	Project      string
	Name         string
	RunnerFolder string
	Labels       string
	Confirm      bool
	Replace      bool
	Service      bool
}

type tokenResponse struct {
	Token string `json:"token"`
}

func RemoveNamedRunner(options RemoveRunnerOptions, inventory Inventory) string {
	repo, err := repoFromProjectOption(options.Project, options.Repo)
	if err != nil {
		return fmt.Sprintf("remove %s failed: %v\n", options.Name, err)
	}
	for _, runner := range inventory.Runners {
		if !strings.EqualFold(runner.Name, options.Name) {
			continue
		}
		if repo != "" && !strings.EqualFold(runner.Repo, repo) {
			continue
		}
		return RemoveRunner(runner, options) + "\n"
	}
	if repo != "" {
		return fmt.Sprintf("no runner found named %s for %s\n", options.Name, repo)
	}
	return fmt.Sprintf("no runner found named %s\n", options.Name)
}

func RemoveRunner(runner Runner, options RemoveRunnerOptions) string {
	if runner.Busy && !options.Force {
		return fmt.Sprintf("%s is busy; removal skipped", runner.Name)
	}
	if runner.Path == "" || runner.Path == "(unit only)" {
		return fmt.Sprintf("%s path is unknown; removal skipped", runner.Name)
	}
	if !options.Confirm {
		return renderRemovePlan(runner, options)
	}
	if options.DeleteFolder && !isSafeRunnerRoot(runner) {
		return fmt.Sprintf("%s folder is outside known runner roots; delete refused: %s", runner.Name, runner.Path)
	}
	if needsElevatedWindowsRemoval(runner) && !isElevatedWindowsProcess() {
		if err := launchElevatedRemoveRunner(runner, options); err != nil {
			return fmt.Sprintf("remove %s requires elevated PowerShell, but UAC launch failed: %v", runner.Name, err)
		}
		return fmt.Sprintf("remove %s requested in elevated PowerShell", runner.Name)
	}
	if shouldStopForCleanup(runner) {
		if err := stopRunnerForCleanup(runner); err != nil {
			return fmt.Sprintf("remove %s failed while stopping: %v", runner.Name, err)
		}
	}
	if err := uninstallRunnerService(runner); err != nil {
		return fmt.Sprintf("remove %s failed while uninstalling service: %v", runner.Name, err)
	}
	if err := unregisterRunner(runner); err != nil {
		return fmt.Sprintf("remove %s failed while unregistering: %v", runner.Name, err)
	}
	if options.DeleteFolder {
		if err := deleteRunnerFolder(runner); err != nil {
			return fmt.Sprintf("remove %s unregistered, but folder delete failed: %v", runner.Name, err)
		}
		return fmt.Sprintf("removed %s and deleted folder", runner.Name)
	}
	return fmt.Sprintf("removed %s; folder preserved at %s", runner.Name, runner.Path)
}

func AddRunner(options AddRunnerOptions) string {
	if strings.TrimSpace(options.Project) == "" {
		return "add runner failed: project is required"
	}
	if strings.TrimSpace(options.Name) == "" {
		return "add runner failed: runner name is required"
	}
	if strings.TrimSpace(options.RunnerFolder) == "" {
		return "add runner failed: runner folder is required"
	}
	repo, err := ProjectRepoFromFolder(options.Project)
	if err != nil {
		return fmt.Sprintf("add runner failed: %v", err)
	}
	if err := validateRunnerConfigScript(options.RunnerFolder); err != nil {
		return fmt.Sprintf("add %s failed: %v", options.Name, err)
	}
	if isRunnerConfigured(options.RunnerFolder) && !options.Replace {
		return fmt.Sprintf("add %s failed: runner folder is already configured; use --replace to reconfigure it", options.Name)
	}
	if !options.Confirm {
		return renderAddPlan(repo, options)
	}
	token, err := registrationToken(repo)
	if err != nil {
		return fmt.Sprintf("add %s failed while getting registration token: %v", options.Name, err)
	}
	if err := configureRunner(repo, token, options); err != nil {
		return fmt.Sprintf("add %s failed while configuring runner: %v", options.Name, err)
	}
	if options.Service {
		if err := installAndStartRunnerService(options.RunnerFolder); err != nil {
			return fmt.Sprintf("added %s, but service install/start failed: %v", options.Name, err)
		}
	}
	return fmt.Sprintf("added %s for %s", options.Name, repo)
}

func renderRemovePlan(runner Runner, options RemoveRunnerOptions) string {
	var b strings.Builder
	fmt.Fprintf(&b, "dry-run remove %s for %s\n", runner.Name, runner.Repo)
	fmt.Fprintf(&b, "- path: %s\n", runner.Path)
	fmt.Fprintf(&b, "- service: %s\n", emptyAsDash(runner.ServiceName))
	fmt.Fprintf(&b, "- stop runner if running/active\n")
	if runner.ServiceName != "" {
		fmt.Fprintf(&b, "- uninstall service/unit\n")
	}
	fmt.Fprintf(&b, "- unregister with GitHub remove token\n")
	if options.DeleteFolder {
		fmt.Fprintf(&b, "- delete runner folder after unregistering\n")
	} else {
		fmt.Fprintf(&b, "- preserve runner folder\n")
	}
	fmt.Fprintf(&b, "run with --confirm to execute")
	return strings.TrimSpace(b.String())
}

func renderAddPlan(repo string, options AddRunnerOptions) string {
	var b strings.Builder
	fmt.Fprintf(&b, "dry-run add %s for %s\n", options.Name, repo)
	fmt.Fprintf(&b, "- project: %s\n", options.Project)
	fmt.Fprintf(&b, "- folder: %s\n", options.RunnerFolder)
	fmt.Fprintf(&b, "- labels: %s\n", emptyAsDash(options.Labels))
	fmt.Fprintf(&b, "- configure existing runner distribution folder with GitHub registration token\n")
	if options.Service {
		fmt.Fprintf(&b, "- install and start runner service\n")
	}
	fmt.Fprintf(&b, "run with --confirm to execute")
	return strings.TrimSpace(b.String())
}

func repoFromProjectOption(project string, repo string) (string, error) {
	if strings.TrimSpace(repo) != "" {
		return repoFromGitHubURL(repo), nil
	}
	if strings.TrimSpace(project) == "" {
		return "", nil
	}
	return ProjectRepoFromFolder(project)
}

func launchElevatedRemoveRunner(runner Runner, options RemoveRunnerOptions) error {
	target, args := elevatedRemoveRunnerTarget(runner, options)
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

func elevatedRemoveRunnerTarget(runner Runner, options RemoveRunnerOptions) (string, []string) {
	args := []string{
		"--remove-runner",
		runner.Name,
		"--repo",
		runner.Repo,
		"--confirm",
	}
	if options.Force {
		args = append(args, "--force")
	}
	if options.DeleteFolder {
		args = append(args, "--delete-folder")
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

func ProjectRepoFromFolder(project string) (string, error) {
	return ProjectRepoFromFolderAt(defaultReposRoot, project)
}

func ProjectRepoFromFolderAt(root string, project string) (string, error) {
	project = strings.TrimSpace(project)
	if project == "" {
		return "", fmt.Errorf("project folder name is empty")
	}
	if project == "." || project == ".." || filepath.Clean(project) != project || strings.ContainsAny(project, `\/:`) {
		return "", fmt.Errorf("project must be a folder name, got %q", project)
	}
	projectDir := filepath.Join(root, project)
	if info, err := os.Stat(projectDir); err != nil || !info.IsDir() {
		return "", fmt.Errorf("project folder not found: %s", projectDir)
	}
	out, err := exec.Command("git", "-C", projectDir, "remote", "get-url", "origin").Output()
	if err != nil {
		return "", fmt.Errorf("cannot read git origin in %s: %w", projectDir, err)
	}
	repo := repoFromGitHubURL(strings.TrimSpace(string(out)))
	if repo == "" || !strings.Contains(repo, "/") {
		return "", fmt.Errorf("origin is not a GitHub owner/repo remote in %s", projectDir)
	}
	return repo, nil
}

func registrationToken(repo string) (string, error) {
	return runnerToken(repo, "registration-token")
}

func removeToken(repo string) (string, error) {
	return runnerToken(repo, "remove-token")
}

func runnerToken(repo string, kind string) (string, error) {
	data, err := ghAPIMethod("POST", fmt.Sprintf("repos/%s/actions/runners/%s", repo, kind))
	if err != nil {
		return "", err
	}
	var response tokenResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return "", err
	}
	if response.Token == "" {
		return "", fmt.Errorf("GitHub returned an empty %s", kind)
	}
	return response.Token, nil
}

func unregisterRunner(runner Runner) error {
	token, err := removeToken(runner.Repo)
	if err != nil {
		return err
	}
	switch runner.Transport {
	case "windows", "local":
		return runLocalRunnerConfig(runner.Path, "remove", []string{"--unattended", "--token", token})
	case "wsl":
		return runWSLRunnerConfig(runner.Path, "remove", []string{"--unattended", "--token", token})
	default:
		return runLocalRunnerConfig(runner.Path, "remove", []string{"--unattended", "--token", token})
	}
}

func configureRunner(repo string, token string, options AddRunnerOptions) error {
	args := []string{
		"--unattended",
		"--url", "https://github.com/" + repo,
		"--token", token,
		"--name", options.Name,
		"--work", "_work",
	}
	if options.Labels != "" {
		args = append(args, "--labels", options.Labels)
	}
	if options.Replace {
		args = append(args, "--replace")
	}
	if isWSLPath(options.RunnerFolder) {
		return runWSLRunnerConfig(options.RunnerFolder, "configure", args)
	}
	return runLocalRunnerConfig(options.RunnerFolder, "configure", args)
}

func runLocalRunnerConfig(folder string, action string, args []string) error {
	scriptName := "config.sh"
	if runtime.GOOS == "windows" || fileExists(filepath.Join(folder, "config.cmd")) {
		scriptName = "config.cmd"
	}
	command := filepath.Join(folder, scriptName)
	runArgs := append(configActionArgs(action), args...)
	var cmd *exec.Cmd
	if strings.EqualFold(scriptName, "config.cmd") {
		cmd = exec.Command("cmd.exe", append([]string{"/c", command}, runArgs...)...)
	} else {
		cmd = exec.Command(command, runArgs...)
	}
	cmd.Dir = folder
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func runWSLRunnerConfig(folder string, action string, args []string) error {
	command := "./config.sh"
	runArgs := append(configActionArgs(action), args...)
	quoted := []string{"cd", shellQuote(folder), "&&", shellQuote(command)}
	for _, arg := range runArgs {
		quoted = append(quoted, shellQuote(arg))
	}
	cmd := exec.Command("wsl.exe", "bash", "-lc", strings.Join(quoted, " "))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func configActionArgs(action string) []string {
	if action == "remove" {
		return []string{"remove"}
	}
	return []string{}
}

func uninstallRunnerService(runner Runner) error {
	if runner.ServiceName == "" {
		return nil
	}
	switch runner.ControlMode {
	case "windows-service":
		return runWindowsRunnerSvc(runner.Path, "uninstall")
	case "wsl-systemd":
		return runWSLRunnerSvc(runner.Path, "uninstall")
	case "systemd":
		return runCommandWithOutput(filepath.Join(runner.Path, "svc.sh"), "uninstall")
	default:
		return fmt.Errorf("unsupported control mode %q", runner.ControlMode)
	}
}

func installAndStartRunnerService(folder string) error {
	if isWSLPath(folder) {
		if err := runWSLRunnerSvc(folder, "install"); err != nil {
			return err
		}
		return runWSLRunnerSvc(folder, "start")
	}
	if err := runWindowsRunnerSvc(folder, "install"); err != nil {
		return err
	}
	return runWindowsRunnerSvc(folder, "start")
}

func runWindowsRunnerSvc(folder string, action string) error {
	cmd := exec.Command(filepath.Join(folder, "svc.cmd"), action)
	cmd.Dir = folder
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: service control may require elevated PowerShell: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func runWSLRunnerSvc(folder string, action string) error {
	command := strings.Join([]string{"cd", shellQuote(folder), "&&", "./svc.sh", shellQuote(action)}, " ")
	cmd := exec.Command("wsl.exe", "bash", "-lc", command)
	if out, err := cmd.CombinedOutput(); err != nil {
		return runWSLShellWithSudo(command, err, out)
	}
	return nil
}

func runWSLShellWithSudo(command string, originalErr error, originalOut []byte) error {
	passwordFile := wslSudoPasswordFileWindowsPath()
	password, readErr := os.ReadFile(passwordFile)
	if readErr != nil {
		return fmt.Errorf("%w: %s; sudo fallback failed: cannot read sudo password file %s: %v", originalErr, strings.TrimSpace(string(originalOut)), passwordFile, readErr)
	}
	cmd := exec.Command("wsl.exe", "--", "sudo", "-S", "-p", "", "bash", "-lc", command)
	cmd.Stdin = strings.NewReader(string(password))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s; sudo fallback failed: %v: %s", originalErr, strings.TrimSpace(string(originalOut)), err, strings.TrimSpace(string(out)))
	}
	return nil
}

func validateRunnerConfigScript(folder string) error {
	if isWSLPath(folder) {
		out, err := exec.Command("wsl.exe", "test", "-x", filepath.ToSlash(filepath.Join(folder, "config.sh"))).CombinedOutput()
		if err != nil {
			return fmt.Errorf("config.sh not found or not executable in %s: %s", folder, strings.TrimSpace(string(out)))
		}
		return nil
	}
	if fileExists(filepath.Join(folder, "config.cmd")) || fileExists(filepath.Join(folder, "config.sh")) {
		return nil
	}
	return fmt.Errorf("config.cmd/config.sh not found in %s", folder)
}

func deleteRunnerFolder(runner Runner) error {
	if runner.Transport == "wsl" || isWSLPath(runner.Path) {
		command := "rm -rf -- " + shellQuote(runner.Path)
		cmd := exec.Command("wsl.exe", "bash", "-lc", command)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
		}
		return nil
	}
	return os.RemoveAll(runner.Path)
}

func isSafeRunnerRoot(runner Runner) bool {
	pathValue := filepath.Clean(runner.Path)
	if runner.Transport == "wsl" || isWSLPath(runner.Path) {
		return strings.HasPrefix(runner.Path, "/home/gsv777/Runners/") ||
			strings.HasPrefix(runner.Path, "/opt/Runners/") ||
			strings.HasPrefix(runner.Path, "/srv/Runners/")
	}
	safeRoot := filepath.Clean(`C:\Runners`) + string(os.PathSeparator)
	return strings.HasPrefix(strings.ToLower(pathValue), strings.ToLower(safeRoot))
}

func needsElevatedWindowsRemoval(runner Runner) bool {
	return runtime.GOOS == "windows" && runner.ControlMode == "windows-service"
}

func isWSLPath(pathValue string) bool {
	return strings.HasPrefix(pathValue, "/")
}

func fileExists(pathValue string) bool {
	_, err := os.Stat(pathValue)
	return err == nil
}

func isRunnerConfigured(folder string) bool {
	if isWSLPath(folder) {
		cmd := exec.Command("wsl.exe", "test", "-f", filepath.ToSlash(filepath.Join(folder, ".runner")))
		return cmd.Run() == nil
	}
	return fileExists(filepath.Join(folder, ".runner"))
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func emptyAsDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}
