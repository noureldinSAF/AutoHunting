package github

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Source struct {
	apiKeys []string
}

type textMatch struct {
	Fragment string `json:"fragment"`
}

type item struct {
	Name        string      `json:"name"`
	HTMLURL     string      `json:"html_url"`
	TextMatches []textMatch `json:"text_matches"`
}

type response struct {
	TotalCount int    `json:"total_count"`
	Items      []item `json:"items"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "github"
}

func (s *Source) RequiresAPIKey() bool {
	return true
}

func (s *Source) randomKey() string {
	if len(s.apiKeys) == 0 {
		return ""
	}
	return s.apiKeys[rand.Intn(len(s.apiKeys))]
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var allResults []string
	seen := make(map[string]bool)

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for github")
	}

	token := s.randomKey()
	if token == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	searchURL := fmt.Sprintf("https://api.github.com/search/code?per_page=100&q=%s&sort=created&order=asc", url.QueryEscape(query))
	domainRegexp := domainRegexp(query)

	results, err := s.enumerate(client, searchURL, domainRegexp, token)
	if err != nil {
		return allResults, err
	}

	for _, result := range results {
		if !seen[result] {
			seen[result] = true
			allResults = append(allResults, result)
		}
	}

	return allResults, nil
}

func (s *Source) enumerate(client *http.Client, searchURL string, domainRegexp *regexp.Regexp, token string) ([]string, error) {
	var allResults []string

	req, err := http.NewRequest(http.MethodGet, searchURL, nil)
	if err != nil {
		return allResults, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3.text-match+json")
	req.Header.Set("Authorization", "token "+token)

	resp, err := client.Do(req)
	if err != nil {
		return allResults, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return allResults, fmt.Errorf("rate limit exceeded")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return allResults, err
	}

	var data response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return allResults, err
	}

	results, err := s.processItems(client, data.Items, domainRegexp, token)
	if err != nil {
		return allResults, err
	}

	allResults = append(allResults, results...)

	// Check for next page in Link header
	linkHeader := resp.Header.Get("Link")
	if strings.Contains(linkHeader, `rel="next"`) {
		// Extract next URL from Link header (simplified)
		// In production, use a proper Link header parser
	}

	return allResults, nil
}

func (s *Source) processItems(client *http.Client, items []item, domainRegexp *regexp.Regexp, token string) ([]string, error) {
	var results []string

	for _, responseItem := range items {
		// Get raw file content
		rawURL := strings.ReplaceAll(responseItem.HTMLURL, "https://github.com/", "https://raw.githubusercontent.com/")
		rawURL = strings.ReplaceAll(rawURL, "/blob/", "/")

		req, err := http.NewRequest(http.MethodGet, rawURL, nil)
		if err != nil {
			continue
		}

		req.Header.Set("Authorization", "token "+token)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode == http.StatusOK {
			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}
				for _, subdomain := range domainRegexp.FindAllString(normalizeContent(line), -1) {
					results = append(results, subdomain)
				}
			}
		}
		resp.Body.Close()

		// Process text matches
		for _, textMatch := range responseItem.TextMatches {
			for _, subdomain := range domainRegexp.FindAllString(normalizeContent(textMatch.Fragment), -1) {
				results = append(results, subdomain)
			}
		}
	}

	return results, nil
}

func normalizeContent(content string) string {
	normalizedContent, _ := url.QueryUnescape(content)
	normalizedContent = strings.ReplaceAll(normalizedContent, "\\t", "")
	normalizedContent = strings.ReplaceAll(normalizedContent, "\\n", "")
	return normalizedContent
}

func domainRegexp(domain string) *regexp.Regexp {
	rdomain := strings.ReplaceAll(domain, ".", "\\.")
	return regexp.MustCompile("(\\w[a-zA-Z0-9][a-zA-Z0-9-\\.]*)" + rdomain)
}
