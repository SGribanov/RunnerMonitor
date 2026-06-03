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
	startCurrent := flag.Bool("start-current", false, "start service-managed runners for the current git origin repository")
	stopCurrent := flag.Bool("stop-current", false, "stop service-managed runners for the current git origin repository")
	restartCurrent := flag.Bool("restart-current", false, "restart service-managed runners for the current git origin repository")
	disableAutostart := flag.Bool("disable-autostart", false, "disable boot autostart for service-managed runners without stopping them")
	configureRemote := flag.String("configure-remote", "", "prompt for SSH remote runner host settings and save them")
	connectRemote := flag.String("connect-remote", "", "open the saved SSH remote runner host TUI")
	flag.Parse()

	needsInventory := *once || *audit || *disableAutostart || *startCurrent || *stopCurrent || *restartCurrent ||
		*startRepo != "" || *stopRepo != "" || *restartRepo != ""
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
	if *startCurrent || *stopCurrent || *restartCurrent {
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
		fmt.Print(app.RunRepoLifecycle(action, repo, inventory))
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
