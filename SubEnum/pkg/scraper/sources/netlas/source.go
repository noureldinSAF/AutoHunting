package netlas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type Source struct {
	apiKeys []string
}

type Item struct {
	Data struct {
		Domain string `json:"domain"`
	} `json:"data"`
}

type DomainsCountResponse struct {
	Count int `json:"count"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "netlas"
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
		return nil, fmt.Errorf("no API keys provided for netlas")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	// To get count of domains
	endpoint := "https://app.netlas.io/api/domains_count/"
	params := url.Values{}
	countQuery := fmt.Sprintf("domain:*.%s AND NOT domain:%s", query, query)
	params.Set("q", countQuery)
	countUrl := endpoint + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, countUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("X-API-Key", randomApiKey)

	resp1, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp1.StatusCode != 200 {
		resp1.Body.Close()
		return nil, fmt.Errorf("request rate limited with status code %d", resp1.StatusCode)
	}

	body, err := io.ReadAll(resp1.Body)
	resp1.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error reading response body")
	}

	var domainsCount DomainsCountResponse
	err = json.Unmarshal(body, &domainsCount)
	if err != nil {
		return nil, err
	}

	// Make a single POST request to get all domains via download method
	apiUrl := "https://app.netlas.io/api/domains/download/"
	requestBody := map[string]interface{}{
		"q":           countQuery,
		"fields":      []string{"*"},
		"source_type": "include",
		"size":        domainsCount.Count,
	}
	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body")
	}

	req2, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewReader(jsonRequestBody))
	if err != nil {
		return nil, err
	}

	req2.Header.Set("accept", "application/json")
	req2.Header.Set("X-API-Key", randomApiKey)
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := client.Do(req2)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body")
	}

	if resp2.StatusCode == 429 {
		return nil, fmt.Errorf("request rate limited with status code %d", resp2.StatusCode)
	}

	if resp2.StatusCode < 200 || resp2.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: status=%d body=%s", resp2.StatusCode, string(body2))
	}

	var data []Item
	err = json.Unmarshal(body2, &data)
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		if item.Data.Domain != "" {
			results = append(results, item.Data.Domain)
		}
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

