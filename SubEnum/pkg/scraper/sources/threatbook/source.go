package threatbook

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Source struct {
	apiKeys []string
}

type threatBookResponse struct {
	ResponseCode int64  `json:"response_code"`
	VerboseMsg   string `json:"verbose_msg"`
	Data         struct {
		Domain     string `json:"domain"`
		SubDomains struct {
			Total string   `json:"total"`
			Data  []string `json:"data"`
		} `json:"sub_domains"`
	} `json:"data"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "threatbook"
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
		return nil, fmt.Errorf("no API keys provided for threatbook")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.threatbook.cn/v3/domain/sub_domains?apikey=%s&resource=%s", randomApiKey, query), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response threatBookResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	if response.ResponseCode != 0 {
		return nil, fmt.Errorf("code %d, %s", response.ResponseCode, response.VerboseMsg)
	}

	total, err := strconv.ParseInt(response.Data.SubDomains.Total, 10, 64)
	if err != nil {
		return nil, err
	}

	if total > 0 {
		results = append(results, response.Data.SubDomains.Data...)
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

