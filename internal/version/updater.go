package version

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/semver/v3"
)

type UpdateInfo struct {
	Version      string
	ReleaseURL   string
	ReleaseNotes string
	ForceUpdate  bool // Will be true for minor or major version changes
	HasUpdate    bool
}

// CheckForUpdates checks if a newer version is available
func CheckForUpdates(ctx context.Context, currentVersion string) (*UpdateInfo, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET",
		"https://api.github.com/repos/martijnspitter/tui-todo/releases/latest", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
		Body    string `json:"body"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	// Parse versions and compare
	if currentVersion == "dev" {
		currentVersion = "v0.0.1"
	}
	if len(currentVersion) > 0 && currentVersion[0] == 'v' {
		currentVersion = currentVersion[1:]
	}
	currentSemver, err := semver.NewVersion(currentVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid current version: %w", err)
	}

	// Remove "v" prefix if present
	tagVersion := release.TagName
	if len(tagVersion) > 0 && tagVersion[0] == 'v' {
		tagVersion = tagVersion[1:]
	}

	latestSemver, err := semver.NewVersion(tagVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid tag version: %w", err)
	}

	if latestSemver.GreaterThan(currentSemver) {
		// Determine if this is a major or minor version change
		isMajorUpdate := latestSemver.Major() > currentSemver.Major()
		isMinorUpdate := latestSemver.Minor() > currentSemver.Minor()
		forceUpdate := isMajorUpdate || isMinorUpdate

		return &UpdateInfo{
			Version:      currentVersion,
			ReleaseURL:   release.HTMLURL,
			ReleaseNotes: release.Body,
			ForceUpdate:  forceUpdate,
			HasUpdate:    true,
		}, nil
	}

	return &UpdateInfo{
		Version:      currentVersion,
		ReleaseURL:   release.HTMLURL,
		ReleaseNotes: release.Body,
		ForceUpdate:  false,
		HasUpdate:    false,
	}, nil
}
