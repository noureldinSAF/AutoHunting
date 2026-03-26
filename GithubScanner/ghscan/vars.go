package main

import (
	"regexp"
	"sync"
)

var mu sync.Mutex

// Scan represents a scan process.
type Scan struct {
	APIKeys        []string
	Target         []string
	LimitedKeys    []string
	CurrentUsedKey string
	Errorfs        chan error
}

// RepoInfo is a repository returned by ListOrgRepos (org repos API).
type RepoInfo struct {
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
	HTMLURL  string `json:"html_url"`
}

// Regex config loaded from regexes.yaml
type regexSignaturesFile struct {
	Signatures []struct {
		Pattern struct {
			Name      string `yaml:"name"`
			Value     string `yaml:"value"`
			Sensitive bool   `yaml:"sensitive"`
		} `yaml:"pattern"`
	} `yaml:"signatures"`
}

// CompiledPattern holds a compiled regex and its metadata for scanning.
type CompiledPattern struct {
	Name  string
	Regex *regexp.Regexp
}

// MatchResult is a single regex match found in file content.
type MatchResult struct {
	Rule       string `json:"rule"`
	Matched    string `json:"matched"`
	CommitURL  string `json:"commit_url"`
}

// OrgReport represents scan results for a single organization.
type OrgReport struct {
	Org   string       `json:"org"`
	Repos []RepoReport `json:"repos"`
}

// RepoReport represents scan results for a single repository.
type RepoReport struct {
	RepoURL string        `json:"repo_url"`
	Matches []MatchResult `json:"matches"`
}
