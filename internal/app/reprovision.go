package app

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

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
	if runner.IsReadOnlyGitHubRow() {
		return fmt.Sprintf("%s is %s and read-only; removal skipped", runner.Name, runnerReadOnlyKind(runner))
	}
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
	return ProjectRepoFromFolderAt(effectiveSettings().ProjectsRoot, project)
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
		return runLocalRunnerConfig(runner.Path, "remove", []string{"--unattended"}, token)
	case "wsl":
		return runWSLRunnerConfig(runner.Path, "remove", []string{"--unattended"}, token)
	default:
		return runLocalRunnerConfig(runner.Path, "remove", []string{"--unattended"}, token)
	}
}

func configureRunner(repo string, token string, options AddRunnerOptions) error {
	args := []string{
		"--unattended",
		"--url", "https://github.com/" + repo,
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
		return runWSLRunnerConfig(options.RunnerFolder, "configure", args, token)
	}
	return runLocalRunnerConfig(options.RunnerFolder, "configure", args, token)
}

func runLocalRunnerConfig(folder string, action string, args []string, token string) error {
	scriptName := "config.sh"
	if runtime.GOOS == "windows" || fileExists(filepath.Join(folder, "config.cmd")) {
		scriptName = "config.cmd"
	}
	command := filepath.Join(folder, scriptName)
	runArgs := append(configActionArgs(action), args...)
	if token != "" {
		runArgs = append(runArgs, "--token", token)
	}
	var cmd *exec.Cmd
	if strings.EqualFold(scriptName, "config.cmd") {
		if token != "" {
			cmd = exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", configPowerShellScriptWithToken(command, runArgs))
			cmd.Env = envWithRunnerToken(token, false)
		} else {
			cmd = exec.Command("cmd.exe", append([]string{"/c", command}, runArgs...)...)
		}
	} else if token != "" {
		cmd = exec.Command("sh", "-lc", configShellLineWithToken(command, runArgs))
		cmd.Env = append(os.Environ(), "RUNNER_MONITOR_RUNNER_TOKEN="+token)
	} else {
		cmd = exec.Command(command, runArgs...)
	}
	cmd.Dir = folder
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, maskSecret(strings.TrimSpace(string(out)), token))
	}
	return nil
}

func runWSLRunnerConfig(folder string, action string, args []string, token string) error {
	command := "./config.sh"
	runArgs := append(configActionArgs(action), args...)
	quoted := []string{"cd", shellQuote(folder), "&&", shellQuote(command)}
	for _, arg := range runArgs {
		quoted = append(quoted, shellQuote(arg))
	}
	if token != "" {
		quoted = append(quoted, "--token", `"$RUNNER_MONITOR_RUNNER_TOKEN"`)
	}
	cmd := exec.Command("wsl.exe", "bash", "-lc", strings.Join(quoted, " "))
	if token != "" {
		cmd.Env = envWithRunnerToken(token, true)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, maskSecret(strings.TrimSpace(string(out)), token))
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

func configPowerShellScriptWithToken(command string, args []string) string {
	parts := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		if args[i] == "--token" && i+1 < len(args) {
			parts = append(parts, powerShellQuote("--token"), "$env:RUNNER_MONITOR_RUNNER_TOKEN")
			i++
			continue
		}
		parts = append(parts, powerShellQuote(args[i]))
	}
	return fmt.Sprintf("$argsList = @(%s); & %s @argsList", strings.Join(parts, ","), powerShellQuote(command))
}

func configShellLineWithToken(command string, args []string) string {
	parts := []string{shellQuote(command)}
	for i := 0; i < len(args); i++ {
		if args[i] == "--token" && i+1 < len(args) {
			parts = append(parts, "--token", `"$RUNNER_MONITOR_RUNNER_TOKEN"`)
			i++
			continue
		}
		parts = append(parts, shellQuote(args[i]))
	}
	return strings.Join(parts, " ")
}

func maskSecret(text string, secret string) string {
	if secret == "" {
		return text
	}
	return strings.ReplaceAll(text, secret, "<redacted>")
}

func envWithRunnerToken(token string, passToWSL bool) []string {
	env := append(os.Environ(), "RUNNER_MONITOR_RUNNER_TOKEN="+token)
	if passToWSL {
		env = append(env, "WSLENV="+appendWSLEnv(os.Getenv("WSLENV"), "RUNNER_MONITOR_RUNNER_TOKEN"))
	}
	return env
}

func appendWSLEnv(existing string, name string) string {
	for _, part := range strings.Split(existing, ":") {
		if part == name {
			return existing
		}
	}
	if strings.TrimSpace(existing) == "" {
		return name
	}
	return existing + ":" + name
}

func runWSLShellWithSudo(command string, originalErr error, originalOut []byte) error {
	password, passwordErr := wslSudoPassword()
	if passwordErr != nil {
		return fmt.Errorf("%w: %s; sudo fallback failed: %v", originalErr, strings.TrimSpace(string(originalOut)), passwordErr)
	}
	cmd := exec.Command("wsl.exe", "--", "sudo", "-S", "-p", "", "bash", "-lc", command)
	cmd.Stdin = strings.NewReader(password)
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
	settings := effectiveSettings()
	if runner.Transport == "wsl" || isWSLPath(runner.Path) {
		return isUnderAnySlashRoot(runner.Path, append(settings.WSLRunnerRoots, settings.LinuxRunnerRoots...))
	}
	return isUnderAnyWindowsRoot(pathValue, settings.WindowsRunnerRoots)
}

func isUnderAnyWindowsRoot(pathValue string, roots []string) bool {
	pathValue = strings.ToLower(filepath.Clean(pathValue))
	for _, root := range roots {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}
		root = strings.ToLower(filepath.Clean(root))
		if pathValue != root && strings.HasPrefix(pathValue, root+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

func isUnderAnySlashRoot(pathValue string, roots []string) bool {
	pathValue = pathCleanSlash(pathValue)
	for _, root := range roots {
		root = pathCleanSlash(root)
		if root == "" || root == "." || root == "/" {
			continue
		}
		if pathValue != root && strings.HasPrefix(pathValue, strings.TrimRight(root, "/")+"/") {
			return true
		}
	}
	return false
}

func pathCleanSlash(pathValue string) string {
	pathValue = strings.ReplaceAll(strings.TrimSpace(pathValue), `\`, "/")
	if pathValue == "" {
		return ""
	}
	return path.Clean(pathValue)
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
