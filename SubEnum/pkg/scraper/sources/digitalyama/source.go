package digitalyama

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Source struct {
	apiKeys []string
}

type digitalYamaResponse struct {
	Query        string   `json:"query"`
	Count        int      `json:"count"`
	Subdomains   []string `json:"subdomains"`
	UsageSummary struct {
		QueryCost        float64 `json:"query_cost"`
		CreditsRemaining float64 `json:"credits_remaining"`
	} `json:"usage_summary"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "digitalyama"
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
		return nil, fmt.Errorf("no API keys provided for digitalyama")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	searchURL := fmt.Sprintf("https://api.digitalyama.com/subdomain_finder?domain=%s", query)
	req, err := http.NewRequest(http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", randomApiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var errResponse struct {
			Detail []struct {
				Loc  []string `json:"loc"`
				Msg  string   `json:"msg"`
				Type string   `json:"type"`
			} `json:"detail"`
		}
		err = jsoniter.Unmarshal(body, &errResponse)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
		}
		if len(errResponse.Detail) > 0 {
			return nil, fmt.Errorf("%s (code %d)", errResponse.Detail[0].Msg, resp.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var response digitalYamaResponse
	err = jsoniter.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	results = append(results, response.Subdomains...)
	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

