package merklemap

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

type merklemapResponse struct {
	Results []struct {
		Hostname          string `json:"hostname"`
		SubjectCommonName string `json:"subject_common_name"`
	} `json:"results"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "merklemap"
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
		return nil, fmt.Errorf("no API keys provided for merklemap")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	page := 0
	for {
		url := fmt.Sprintf("https://api.merklemap.com/v1/search?query=*.%s&page=%d&type=wildcard", query, page)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return allResults, err
		}

		req.Header.Set("Authorization", "Bearer "+randomApiKey)

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

		var response merklemapResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return allResults, err
		}

		if len(response.Results) == 0 {
			break
		}

		for _, result := range response.Results {
			if result.Hostname != "" && strings.HasSuffix(result.Hostname, "."+query) {
				allResults = append(allResults, result.Hostname)
			}
			if result.SubjectCommonName != "" && strings.HasSuffix(result.SubjectCommonName, "."+query) {
				allResults = append(allResults, result.SubjectCommonName)
			}
		}

		page++
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

