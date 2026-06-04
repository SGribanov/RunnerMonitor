package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/SGribanov/RunnerMonitor/internal/app"
)

func main() {
	once := flag.Bool("once", false, "print runner inventory once and exit")
	audit := flag.Bool("audit", false, "print runner cleanup audit once and exit")
	startRepo := flag.String("start-repo", "", "start service-managed runners for owner/repo")
	stopRepo := flag.String("stop-repo", "", "stop service-managed runners for owner/repo")
	restartRepo := flag.String("restart-repo", "", "restart service-managed runners for owner/repo")
	clearRepo := flag.String("clear-repo", "", "clear idle runner work directories for owner/repo")
	clearRunner := flag.String("clear-runner", "", "clear one runner by name")
	project := flag.String("project", "", "project folder name under C:\\Repos")
	repo := flag.String("repo", "", "GitHub owner/repo filter for runner operations")
	removeRunner := flag.String("remove-runner", "", "remove one runner by name; dry-run unless --confirm is set")
	addRunner := flag.String("add-runner", "", "add/configure one runner by name; dry-run unless --confirm is set")
	runnerFolder := flag.String("runner-folder", "", "existing runner distribution folder for --add-runner")
	labels := flag.String("labels", "", "comma-separated labels for --add-runner")
	confirm := flag.Bool("confirm", false, "execute destructive or provisioning operation")
	force := flag.Bool("force", false, "allow selected destructive operation to proceed past busy checks")
	deleteFolder := flag.Bool("delete-folder", false, "delete runner folder after successful unregister")
	replace := flag.Bool("replace", false, "pass --replace when configuring a runner")
	service := flag.Bool("service", false, "install and start runner service after configuring")
	startCurrent := flag.Bool("start-current", false, "start service-managed runners for the current git origin repository")
	stopCurrent := flag.Bool("stop-current", false, "stop service-managed runners for the current git origin repository")
	restartCurrent := flag.Bool("restart-current", false, "restart service-managed runners for the current git origin repository")
	clearCurrent := flag.Bool("clear-current", false, "clear idle runner work directories for the current git origin repository")
	clearIdle := flag.Bool("clear-idle", false, "clear all idle runner work directories")
	disableAutostart := flag.Bool("disable-autostart", false, "disable boot autostart for service-managed runners without stopping them")
	configureRemote := flag.String("configure-remote", "", "prompt for SSH remote runner host settings and save them")
	connectRemote := flag.String("connect-remote", "", "open the saved SSH remote runner host TUI")
	flag.Parse()

	needsInventory := *once || *audit || *disableAutostart || *startCurrent || *stopCurrent || *restartCurrent || *clearCurrent || *clearIdle ||
		*startRepo != "" || *stopRepo != "" || *restartRepo != "" || *clearRepo != "" || *clearRunner != "" || *removeRunner != ""
	var inventory app.Inventory
	if needsInventory {
		var err error
		inventory, err = app.Refresh()
		if err != nil {
			fmt.Fprintf(os.Stderr, "initial refresh warning: %v\n", err)
		}
	}
	if *once {
		fmt.Print(app.RenderInventory(inventory))
		return
	}
	if *audit {
		fmt.Print(app.RenderAudit(inventory))
		return
	}
	if *disableAutostart {
		fmt.Print(app.DisableAutostart(inventory))
		return
	}
	if *startCurrent || *stopCurrent || *restartCurrent || *clearCurrent {
		repo, err := app.CurrentGitHubRepo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot detect current GitHub repo: %v\n", err)
			os.Exit(1)
		}
		action := "start"
		if *stopCurrent {
			action = "stop"
		}
		if *restartCurrent {
			action = "restart"
		}
		if *clearCurrent {
			fmt.Print(app.ClearRepoRunners(repo, inventory))
		} else {
			fmt.Print(app.RunRepoLifecycle(action, repo, inventory))
		}
		return
	}
	if *startRepo != "" {
		fmt.Print(app.RunRepoLifecycle("start", *startRepo, inventory))
		return
	}
	if *stopRepo != "" {
		fmt.Print(app.RunRepoLifecycle("stop", *stopRepo, inventory))
		return
	}
	if *restartRepo != "" {
		fmt.Print(app.RunRepoLifecycle("restart", *restartRepo, inventory))
		return
	}
	if *clearRepo != "" {
		fmt.Print(app.ClearRepoRunners(*clearRepo, inventory))
		return
	}
	if *clearRunner != "" {
		fmt.Print(app.ClearNamedRunner(*clearRunner, inventory))
		return
	}
	if *removeRunner != "" {
		fmt.Print(app.RemoveNamedRunner(app.RemoveRunnerOptions{
			Name:         *removeRunner,
			Project:      *project,
			Repo:         *repo,
			Confirm:      *confirm,
			Force:        *force,
			DeleteFolder: *deleteFolder,
		}, inventory))
		return
	}
	if *addRunner != "" {
		fmt.Println(app.AddRunner(app.AddRunnerOptions{
			Project:      *project,
			Name:         *addRunner,
			RunnerFolder: *runnerFolder,
			Labels:       *labels,
			Confirm:      *confirm,
			Replace:      *replace,
			Service:      *service,
		}))
		return
	}
	if *clearIdle {
		fmt.Print(app.ClearIdleRunners(inventory))
		return
	}
	if *configureRemote != "" {
		if err := app.ConfigureRemoteHost(*configureRemote, os.Stdin, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "configure remote failed: %v\n", err)
			os.Exit(1)
		}
		return
	}
	if *connectRemote != "" {
		if err := app.ConnectRemoteHost(*connectRemote, os.Stdin, os.Stdout, os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "connect remote failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	program := tea.NewProgram(app.NewLoadingModel())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runner-monitor failed: %v\n", err)
		os.Exit(1)
	}
}
