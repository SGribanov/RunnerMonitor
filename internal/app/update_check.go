package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type latestRelease struct {
	TagName string `json:"tagName"`
	URL     string `json:"url"`
}

func CheckForUpdate(currentVersion string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	release, err := latestReleaseInfo(ctx)
	if err != nil {
		return ""
	}
	return updateNotice(currentVersion, release.TagName, release.URL)
}

func latestReleaseInfo(ctx context.Context) (latestRelease, error) {
	cmd := exec.CommandContext(ctx, "gh", "release", "view", "--repo", ReleaseRepo, "--json", "tagName,url")
	out, err := cmd.Output()
	if err != nil {
		return latestRelease{}, err
	}
	var release latestRelease
	if err := json.Unmarshal(out, &release); err != nil {
		return latestRelease{}, err
	}
	if strings.TrimSpace(release.TagName) == "" {
		return latestRelease{}, fmt.Errorf("latest release has no tag")
	}
	return release, nil
}

func updateNotice(currentVersion string, latestVersion string, releaseURL string) string {
	current := strings.TrimSpace(currentVersion)
	latest := strings.TrimSpace(latestVersion)
	if current == "" || latest == "" || !isNewerVersion(latest, current) {
		return ""
	}
	if strings.TrimSpace(releaseURL) == "" {
		return fmt.Sprintf("update available: %s -> %s", current, latest)
	}
	return fmt.Sprintf("update available: %s -> %s (%s)", current, latest, releaseURL)
}

func isNewerVersion(candidate string, current string) bool {
	candidateParts := versionParts(candidate)
	currentParts := versionParts(current)
	for i := 0; i < max(len(candidateParts), len(currentParts)); i++ {
		candidatePart := 0
		if i < len(candidateParts) {
			candidatePart = candidateParts[i]
		}
		currentPart := 0
		if i < len(currentParts) {
			currentPart = currentParts[i]
		}
		if candidatePart > currentPart {
			return true
		}
		if candidatePart < currentPart {
			return false
		}
	}
	return false
}

func versionParts(version string) []int {
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	parts := strings.Split(version, ".")
	values := make([]int, 0, len(parts))
	for _, part := range parts {
		part = leadingDigits(part)
		if part == "" {
			values = append(values, 0)
			continue
		}
		value, err := strconv.Atoi(part)
		if err != nil {
			value = 0
		}
		values = append(values, value)
	}
	return values
}

func leadingDigits(value string) string {
	for i, r := range value {
		if r < '0' || r > '9' {
			return value[:i]
		}
	}
	return value
}
