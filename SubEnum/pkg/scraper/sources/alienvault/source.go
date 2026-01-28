package alienvault

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

type alienvaultResponse struct {
	Detail     string `json:"detail"`
	Error      string `json:"error"`
	PassiveDNS []struct {
		Hostname string `json:"hostname"`
	} `json:"passive_dns"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "alienvault"
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
		return nil, fmt.Errorf("no API keys provided for alienvault")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/passive_dns", query), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+randomApiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response alienvaultResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, fmt.Errorf("%s, %s", response.Detail, response.Error)
	}

	for _, record := range response.PassiveDNS {
		results = append(results, record.Hostname)
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

