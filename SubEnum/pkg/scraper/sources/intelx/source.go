package intelx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Source struct {
	apiKeys []apiKey
}

type apiKey struct {
	host string
	key  string
}

type searchResponseType struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
}

type selectorType struct {
	Selectvalue string `json:"selectorvalue"`
}

type searchResultType struct {
	Selectors []selectorType `json:"selectors"`
	Status    int            `json:"status"`
}

type requestBody struct {
	Term       string `json:"term"`
	Maxresults int    `json:"maxresults"`
	Media      int    `json:"media"`
	Target     int    `json:"target"`
	Terminate  []int  `json:"terminate"`
	Timeout    int    `json:"timeout"`
}

func New(apiKeys []string) *Source {
	keys := make([]apiKey, 0)
	for _, key := range apiKeys {
		parts := strings.Split(key, ":")
		if len(parts) == 2 {
			keys = append(keys, apiKey{host: parts[0], key: parts[1]})
		}
	}
	return &Source{apiKeys: keys}
}

func (s *Source) Name() string {
	return "intelx"
}

func (s *Source) RequiresAPIKey() bool {
	return true
}

func (s *Source) randomKey() apiKey {
	if len(s.apiKeys) == 0 {
		return apiKey{}
	}
	return s.apiKeys[rand.Intn(len(s.apiKeys))]
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var allResults []string

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for intelx")
	}

	randomApiKey := s.randomKey()
	if randomApiKey.host == "" || randomApiKey.key == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	searchURL := fmt.Sprintf("https://%s/phonebook/search?k=%s", randomApiKey.host, randomApiKey.key)
	reqBody := requestBody{
		Term:       query,
		Maxresults: 100000,
		Media:      0,
		Target:     1,
		Timeout:    20,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, searchURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var response searchResponseType
	err = jsoniter.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	resultsURL := fmt.Sprintf("https://%s/phonebook/search/result?k=%s&id=%s&limit=10000", randomApiKey.host, randomApiKey.key, response.ID)
	status := 0
	for status == 0 || status == 3 {
		req, err := http.NewRequest(http.MethodGet, resultsURL, nil)
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

		var response searchResultType
		err = jsoniter.Unmarshal(body, &response)
		if err != nil {
			return allResults, err
		}

		status = response.Status
		for _, hostname := range response.Selectors {
			allResults = append(allResults, hostname.Selectvalue)
		}
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

