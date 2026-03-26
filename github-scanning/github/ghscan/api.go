package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	githubAPIBase    = "https://api.github.com"
	githubAPIVersion = "2022-11-28"
)

// ListOrgRepos returns all repository names and clone URLs for the given org.
func (s *Scan) ListOrgRepos(org string) ([]RepoInfo, error) {
	var all []RepoInfo
	page := 1
	const perPage = 100

	for {
		url := fmt.Sprintf("%s/orgs/%s/repos?per_page=%d&page=%d", githubAPIBase, org, perPage, page)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		s.setAuthHeaders(req)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("org repos: HTTP %d: %s", resp.StatusCode, string(body))
		}

		var pageRepos []RepoInfo
		if err := json.Unmarshal(body, &pageRepos); err != nil {
			return nil, err
		}

		if len(pageRepos) == 0 {
			break
		}
		all = append(all, pageRepos...)

		if len(pageRepos) < perPage {
			break
		}
		page++
	}

	return all, nil
}

// getRemainingRequests retrieves the remaining number of requests allowed for the current API key.
func (s *Scan) getRemainingRequests() int {
	req, err := http.NewRequest(http.MethodGet, githubAPIBase+"/rate_limit", nil)
	if err != nil {
		s.sendError(err)
		return 0
	}

	s.setAuthHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.sendError(err)
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.sendError(err)
		return 0
	}

	var rateLimitResponse struct {
		Resources struct {
			CodeSearch struct {
				Remaining int `json:"remaining"`
			} `json:"code_search"`
		} `json:"resources"`
	}

	if err := json.Unmarshal(body, &rateLimitResponse); err != nil {
		s.sendError(err)
		return 0
	}

	return rateLimitResponse.Resources.CodeSearch.Remaining
}

// isAliveKey checks if the provided API key is valid and active.
func isAliveKey(key string) bool {
	req, err := http.NewRequest(http.MethodGet, githubAPIBase+"/user", nil)
	if err != nil {
		return false
	}

	req.Header.Set("Authorization", "Bearer "+key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	return strings.Contains(string(body), "login")
}

// setAuthHeaders sets the standard GitHub API headers including authentication.
func (s *Scan) setAuthHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+s.CurrentUsedKey)
	req.Header.Set("X-GitHub-Api-Version", githubAPIVersion)
}

// sendError sends an error to the error channel without blocking.
func (s *Scan) sendError(err error) {
	select {
	case s.Errorfs <- err:
	default:
	}
}
