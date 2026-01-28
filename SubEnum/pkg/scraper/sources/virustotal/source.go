package virustotal

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

type response struct {
	Data []Object `json:"data"`
	Meta Meta     `json:"meta"`
}

type Object struct {
	Id string `json:"id"`
}

type Meta struct {
	Cursor string `json:"cursor"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "virustotal"
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
		return nil, fmt.Errorf("no API keys provided for virustotal")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	var cursor = ""
	for {
		var url = fmt.Sprintf("https://www.virustotal.com/api/v3/domains/%s/subdomains?limit=40", query)
		if cursor != "" {
			url = fmt.Sprintf("%s&cursor=%s", url, cursor)
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return allResults, err
		}

		req.Header.Set("x-apikey", randomApiKey)

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

		var data response
		err = json.Unmarshal(body, &data)
		if err != nil {
			return allResults, err
		}

		for _, subdomain := range data.Data {
			allResults = append(allResults, subdomain.Id)
		}

		cursor = data.Meta.Cursor
		if cursor == "" {
			break
		}
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

