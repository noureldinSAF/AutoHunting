package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cyinnove/logify"
	"gopkg.in/yaml.v3"
)

// RepoFromSlug converts "owner/repo" or a GitHub URL into RepoInfo.
func RepoFromSlug(slug string) (RepoInfo, error) {
	s := strings.TrimSpace(slug)
	if s == "" {
		return RepoInfo{}, fmt.Errorf("empty repo slug")
	}

	// Allow full URL: https://github.com/owner/repo or https://github.com/owner/repo/
	if strings.Contains(s, "github.com/") {
		return parseGitHubURL(s)
	}

	// owner/repo format
	return parseOwnerRepo(s)
}

// parseGitHubURL parses a GitHub URL into RepoInfo.
func parseGitHubURL(url string) (RepoInfo, error) {
	parts := strings.Split(url, "github.com/")
	if len(parts) != 2 {
		return RepoInfo{}, fmt.Errorf("invalid github URL: %s", url)
	}

	path := strings.Trim(parts[1], "/")
	segments := strings.SplitN(path, "/", 2)
	if len(segments) != 2 {
		return RepoInfo{}, fmt.Errorf("invalid github URL path: %s", url)
	}

	owner, repo := segments[0], strings.TrimSuffix(segments[1], ".git")
	fullName := owner + "/" + repo

	return RepoInfo{
		FullName: fullName,
		CloneURL: fmt.Sprintf("https://github.com/%s.git", fullName),
		HTMLURL:  fmt.Sprintf("https://github.com/%s", fullName),
	}, nil
}

// parseOwnerRepo parses an "owner/repo" string into RepoInfo.
func parseOwnerRepo(s string) (RepoInfo, error) {
	segments := strings.SplitN(s, "/", 2)
	if len(segments) != 2 || segments[0] == "" || segments[1] == "" {
		return RepoInfo{}, fmt.Errorf("repo must be owner/repo: %s", s)
	}

	fullName := segments[0] + "/" + segments[1]
	return RepoInfo{
		FullName: fullName,
		CloneURL: fmt.Sprintf("https://github.com/%s.git", fullName),
		HTMLURL:  fmt.Sprintf("https://github.com/%s", fullName),
	}, nil
}

// GetRawFile converts a GitHub blob/tree URL to a raw content URL.
func GetRawFile(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) < 7 {
		return ""
	}

	var blobIndex int
	for i, part := range parts {
		if part == "blob" || part == "tree" {
			blobIndex = i
			break
		}
	}

	if blobIndex == 0 || blobIndex+2 >= len(parts) {
		return ""
	}

	owner := parts[3]
	repo := parts[4]
	branchOrTag := parts[blobIndex+1]
	filePath := strings.Join(parts[blobIndex+2:], "/")

	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, branchOrTag, filePath)
}

// GetRepoURL converts a raw GitHub URL to a regular GitHub blob URL.
func GetRepoURL(rawURL string) string {
	parts := strings.Split(rawURL, "/")
	if len(parts) < 7 {
		return ""
	}

	owner := parts[3]
	repo := parts[4]
	branchOrTag := parts[5]
	filePath := strings.Join(parts[6:], "/")

	return fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s", owner, repo, branchOrTag, filePath)
}

// LoadRegexes reads regexes.yaml and returns compiled patterns for scanning.
func LoadRegexes(path string) ([]CompiledPattern, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read regexes file: %w", err)
	}

	var file regexSignaturesFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse regexes yaml: %w", err)
	}

	patterns := make([]CompiledPattern, 0, len(file.Signatures))
	for _, sig := range file.Signatures {
		value := strings.TrimSpace(sig.Pattern.Value)
		if value == "" {
			continue
		}

		re, err := regexp.Compile(value)
		if err != nil {
			logify.Warningf("skip invalid regex %q: %v", sig.Pattern.Name, err)
			continue
		}

		patterns = append(patterns, CompiledPattern{
			Name:  sig.Pattern.Name,
			Regex: re,
		})
	}

	return patterns, nil
}

   