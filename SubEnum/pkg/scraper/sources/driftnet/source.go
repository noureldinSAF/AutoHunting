package driftnet

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Source struct {
	apiKeys []string
}

const (
	baseURL      = "https://api.driftnet.io/v1/"
	summaryLimit = 10000
)

type endpointConfig struct {
	endpoint string
	param    string
	context  string
}

var endpoints = []endpointConfig{
	{"ct/log", "field=host:", "cert-dns-name"},
	{"scan/protocols", "field=host:", "cert-dns-name"},
	{"scan/domains", "field=host:", "cert-dns-name"},
	{"domain/rdns", "host=", "dns-ptr"},
}

type summaryResponse struct {
	Summary struct {
		Other  int            `json:"other"`
		Values map[string]int `json:"values"`
	} `json:"summary"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "driftnet"
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
	seen := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for driftnet")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	headers := map[string]string{
		"accept": "application/json",
	}
	if randomApiKey != "" {
		headers["authorization"] = "Bearer " + randomApiKey
	}

	for _, epConfig := range endpoints {
		wg.Add(1)
		go func(ep endpointConfig) {
			defer wg.Done()
			results := s.queryEndpoint(client, query, headers, ep)
			mu.Lock()
			for _, result := range results {
				if !seen[result] {
					seen[result] = true
					allResults = append(allResults, result)
				}
			}
			mu.Unlock()
		}(epConfig)
	}

	wg.Wait()
	return allResults, nil
}

func (s *Source) queryEndpoint(client *http.Client, domain string, headers map[string]string, epConfig endpointConfig) []string {
	var results []string

	requestURL := fmt.Sprintf("%s%s?%s%s&summarize=host&summary_context=%s&summary_limit=%d",
		baseURL, epConfig.endpoint, epConfig.param, url.QueryEscape(domain), epConfig.context, summaryLimit)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return results
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return results
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return results
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return results
	}

	var summary summaryResponse
	err = json.Unmarshal(body, &summary)
	if err != nil {
		return results
	}

	for subdomain := range summary.Summary.Values {
		if strings.HasSuffix(subdomain, "."+domain) {
			results = append(results, subdomain)
		}
	}

	return results
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

