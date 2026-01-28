package quake

import (
	"bytes"
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

type quakeResults struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Service struct {
			HTTP struct {
				Host string `json:"host"`
			} `json:"http"`
		}
	} `json:"data"`
	Meta struct {
		Pagination struct {
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "quake"
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
		return nil, fmt.Errorf("no API keys provided for quake")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	var pageSize = 500
	var start = 0
	var totalResults = -1

	for {
		requestBody := []byte(fmt.Sprintf(`{"query":"domain: %s", "include":["service.http.host"], "latest": true, "size":%d, "start":%d}`, query, pageSize, start))
		
		req, err := http.NewRequest(http.MethodPost, "https://quake.360.net/api/v3/search/quake_service", bytes.NewReader(requestBody))
		if err != nil {
			return allResults, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-QuakeToken", randomApiKey)

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

		var response quakeResults
		err = json.Unmarshal(body, &response)
		if err != nil {
			return allResults, err
		}

		if response.Code != 0 {
			return allResults, fmt.Errorf("%s", response.Message)
		}

		if totalResults == -1 {
			totalResults = response.Meta.Pagination.Total
		}

		for _, quakeDomain := range response.Data {
			subdomain := quakeDomain.Service.HTTP.Host
			if strings.ContainsAny(subdomain, "暂无权限") {
				continue
			}
			allResults = append(allResults, subdomain)
		}

		if len(response.Data) == 0 || start+pageSize >= totalResults {
			break
		}

		start += pageSize
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

