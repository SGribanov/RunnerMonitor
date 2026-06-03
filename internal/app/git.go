package app

import (
	"fmt"
	"os/exec"
	"strings"
)

func CurrentGitHubRepo() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", err
	}
	repo := repoFromGitHubURL(strings.TrimSpace(string(out)))
	if repo == "" || !strings.Contains(repo, "/") {
		return "", fmt.Errorf("origin is not a GitHub owner/repo remote")
	}
	return repo, nil
}
