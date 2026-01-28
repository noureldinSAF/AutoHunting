package huntermap

import (
	"encoding/base64"
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

type huntermapResponse struct {
	Data struct {
		List []struct {
			Domain string `json:"domain"`
		} `json:"list"`
		Total int `json:"total"`
	} `json:"data"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "huntermap"
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
		return nil, fmt.Errorf("no API keys provided for huntermap")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	// Calculate time range (last year)
	endTime := time.Now().Format("2006-01-02")
	startTime := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")

	// Base64 encode the domain
	queryEncoded := base64.URLEncoding.EncodeToString([]byte(query))

	page := 1
	pageSize := 100
	totalFound := 0

	for {
		url := fmt.Sprintf("https://api.hunter.how/search?api-key=%s&query=%s&start_time=%s&end_time=%s&page=%d&page_size=%d",
			randomApiKey, queryEncoded, startTime, endTime, page, pageSize)

		req, err := http.NewRequest(http.MethodGet, url, nil)
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

		var response huntermapResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return allResults, err
		}

		if len(response.Data.List) == 0 {
			break
		}

		for _, item := range response.Data.List {
			subdomain := item.Domain
			if strings.HasSuffix(subdomain, "."+query) {
				allResults = append(allResults, subdomain)
				totalFound++
			}
		}

		if totalFound >= response.Data.Total {
			break
		}

		page++
		time.Sleep(2 * time.Second) // Rate limiting
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

