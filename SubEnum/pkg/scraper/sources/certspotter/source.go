package certspotter

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

type certspotterObject struct {
	ID       string   `json:"id"`
	DNSNames []string `json:"dns_names"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "certspotter"
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
		return nil, fmt.Errorf("no API keys provided for certspotter")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	headers := map[string]string{"Authorization": "Bearer " + randomApiKey}

	// First request
	reqURL := fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&include_subdomains=true&expand=dns_names", query)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var response []certspotterObject
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: status=%d body=%s", resp.StatusCode, string(body))
	}

	err = jsoniter.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	for _, cert := range response {
		allResults = append(allResults, cert.DNSNames...)
	}

	// If no results, return early
	if len(response) == 0 {
		return allResults, nil
	}

	// Pagination
	id := response[len(response)-1].ID
	for {
		reqURL := fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&include_subdomains=true&expand=dns_names&after=%s", query, id)

		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			return allResults, err
		}

		for k, v := range headers {
			req.Header.Set(k, v)
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

		var response []certspotterObject
		err = jsoniter.Unmarshal(body, &response)
		if err != nil {
			return allResults, err
		}

		if len(response) == 0 {
			break
		}

		for _, cert := range response {
			allResults = append(allResults, cert.DNSNames...)
		}

		id = response[len(response)-1].ID
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

