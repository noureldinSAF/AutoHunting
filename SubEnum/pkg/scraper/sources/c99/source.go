package c99

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Source struct {
	apiKeys []string
}

type dnsdbLookupResponse struct {
	Success    bool `json:"success"`
	Subdomains []struct {
		Subdomain  string `json:"subdomain"`
		IP         string `json:"ip"`
		Cloudflare bool   `json:"cloudflare"`
	} `json:"subdomains"`
	Error string `json:"error"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "c99"
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
		return nil, fmt.Errorf("no API keys provided for c99")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	searchURL := fmt.Sprintf("https://api.c99.nl/subdomainfinder?key=%s&domain=%s&json", randomApiKey, query)
	req, err := http.NewRequest(http.MethodGet, searchURL, nil)
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: status=%d body=%s", resp.StatusCode, string(body))
	}

	var response dnsdbLookupResponse
	err = jsoniter.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, fmt.Errorf("%v", response.Error)
	}

	for _, data := range response.Subdomains {
		if !strings.HasPrefix(data.Subdomain, ".") {
			results = append(results, data.Subdomain)
		}
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

