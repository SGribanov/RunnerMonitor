package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type workflowRunsForJobsResponse struct {
	WorkflowRuns []struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		HTMLURL   string `json:"html_url"`
		CreatedAt string `json:"created_at"`
	} `json:"workflow_runs"`
}

type workflowJobsResponse struct {
	Jobs []struct {
		ID         int64    `json:"id"`
		Name       string   `json:"name"`
		Status     string   `json:"status"`
		Conclusion string   `json:"conclusion"`
		RunnerName string   `json:"runner_name"`
		Labels     []string `json:"labels"`
		HTMLURL    string   `json:"html_url"`
		StartedAt  string   `json:"started_at"`
		CreatedAt  string   `json:"created_at"`
	} `json:"jobs"`
}

var githubHostedJobsCache = struct {
	sync.Mutex
	key      string
	storedAt time.Time
	runners  []Runner
	err      error
}{}

func LoadGitHubHostedJobs(repos []string) ([]Runner, error) {
	var runners []Runner
	var warnings []string

	for _, repo := range repos {
		for _, status := range []string{"queued", "in_progress"} {
			data, err := ghAPI(fmt.Sprintf("repos/%s/actions/runs?status=%s&per_page=20", repo, status))
			if err != nil {
				warnings = append(warnings, fmt.Sprintf("%s hosted %s runs: %v", repo, status, err))
				continue
			}
			var runs workflowRunsForJobsResponse
			if err := json.Unmarshal(data, &runs); err != nil {
				warnings = append(warnings, fmt.Sprintf("%s hosted %s runs json: %v", repo, status, err))
				continue
			}
			for _, run := range runs.WorkflowRuns {
				jobs, err := loadGitHubHostedJobsForRun(repo, run.ID, run.Name, firstNonEmpty(run.HTMLURL, ""), parseGitHubTime(run.CreatedAt))
				if err != nil {
					warnings = append(warnings, fmt.Sprintf("%s run %d jobs: %v", repo, run.ID, err))
					continue
				}
				runners = append(runners, jobs...)
			}
		}
	}

	if len(warnings) > 0 {
		return runners, errors.New(strings.Join(warnings, "; "))
	}
	return runners, nil
}

func LoadGitHubHostedJobsCached(repos []string, maxAge time.Duration) ([]Runner, error) {
	if maxAge <= 0 {
		return LoadGitHubHostedJobs(repos)
	}
	key := strings.Join(repos, "\x00")
	now := time.Now()

	githubHostedJobsCache.Lock()
	if githubHostedJobsCache.key == key && !githubHostedJobsCache.storedAt.IsZero() && now.Sub(githubHostedJobsCache.storedAt) <= maxAge {
		runners := cloneRunners(githubHostedJobsCache.runners)
		err := githubHostedJobsCache.err
		githubHostedJobsCache.Unlock()
		return runners, err
	}
	githubHostedJobsCache.Unlock()

	runners, err := LoadGitHubHostedJobs(repos)

	githubHostedJobsCache.Lock()
	githubHostedJobsCache.key = key
	githubHostedJobsCache.storedAt = now
	githubHostedJobsCache.runners = cloneRunners(runners)
	githubHostedJobsCache.err = err
	cached := cloneRunners(githubHostedJobsCache.runners)
	cachedErr := githubHostedJobsCache.err
	githubHostedJobsCache.Unlock()
	return cached, cachedErr
}

func loadGitHubHostedJobsForRun(repo string, runID int64, workflowName string, runURL string, runCreatedAt time.Time) ([]Runner, error) {
	data, err := ghAPI(fmt.Sprintf("repos/%s/actions/runs/%d/jobs?per_page=100", repo, runID))
	if err != nil {
		return nil, err
	}
	return hostedJobsFromResponse(repo, workflowName, runURL, runCreatedAt, data)
}

func hostedJobsFromResponse(repo, workflowName, runURL string, runCreatedAt time.Time, data []byte) ([]Runner, error) {
	var response workflowJobsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	runners := make([]Runner, 0, len(response.Jobs))
	for _, job := range response.Jobs {
		if job.Status == "completed" || hasSelfHostedLabel(job.Labels) {
			continue
		}
		createdAt := firstTime(parseGitHubTime(job.CreatedAt), parseGitHubTime(job.StartedAt), runCreatedAt)
		runners = append(runners, hostedJobRunner(repo, workflowName, job.Name, job.Status, job.RunnerName, job.Labels, firstNonEmpty(job.HTMLURL, runURL), createdAt))
	}
	return runners, nil
}

func hostedJobRunner(repo, workflowName, jobName, status, runnerName string, labels []string, url string, createdAt time.Time) Runner {
	name := strings.TrimSpace(jobName)
	if strings.TrimSpace(workflowName) != "" && !strings.EqualFold(workflowName, jobName) {
		name = strings.TrimSpace(workflowName) + " / " + name
	}
	if name == "" {
		name = "GitHub-hosted job"
	}
	if strings.TrimSpace(runnerName) != "" {
		name += " @ " + strings.TrimSpace(runnerName)
	}
	queue := 0
	stale := 0
	if strings.EqualFold(status, "queued") {
		queue = 1
		if !createdAt.IsZero() && createdAt.Before(time.Now().Add(-30*time.Minute)) {
			stale = 1
		}
	}
	return Runner{
		Name:            name,
		Repo:            repo,
		OS:              hostedOS(labels),
		Host:            "github",
		Path:            url,
		Transport:       "github-hosted",
		LocalState:      "hosted",
		ControlMode:     "github-hosted",
		GitHubStatus:    firstNonEmpty(status, "unknown"),
		Busy:            strings.EqualFold(status, "in_progress"),
		Labels:          append([]string(nil), labels...),
		QueueCount:      queue,
		StaleQueueCount: stale,
	}
}

func hasSelfHostedLabel(labels []string) bool {
	for _, label := range labels {
		if strings.EqualFold(strings.TrimSpace(label), "self-hosted") {
			return true
		}
	}
	return false
}

func hostedOS(labels []string) string {
	for _, label := range labels {
		normalized := strings.ToLower(strings.TrimSpace(label))
		switch {
		case strings.Contains(normalized, "windows"):
			return "Windows"
		case strings.Contains(normalized, "ubuntu") || strings.Contains(normalized, "linux"):
			return "Linux"
		case strings.Contains(normalized, "macos"):
			return "macOS"
		}
	}
	return "GitHub"
}

func parseGitHubTime(value string) time.Time {
	if strings.TrimSpace(value) == "" {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func firstTime(values ...time.Time) time.Time {
	for _, value := range values {
		if !value.IsZero() {
			return value
		}
	}
	return time.Time{}
}

func cloneRunners(runners []Runner) []Runner {
	clone := make([]Runner, len(runners))
	for i, runner := range runners {
		runner.Labels = append([]string(nil), runner.Labels...)
		clone[i] = runner
	}
	return clone
}
