package leakix

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

type subResponse struct {
	Subdomain   string    `json:"subdomain"`
	DistinctIps int       `json:"distinct_ips"`
	LastSeen    time.Time `json:"last_seen"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "leakix"
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

	randomApiKey := s.randomKey()
	// API key is optional for leakix
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://leakix.net/api/subdomains/%s", query), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	if randomApiKey != "" {
		req.Header.Set("api-key", randomApiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var subdomains []subResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&subdomains)
	if err != nil {
		return nil, err
	}

	for _, result := range subdomains {
		results = append(results, result.Subdomain)
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

