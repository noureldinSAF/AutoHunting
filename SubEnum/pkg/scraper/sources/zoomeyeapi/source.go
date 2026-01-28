package zoomeyeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Source struct {
	apiKeys []string
}

type zoomeyeResults struct {
	Status int `json:"status"`
	Total  int `json:"total"`
	List   []struct {
		Name string   `json:"name"`
		Ip   []string `json:"ip"`
	} `json:"list"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "zoomeyeapi"
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
		return nil, fmt.Errorf("no API keys provided for zoomeyeapi")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	randomApiInfo := strings.Split(randomApiKey, ":")
	if len(randomApiInfo) != 2 {
		return nil, fmt.Errorf("invalid API key format, expected host:apikey")
	}

	host := randomApiInfo[0]
	apiKey := randomApiInfo[1]

	headers := map[string]string{
		"API-KEY":      apiKey,
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var pages = 1
	for currentPage := 1; currentPage <= pages; currentPage++ {
		api := fmt.Sprintf("https://api.%s/domain/search?q=%s&type=1&s=1000&page=%d", host, query, currentPage)
		
		req, err := http.NewRequest(http.MethodGet, api, nil)
		if err != nil {
			return allResults, err
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		isForbidden := resp != nil && resp.StatusCode == http.StatusForbidden
		if err != nil {
			if !isForbidden {
				return allResults, err
			}
			return allResults, nil
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return allResults, err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return allResults, fmt.Errorf("api error: status=%d body=%s", resp.StatusCode, string(body))
		}

		var res zoomeyeResults
		err = json.Unmarshal(body, &res)
		if err != nil {
			return allResults, err
		}

		pages = int(res.Total/1000) + 1
		for _, r := range res.List {
			allResults = append(allResults, r.Name)
		}
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

