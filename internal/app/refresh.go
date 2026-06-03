package app

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

func Refresh() (Inventory, error) {
	var warnings []error
	var runners []Runner

	local, err := DiscoverLocal()
	if err != nil {
		warnings = append(warnings, err)
	}
	runners = append(runners, local...)

	statuses, queues, err := LoadGitHubStatus(uniqueRepos(runners))
	if err != nil {
		warnings = append(warnings, err)
	}

	for i := range runners {
		key := runnerKey(runners[i].Repo, runners[i].Name)
		if status, ok := statuses[key]; ok {
			runners[i].GitHubStatus = status.Status
			runners[i].Busy = status.Busy
			runners[i].Labels = status.Labels
			runners[i].Version = status.Version
			if status.OS != "" {
				runners[i].OS = status.OS
			}
		} else if runners[i].GitHubStatus == "" {
			runners[i].GitHubStatus = "unknown"
		}
		if queue, ok := queues[runners[i].Repo]; ok {
			runners[i].QueueCount = queue.Count
			runners[i].StaleQueueCount = queue.StaleCount
		}
	}

	sort.SliceStable(runners, func(i, j int) bool {
		left := strings.ToLower(fmt.Sprintf("%s/%s/%s", runners[i].Host, runners[i].Repo, runners[i].Name))
		right := strings.ToLower(fmt.Sprintf("%s/%s/%s", runners[j].Host, runners[j].Repo, runners[j].Name))
		return left < right
	})

	return Inventory{Runners: runners}, errors.Join(warnings...)
}

func uniqueRepos(runners []Runner) []string {
	seen := map[string]bool{}
	var repos []string
	for _, runner := range runners {
		if runner.Repo == "" || seen[runner.Repo] {
			continue
		}
		seen[runner.Repo] = true
		repos = append(repos, runner.Repo)
	}
	sort.Strings(repos)
	return repos
}
