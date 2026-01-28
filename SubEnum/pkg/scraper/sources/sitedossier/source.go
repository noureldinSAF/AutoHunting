package sitedossier

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"subenum/pkg/utils"
)

type Source struct{}

var reNext = regexp.MustCompile(`<a href="([A-Za-z0-9/.]+)"><b>`)

func (s *Source) Name() string {
	return "sitedossier"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var allResults []string
	seen := make(map[string]bool)
	baseURL := fmt.Sprintf("http://www.sitedossier.com/parentdomain/%s", query)
	
	results, err := s.enumerate(client, baseURL, query, seen)
	if err != nil {
		return allResults, err
	}
	
	allResults = append(allResults, results...)
	return allResults, nil
}

func (s *Source) enumerate(client *http.Client, baseURL string, domain string, seen map[string]bool) ([]string, error) {
	var results []string

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return results, err
	}

	resp, err := client.Do(req)
	isnotfound := resp != nil && resp.StatusCode == http.StatusNotFound
	if err != nil && !isnotfound {
		return results, err
	}
	defer resp.Body.Close()

	// Read all body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}
	src := string(body)

	// Create subdomain extractor for this domain
	extractor, err := utils.NewSubdomainExtractor(domain)
	if err != nil {
		return results, err
	}

	// Extract subdomains from the HTML
	extracted := extractor.Extract(src)
	for _, subdomain := range extracted {
		if !seen[subdomain] {
			seen[subdomain] = true
			results = append(results, subdomain)
		}
	}

	match := reNext.FindStringSubmatch(src)
	if len(match) > 0 {
		nextResults, err := s.enumerate(client, fmt.Sprintf("http://www.sitedossier.com%s", match[1]), domain, seen)
		if err == nil {
			results = append(results, nextResults...)
		}
	}

	return results, nil
}

