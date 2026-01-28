package commoncrawl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"subenum/pkg/utils"
)

type Source struct{}

const (
	indexURL     = "https://index.commoncrawl.org/collinfo.json"
	maxYearsBack = 5
)

var year = time.Now().Year()

type indexResponse struct {
	ID     string `json:"id"`
	APIURL string `json:"cdx-api"`
}

func (s *Source) Name() string {
	return "commoncrawl"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var allResults []string
	seen := make(map[string]bool)

	// Get index list
	req, err := http.NewRequest(http.MethodGet, indexURL, nil)
	if err != nil {
		return allResults, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return allResults, err
	}

	var indexes []indexResponse
	err = json.NewDecoder(resp.Body).Decode(&indexes)
	resp.Body.Close()
	if err != nil {
		return allResults, err
	}

	years := make([]string, 0)
	for i := 0; i < maxYearsBack; i++ {
		years = append(years, strconv.Itoa(year-i))
	}

	searchIndexes := make(map[string]string)
	for _, year := range years {
		for _, index := range indexes {
			if strings.Contains(index.ID, year) {
				if _, ok := searchIndexes[year]; !ok {
					searchIndexes[year] = index.APIURL
					break
				}
			}
		}
	}

	for _, apiURL := range searchIndexes {
		results, err := s.getSubdomains(client, apiURL, query)
		if err != nil {
			continue
		}
		for _, result := range results {
			if !seen[result] {
				seen[result] = true
				allResults = append(allResults, result)
			}
		}
	}

	return allResults, nil
}

func (s *Source) getSubdomains(client *http.Client, searchURL, domain string) ([]string, error) {
	var results []string
	seen := make(map[string]bool)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?url=*.%s", searchURL, domain), nil)
	if err != nil {
		return results, err
	}

	req.Header.Set("Host", "index.commoncrawl.org")

	resp, err := client.Do(req)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	// Read entire body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}

	bodyStr := string(body)
	// Decode URL-encoded content
	bodyStr, _ = url.QueryUnescape(bodyStr)

	// Create subdomain extractor for this domain
	extractor, err := utils.NewSubdomainExtractor(domain)
	if err != nil {
		return results, err
	}

	// Extract subdomains from the response
	extracted := extractor.Extract(bodyStr)

	// Deduplicate and clean results
	for _, subdomain := range extracted {
		// Fix for triple encoded URL artifacts
		cleaned := strings.ToLower(subdomain)
		cleaned = strings.TrimPrefix(cleaned, "25")
		cleaned = strings.TrimPrefix(cleaned, "2f")
		if cleaned != "" && !seen[cleaned] {
			seen[cleaned] = true
			results = append(results, cleaned)
		}
	}

	return results, nil
}

