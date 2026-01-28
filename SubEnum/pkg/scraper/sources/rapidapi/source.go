package rapidapi

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
	apiKeys   []string
	whoisKeys []string
}

type rapidapiResponse struct {
	Result struct {
		Records []struct {
			Domain string `json:"domain"`
		} `json:"records"`
	} `json:"result"`
}

func New(apiKeys []string) *Source {
	// Note: rapidapi needs both rapidapi keys and whoisxmlapi keys
	// Format: "rapidapi_key:whoisxmlapi_key"
	rapidKeys := make([]string, 0)
	whoisKeys := make([]string, 0)
	
	for _, key := range apiKeys {
		parts := strings.Split(key, ":")
		if len(parts) == 2 {
			rapidKeys = append(rapidKeys, parts[0])
			whoisKeys = append(whoisKeys, parts[1])
		}
	}
	
	return &Source{
		apiKeys:   rapidKeys,
		whoisKeys: whoisKeys,
	}
}

func (s *Source) Name() string {
	return "rapidapi"
}

func (s *Source) RequiresAPIKey() bool {
	return true
}

func (s *Source) randomKey() (string, string) {
	if len(s.apiKeys) == 0 || len(s.whoisKeys) == 0 {
		return "", ""
	}
	rapidIdx := rand.Intn(len(s.apiKeys))
	whoisIdx := rand.Intn(len(s.whoisKeys))
	return s.apiKeys[rapidIdx], s.whoisKeys[whoisIdx]
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	if len(s.apiKeys) == 0 || len(s.whoisKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for rapidapi (needs both rapidapi and whoisxmlapi keys)")
	}

	randomApiKey, randomWhoisKey := s.randomKey()
	if randomApiKey == "" || randomWhoisKey == "" {
		return nil, fmt.Errorf("no valid API keys available")
	}

	url := fmt.Sprintf("https://subdomains-lookup.p.rapidapi.com/api/v1?domainName=%s&apiKey=%s&outputFormat=JSON", query, randomWhoisKey)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-RapidAPI-Key", randomApiKey)
	req.Header.Set("X-RapidAPI-Host", "subdomains-lookup.p.rapidapi.com")

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

	var response rapidapiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	for _, record := range response.Result.Records {
		if record.Domain != "" {
			results = append(results, record.Domain)
		}
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

