package pugrecon

import (
	"bytes"
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

type pugreconResult struct {
	Name string `json:"name"`
}

type pugreconAPIResponse struct {
	Results        []pugreconResult `json:"results"`
	QuotaRemaining int              `json:"quota_remaining"`
	Limited        bool             `json:"limited"`
	TotalResults   int              `json:"total_results"`
	Message        string           `json:"message"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "pugrecon"
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
		return nil, fmt.Errorf("no API keys provided for pugrecon")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	postData := map[string]string{"domain_name": query}
	bodyBytes, err := json.Marshal(postData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://pugrecon.com/api/v1/domains", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+randomApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var apiResp pugreconAPIResponse
		if json.Unmarshal(body, &apiResp) == nil && apiResp.Message != "" {
			return nil, fmt.Errorf("received status code %d: %s", resp.StatusCode, apiResp.Message)
		}
		return nil, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	var response pugreconAPIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	for _, subdomain := range response.Results {
		results = append(results, subdomain.Name)
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

