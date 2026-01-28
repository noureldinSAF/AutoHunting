package dnsdumpster

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

type response struct {
	A []struct {
		Host string `json:"host"`
	} `json:"a"`
	Ns []struct {
		Host string `json:"host"`
	} `json:"ns"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "dnsdumpster"
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
		return nil, fmt.Errorf("no API keys provided for dnsdumpster")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.dnsdumpster.com/domain/%s", query), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", randomApiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	for _, record := range append(response.A, response.Ns...) {
		results = append(results, record.Host)
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

