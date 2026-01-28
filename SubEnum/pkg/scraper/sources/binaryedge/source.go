package binaryedge

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

type binaryedgeResponse struct {
	Events []string `json:"events"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "binaryedge"
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
		return nil, fmt.Errorf("no API keys provided for binaryedge")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	pageNum := 1
	pageSize := 100

	for {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.binaryedge.io/v2/query/domains/subdomain/%s?page=%d&pagesize=%d", query, pageNum, pageSize), nil)
		if err != nil {
			return allResults, err
		}

		req.Header.Set("X-Key", randomApiKey)

		resp, err := client.Do(req)
		if err != nil {
			return allResults, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return allResults, err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return allResults, fmt.Errorf("api error: status=%d body=%s", resp.StatusCode, string(body))
		}

		var response binaryedgeResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return allResults, err
		}

		if len(response.Events) == 0 {
			break
		}

		allResults = append(allResults, response.Events...)
		pageNum++
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

