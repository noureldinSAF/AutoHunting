package crtsh

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	// postgres driver

	jsoniter "github.com/json-iterator/go"

	_ "github.com/lib/pq"
)

type Source struct{}

type subdomain struct {
	ID        int    `json:"id"`
	NameValue string `json:"name_value"`
}

func (s *Source) Name() string {
	return "crtsh"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {

	// Fallback to HTTP
	httpResults, err := s.getSubdomainsFromHTTP(query, client)
	if err != nil {
		return nil, err
	}

	return httpResults, nil
}

func (s *Source) getSubdomainsFromHTTP(domain string, client *http.Client) ([]string, error) {
	var results []string

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var subdomains []subdomain
	err = jsoniter.Unmarshal(body, &subdomains)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	for _, subdomain := range subdomains {
		for sub := range strings.SplitSeq(subdomain.NameValue, "\n") {
			sub = strings.TrimSpace(sub)
			if sub == "" {
				continue
			}
			// Filter out wildcard entries (e.g., *.domain.com)
			if strings.HasPrefix(sub, "*.") {
				continue
			}
			// Skip if it's a duplicate
			if !seen[sub] {
				seen[sub] = true
				results = append(results, sub)
			}
		}
	}

	return results, nil
}
