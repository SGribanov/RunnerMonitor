package app

import (
	"errors"
	"sort"
	"strings"
	"time"
)

const autoRefreshGitHubCacheTTL = 30 * time.Second

var discoverLocal = DiscoverLocal

func Refresh() (Inventory, error) {
	return refreshWithGitHubData(LoadGitHubStatus, LoadGitHubHostedJobs)
}

func RefreshWithGitHubCache(maxAge time.Duration) (Inventory, error) {
	return refreshWithGitHubData(
		func(repos []string) (map[string]GitHubRunnerStatus, map[string]QueueStatus, error) {
			return LoadGitHubStatusCached(repos, maxAge)
		},
		func(repos []string) ([]Runner, error) {
			return LoadGitHubHostedJobsCached(repos, maxAge)
		},
	)
}

func refreshWithGitHubStatus(loadGitHubStatus func([]string) (map[string]GitHubRunnerStatus, map[string]QueueStatus, error)) (Inventory, error) {
	return refreshWithGitHubData(loadGitHubStatus, func([]string) ([]Runner, error) { return nil, nil })
}

func refreshWithGitHubData(
	loadGitHubStatus func([]string) (map[string]GitHubRunnerStatus, map[string]QueueStatus, error),
	loadHostedJobs func([]string) ([]Runner, error),
) (Inventory, error) {
	var warnings []error
	var runners []Runner

	local, err := discoverLocal()
	if err != nil {
		warnings = append(warnings, err)
	}
	runners = append(runners, local...)

	repos := monitoredRepos(runners, effectiveSettings().GitHubHostedRepos)
	statuses, queues, err := loadGitHubStatus(repos)
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
	runners = append(runners, remoteOnlyGitHubRunners(statuses, runners, queues)...)

	hosted, err := loadHostedJobs(repos)
	if err != nil {
		warnings = append(warnings, err)
	}
	runners = append(runners, hosted...)

	sortRunners(runners)

	return Inventory{Runners: runners}, errors.Join(warnings...)
}

func remoteOnlyGitHubRunners(statuses map[string]GitHubRunnerStatus, local []Runner, queues map[string]QueueStatus) []Runner {
	seen := map[string]bool{}
	for _, runner := range local {
		seen[runnerKey(runner.Repo, runner.Name)] = true
	}
	keys := make([]string, 0, len(statuses))
	for key := range statuses {
		if !seen[key] {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	runners := make([]Runner, 0, len(keys))
	for _, key := range keys {
		status := statuses[key]
		repo := firstNonEmpty(status.Repo, repoFromRunnerKey(key))
		if repo == "" {
			continue
		}
		queue := queueForRepo(queues, repo)
		runners = append(runners, Runner{
			Name:            status.Name,
			Repo:            repo,
			OS:              status.OS,
			Host:            "github",
			Path:            "(not local)",
			Transport:       "github-remote",
			LocalState:      "remote",
			ControlMode:     "github-remote",
			GitHubStatus:    firstNonEmpty(status.Status, "unknown"),
			Busy:            status.Busy,
			Labels:          append([]string(nil), status.Labels...),
			Version:         status.Version,
			QueueCount:      queue.Count,
			StaleQueueCount: queue.StaleCount,
		})
	}
	return runners
}

func repoFromRunnerKey(key string) string {
	repo, _, ok := strings.Cut(key, "|")
	if !ok {
		return ""
	}
	return repo
}

func queueForRepo(queues map[string]QueueStatus, repo string) QueueStatus {
	if queue, ok := queues[repo]; ok {
		return queue
	}
	for key, queue := range queues {
		if strings.EqualFold(key, repo) {
			return queue
		}
	}
	return QueueStatus{}
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
	return monitoredRepos(runners, nil)
}

func monitoredRepos(runners []Runner, configured []string) []string {
	seen := map[string]bool{}
	var repos []string
	for _, runner := range runners {
		if runner.Repo == "" {
			continue
		}
		key := strings.ToLower(runner.Repo)
		if seen[key] {
			continue
		}
		seen[key] = true
		repos = append(repos, runner.Repo)
	}
	for _, repo := range configured {
		repo = repoFromGitHubURL(repo)
		if repo == "" {
			continue
		}
		key := strings.ToLower(repo)
		if seen[key] {
			continue
		}
		seen[key] = true
		repos = append(repos, repo)
	}
	sort.Strings(repos)
	return repos
}
