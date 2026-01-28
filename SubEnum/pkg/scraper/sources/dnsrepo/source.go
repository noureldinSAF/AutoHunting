package dnsrepo

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Source struct {
	apiKeys []string
}

type DnsRepoResponse []struct {
	Domain string `json:"domain"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "dnsrepo"
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
		return nil, fmt.Errorf("no API keys provided for dnsrepo")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	randomApiInfo := strings.Split(randomApiKey, ":")
	if len(randomApiInfo) != 2 {
		return nil, fmt.Errorf("invalid API key format, expected token:apikey")
	}

	token := randomApiInfo[0]
	apiKey := randomApiInfo[1]

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://dnsarchive.net/api/?apikey=%s&search=%s", apiKey, query), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Access", token)

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

	var result DnsRepoResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	for _, sub := range result {
		results = append(results, strings.TrimSuffix(sub.Domain, "."))
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

