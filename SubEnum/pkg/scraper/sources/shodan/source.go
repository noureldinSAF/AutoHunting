package shodan

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

type dnsdbLookupResponse struct {
	Domain     string   `json:"domain"`
	Subdomains []string `json:"subdomains"`
	Result     int      `json:"result"`
	Error      string   `json:"error"`
	More       bool     `json:"more"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "shodan"
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
		return nil, fmt.Errorf("no API keys provided for shodan")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	page := 1
	for {
		searchURL := fmt.Sprintf("https://api.shodan.io/dns/domain/%s?key=%s&page=%d", query, randomApiKey, page)
		req, err := http.NewRequest(http.MethodGet, searchURL, nil)
		if err != nil {
			return allResults, err
		}

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

		var response dnsdbLookupResponse
		err = jsoniter.Unmarshal(body, &response)
		if err != nil {
			return allResults, err
		}

		if response.Error != "" {
			return allResults, fmt.Errorf("%v", response.Error)
		}

		for _, data := range response.Subdomains {
			value := fmt.Sprintf("%s.%s", data, response.Domain)
			allResults = append(allResults, value)
		}

		if !response.More {
			break
		}
		page++
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

