package gitlab

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Source struct {
	apiKeys []string
}

type item struct {
	Data      string `json:"data"`
	ProjectId int    `json:"project_id"`
	Path      string `json:"path"`
	Ref       string `json:"ref"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "gitlab"
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
		return nil, fmt.Errorf("no API keys provided for gitlab")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	headers := map[string]string{"PRIVATE-TOKEN": randomApiKey}
	searchURL := fmt.Sprintf("https://gitlab.com/api/v4/search?scope=blobs&search=%s&per_page=100", url.QueryEscape(query))
	domainRegexp := domainRegexp(query)

	results, err := s.enumerate(client, searchURL, domainRegexp, headers)
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

func (s *Source) enumerate(client *http.Client, searchURL string, domainRegexp *regexp.Regexp, headers map[string]string) ([]string, error) {
	var allResults []string

	req, err := http.NewRequest(http.MethodGet, searchURL, nil)
	if err != nil {
		return allResults, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return allResults, err
	}
	defer resp.Body.Close()

	var items []item
	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		return allResults, err
	}

	for _, it := range items {
		fileUrl := fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/repository/files/%s/raw?ref=%s",
			it.ProjectId, url.QueryEscape(it.Path), it.Ref)

		req, err := http.NewRequest(http.MethodGet, fileUrl, nil)
		if err != nil {
			continue
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

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
				for _, subdomain := range domainRegexp.FindAllString(line, -1) {
					allResults = append(allResults, subdomain)
				}
			}
		}
		resp.Body.Close()
	}

	return allResults, nil
}

func domainRegexp(domain string) *regexp.Regexp {
	rdomain := strings.ReplaceAll(domain, ".", "\\.")
	return regexp.MustCompile("(\\w[a-zA-Z0-9][a-zA-Z0-9-\\.]*)" + rdomain)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

