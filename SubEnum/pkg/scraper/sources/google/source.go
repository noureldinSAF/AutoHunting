package google

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Source struct {
	apiKeys []string
	cxKeys  []string
}

type googleResponse struct {
	Items []struct {
		DisplayLink string `json:"displayLink"`
	} `json:"items"`
}

func New(apiKeys []string) *Source {
	// Google needs both API key and CX (Custom Search Engine ID)
	// Format: "apikey:cx"
	rapidKeys := make([]string, 0)
	cxKeys := make([]string, 0)

	for _, key := range apiKeys {
		parts := strings.Split(key, ":")
		if len(parts) == 2 {
			rapidKeys = append(rapidKeys, parts[0])
			cxKeys = append(cxKeys, parts[1])
		}
	}

	return &Source{
		apiKeys: rapidKeys,
		cxKeys:  cxKeys,
	}
}

func (s *Source) Name() string {
	return "google"
}

func (s *Source) RequiresAPIKey() bool {
	return true
}

func (s *Source) randomKeys() (string, string) {
	if len(s.apiKeys) == 0 || len(s.cxKeys) == 0 {
		return "", ""
	}
	apiIdx := rand.Intn(len(s.apiKeys))
	cxIdx := rand.Intn(len(s.cxKeys))
	return s.apiKeys[apiIdx], s.cxKeys[cxIdx]
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var allResults []string
	seen := make(map[string]bool)

	if len(s.apiKeys) == 0 || len(s.cxKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for google (needs both API key and CX)")
	}

	randomApiKey, randomCX := s.randomKeys()
	if randomApiKey == "" || randomCX == "" {
		return nil, fmt.Errorf("no valid API keys available")
	}

	dork := fmt.Sprintf("site:*.%s -www", query)
	page := 1

	for page <= 100 { // Google CSE limits to 100 results
		apiURL := fmt.Sprintf("https://customsearch.googleapis.com/customsearch/v1?q=%s&cx=%s&num=10&start=%d&key=%s&alt=json",
			url.QueryEscape(dork), randomCX, page, randomApiKey)

		req, err := http.NewRequest(http.MethodGet, apiURL, nil)
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

		var response googleResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return allResults, err
		}

		if len(response.Items) == 0 {
			break
		}

		for _, item := range response.Items {
			if item.DisplayLink != "" && !seen[item.DisplayLink] {
				seen[item.DisplayLink] = true
				allResults = append(allResults, item.DisplayLink)
			}
		}

		page += 10
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

