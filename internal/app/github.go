package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

type GitHubRunnerStatus struct {
	Name    string
	OS      string
	Status  string
	Busy    bool
	Labels  []string
	Version string
}

type QueueStatus struct {
	Count      int
	StaleCount int
}

type runnersResponse struct {
	Runners []struct {
		Name    string `json:"name"`
		OS      string `json:"os"`
		Status  string `json:"status"`
		Busy    bool   `json:"busy"`
		Version string `json:"version"`
		Labels  []struct {
			Name string `json:"name"`
		} `json:"labels"`
	} `json:"runners"`
}

type runsResponse struct {
	TotalCount   int `json:"total_count"`
	WorkflowRuns []struct {
		CreatedAt time.Time `json:"created_at"`
	} `json:"workflow_runs"`
}

func LoadGitHubStatus(repos []string) (map[string]GitHubRunnerStatus, map[string]QueueStatus, error) {
	statuses := map[string]GitHubRunnerStatus{}
	queues := map[string]QueueStatus{}
	var warnings []string

	for _, repo := range repos {
		runnerData, err := ghAPI(fmt.Sprintf("repos/%s/actions/runners", repo))
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s runners: %v", repo, err))
			continue
		}
		var response runnersResponse
		if err := json.Unmarshal(runnerData, &response); err != nil {
			warnings = append(warnings, fmt.Sprintf("%s runners json: %v", repo, err))
			continue
		}
		for _, runner := range response.Runners {
			var labels []string
			for _, label := range runner.Labels {
				labels = append(labels, label.Name)
			}
			statuses[runnerKey(repo, runner.Name)] = GitHubRunnerStatus{
				Name:    runner.Name,
				OS:      runner.OS,
				Status:  runner.Status,
				Busy:    runner.Busy,
				Labels:  labels,
				Version: runner.Version,
			}
		}

		queueData, err := ghAPI(fmt.Sprintf("repos/%s/actions/runs?status=queued&per_page=100", repo))
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s queue: %v", repo, err))
			continue
		}
		var runs runsResponse
		if err := json.Unmarshal(queueData, &runs); err != nil {
			warnings = append(warnings, fmt.Sprintf("%s queue json: %v", repo, err))
			continue
		}
		queue := QueueStatus{Count: runs.TotalCount}
		staleBefore := time.Now().Add(-30 * time.Minute)
		for _, run := range runs.WorkflowRuns {
			if run.CreatedAt.Before(staleBefore) {
				queue.StaleCount++
			}
		}
		queues[repo] = queue
	}

	if len(warnings) > 0 {
		return statuses, queues, errors.New(strings.Join(warnings, "; "))
	}
	return statuses, queues, nil
}

func ghAPI(endpoint string) ([]byte, error) {
	return exec.Command("gh", "api", endpoint).Output()
}

func ghAPIMethod(method string, endpoint string) ([]byte, error) {
	return exec.Command("gh", "api", "-X", method, endpoint).Output()
}

func repoFromGitHubURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err == nil && parsed.Host == "github.com" {
		return strings.Trim(strings.TrimSuffix(parsed.Path, ".git"), "/")
	}
	if strings.HasPrefix(raw, "git@github.com:") {
		return strings.TrimSuffix(strings.TrimPrefix(raw, "git@github.com:"), ".git")
	}
	return raw
}

func runnerKey(repo, name string) string {
	return strings.ToLower(repo + "|" + name)
}
