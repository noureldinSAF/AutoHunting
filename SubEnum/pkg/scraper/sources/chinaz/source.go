package chinaz

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

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "chinaz"
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
		return nil, fmt.Errorf("no API keys provided for chinaz")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://apidatav2.chinaz.com/single/alexa?key=%s&domain=%s", randomApiKey, query), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: status=%d body=%s", resp.StatusCode, string(body))
	}

	// Parse JSON using jsoniter-like approach with standard library
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return results, nil
	}

	subdomainList, ok := result["ContributingSubdomainList"].([]interface{})
	if !ok {
		return results, nil
	}

	for _, item := range subdomainList {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		dataUrl, ok := itemMap["DataUrl"].(string)
		if ok && dataUrl != "" {
			results = append(results, dataUrl)
		}
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

