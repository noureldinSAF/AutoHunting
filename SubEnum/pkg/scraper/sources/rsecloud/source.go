package rsecloud

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
	Count      int      `json:"count"`
	Data       []string `json:"data"`
	Page       int      `json:"page"`
	PageSize   int      `json:"pagesize"`
	TotalPages int      `json:"total_pages"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "rsecloud"
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
	var allResults []string

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for rsecloud")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	headers := map[string]string{"Content-Type": "application/json", "X-API-Key": randomApiKey}

	fetchSubdomains := func(endpoint string) ([]string, error) {
		var results []string
		page := 1
		for {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.rsecloud.com/api/v2/subdomains/%s/%s?page=%d", endpoint, query, page), nil)
			if err != nil {
				return results, err
			}

			for k, v := range headers {
				req.Header.Set(k, v)
			}

			resp, err := client.Do(req)
			if err != nil {
				return results, err
			}

			var rseCloudResponse response
			err = json.NewDecoder(resp.Body).Decode(&rseCloudResponse)
			resp.Body.Close()
			if err != nil {
				return results, err
			}

			results = append(results, rseCloudResponse.Data...)

			if page >= rseCloudResponse.TotalPages {
				break
			}
			page++
		}
		return results, nil
	}

	active, err := fetchSubdomains("active")
	if err != nil {
		return allResults, err
	}
	allResults = append(allResults, active...)

	passive, err := fetchSubdomains("passive")
	if err != nil {
		return allResults, err
	}
	allResults = append(allResults, passive...)

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

