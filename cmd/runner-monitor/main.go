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
	disableAutostart := flag.Bool("disable-autostart", false, "disable boot autostart for service-managed runners without stopping them")
	flag.Parse()

	inventory, err := app.Refresh()
	if err != nil {
		fmt.Fprintf(os.Stderr, "initial refresh warning: %v\n", err)
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

	program := tea.NewProgram(app.NewModel(inventory))
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runner-monitor failed: %v\n", err)
		os.Exit(1)
	}
}
