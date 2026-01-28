package urlscan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Source struct{}

func (s *Source) Name() string {
	return "urlscan"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	url := fmt.Sprintf("https://urlscan.io/api/v1/search/?q=page.domain:%s&size=10000", query)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var response struct {
		Results []struct {
			Task struct {
				Domain string `json:"domain"`
			} `json:"task"`
			Page struct {
				Domain string `json:"domain"`
				URL    string `json:"url"`
			} `json:"page"`
		} `json:"results"`
		Total   int  `json:"total"`
		HasMore bool `json:"has_more"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w (body: %s)", err, string(body[:min(len(body), 200)]))
	}

	if len(response.Results) == 0 {
		return results, nil
	}

	// Collect unique domains
	seen := make(map[string]bool)

	for _, entry := range response.Results {
		// Extract from page.domain
		if entry.Page.Domain != "" {
			seen[entry.Page.Domain] = true
		}

		// Extract from page.url (might contain subdomains)
		if entry.Page.URL != "" {
			urlDomain := extractDomainFromURL(entry.Page.URL)
			if urlDomain != "" {
				seen[urlDomain] = true
			}
		}
	}

	// Send unique domains that match our target
	for subdomain := range seen {
		if subdomain != "" && (strings.HasSuffix(subdomain, "."+query) || subdomain == query) {
			results = append(results, subdomain)
		}
	}

	return results, nil
}

// extractDomainFromURL extracts the domain from a URL string
func extractDomainFromURL(urlStr string) string {
	// Remove protocol
	if idx := strings.Index(urlStr, "://"); idx != -1 {
		urlStr = urlStr[idx+3:]
	}

	// Remove path
	if idx := strings.Index(urlStr, "/"); idx != -1 {
		urlStr = urlStr[:idx]
	}

	// Remove port
	if idx := strings.Index(urlStr, ":"); idx != -1 {
		urlStr = urlStr[:idx]
	}

	return urlStr
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

