package rapidfinder

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type Source struct {
	apiKeys []string
}

type rapidfinderResponse struct {
	Subdomains []struct {
		Subdomain string `json:"subdomain"`
	} `json:"subdomains"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "rapidfinder"
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
	var results []string

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for rapidfinder")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	url := fmt.Sprintf("https://subdomain-finder5.p.rapidapi.com/v1/subdomain-finder/?domain=%s", query)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-RapidAPI-Key", randomApiKey)
	req.Header.Set("X-RapidAPI-Host", "subdomain-finder5.p.rapidapi.com")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: status=%d body=%s", resp.StatusCode, string(body))
	}

	var response rapidfinderResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	for _, item := range response.Subdomains {
		if item.Subdomain != "" {
			results = append(results, item.Subdomain)
		}
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

