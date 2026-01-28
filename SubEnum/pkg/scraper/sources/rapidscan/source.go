package rapidscan

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Source struct {
	apiKeys []string
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "rapidscan"
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
		return nil, fmt.Errorf("no API keys provided for rapidscan")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	url := fmt.Sprintf("https://subdomain-scan1.p.rapidapi.com/?domain=%s", query)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-RapidAPI-Key", randomApiKey)
	req.Header.Set("X-RapidAPI-Host", "subdomain-scan1.p.rapidapi.com")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var subdomains []string
	err = json.NewDecoder(resp.Body).Decode(&subdomains)
	if err != nil {
		return nil, err
	}

	results = append(results, subdomains...)
	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

