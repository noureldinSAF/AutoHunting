package redhuntlabs

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Source struct {
	apiKeys []string
}

type Response struct {
	Subdomains []string         `json:"subdomains"`
	Metadata   ResponseMetadata `json:"metadata"`
}

type ResponseMetadata struct {
	ResultCount int `json:"result_count"`
	PageSize    int `json:"page_size"`
	PageNumber  int `json:"page_number"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "redhuntlabs"
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
		return nil, fmt.Errorf("no API keys provided for redhuntlabs")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" || !strings.Contains(randomApiKey, ":") {
		return nil, fmt.Errorf("invalid API key format, expected baseUrl:key")
	}

	randomApiInfo := strings.Split(randomApiKey, ":")
	if len(randomApiInfo) != 3 {
		return nil, fmt.Errorf("invalid API key format, expected protocol:host:key")
	}

	baseUrl := randomApiInfo[0] + ":" + randomApiInfo[1]
	requestHeaders := map[string]string{"X-BLOBR-KEY": randomApiInfo[2], "User-Agent": "subfinder"}
	pageSize := 1000

	getUrl := fmt.Sprintf("%s?domain=%s&page=1&page_size=%d", baseUrl, query, pageSize)
	req, err := http.NewRequest(http.MethodGet, getUrl, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range requestHeaders {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("encountered error: %v; note: if you get a 'limit has been reached' error, head over to https://devportal.redhuntlabs.com", err)
	}
	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	allResults = append(allResults, response.Subdomains...)

	if response.Metadata.ResultCount > pageSize {
		totalPages := (response.Metadata.ResultCount + pageSize - 1) / pageSize
		for page := 2; page <= totalPages; page++ {
			getUrl = fmt.Sprintf("%s?domain=%s&page=%d&page_size=%d", baseUrl, query, page, pageSize)
			req, err := http.NewRequest(http.MethodGet, getUrl, nil)
			if err != nil {
				return allResults, err
			}

			for k, v := range requestHeaders {
				req.Header.Set(k, v)
			}

			resp, err := client.Do(req)
			if err != nil {
				return allResults, fmt.Errorf("encountered error: %v; note: if you get a 'limit has been reached' error, head over to https://devportal.redhuntlabs.com", err)
			}

			err = json.NewDecoder(resp.Body).Decode(&response)
			resp.Body.Close()
			if err != nil {
				continue
			}

			allResults = append(allResults, response.Subdomains...)
		}
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

