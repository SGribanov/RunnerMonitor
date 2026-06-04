package app

import (
	"errors"
	"sort"
	"strings"
	"time"
)

const autoRefreshGitHubCacheTTL = 30 * time.Second

func Refresh() (Inventory, error) {
	return refreshWithGitHubStatus(LoadGitHubStatus)
}

func RefreshWithGitHubCache(maxAge time.Duration) (Inventory, error) {
	return refreshWithGitHubStatus(func(repos []string) (map[string]GitHubRunnerStatus, map[string]QueueStatus, error) {
		return LoadGitHubStatusCached(repos, maxAge)
	})
}

func refreshWithGitHubStatus(loadGitHubStatus func([]string) (map[string]GitHubRunnerStatus, map[string]QueueStatus, error)) (Inventory, error) {
	var warnings []error
	var runners []Runner

	local, err := DiscoverLocal()
	if err != nil {
		warnings = append(warnings, err)
	}
	runners = append(runners, local...)

	statuses, queues, err := loadGitHubStatus(uniqueRepos(runners))
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

	sortRunners(runners)

	return Inventory{Runners: runners}, errors.Join(warnings...)
}

func sortRunners(runners []Runner) {
	sortable := make([]struct {
		runner Runner
		key    string
	}, len(runners))
	for i, runner := range runners {
		sortable[i] = struct {
			runner Runner
			key    string
		}{
			runner: runner,
			key:    strings.ToLower(runner.Host + "/" + runner.Repo + "/" + runner.Name),
		}
	}
	sort.SliceStable(sortable, func(i, j int) bool {
		return sortable[i].key < sortable[j].key
	})
	for i, item := range sortable {
		runners[i] = item.runner
	}
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
