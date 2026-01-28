package odin

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

type odinRequest struct {
	Domain string        `json:"domain"`
	Limit  int           `json:"limit"`
	Start  []interface{} `json:"start"`
}

type odinResponse struct {
	Success    bool     `json:"success"`
	Data       []string `json:"data"`
	Pagination struct {
		Last []interface{} `json:"last"`
	} `json:"pagination"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "odin"
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
		return nil, fmt.Errorf("no API keys provided for odin")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	var start []interface{}
	limit := 1000

	for {
		requestBody := odinRequest{
			Domain: query,
			Limit:  limit,
			Start:  start,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return allResults, err
		}

		req, err := http.NewRequest(http.MethodPost, "https://api.odin.io/v1/domain/subdomain/search", bytes.NewReader(jsonData))
		if err != nil {
			return allResults, err
		}

		req.Header.Set("X-API-Key", randomApiKey)
		req.Header.Set("Content-Type", "application/json")

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

		var response odinResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return allResults, err
		}

		if !response.Success || len(response.Data) == 0 {
			break
		}

		allResults = append(allResults, response.Data...)

		if len(response.Pagination.Last) == 0 {
			break
		}
		start = response.Pagination.Last
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

