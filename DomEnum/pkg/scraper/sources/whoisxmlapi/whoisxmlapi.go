package whoisxmlapi

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
)

const (
	apiUrl = "https://api.whoisxmlapi.com/api/v2"
)

type Source struct{
	apiKeys []string
}

func New(apiKeys []string) *Source {
	return &Source{apiKeys: apiKeys}
}

type request struct {
	APIKey string `json:"apiKey"`
	SearchType string `json:"type"`
	Mode string `json:"mode"`
	PunyCode bool `json:"punycode"`
	BasicSearchTerm *BasicSearchTerm `json:"basicSearchTerm"`
	SearchAfter *string `json:"searchAfter"`
}

type responseObj struct {
	NextPageSearchAfter string `json:"nextPageSearchAfter"`
    DomainsCount int `json:"domainsCount"`
    DomainsList []string `json:"domainsList"`
}

type BasicSearchTerm struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

type response struct {}


func (s *Source) Name() string {
	return "whoisxmlapi"
}

func (s *Source) randomKey() string {
	return s.apiKeys[rand.Intn(len(s.apiKeys))]
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {

	if s.apiKeys == nil || len(s.apiKeys) == 0 {
		return []string{}, nil
	}

	var domains []string
	var searchAfter *string

	for {
		req := &request{
			APIKey:      s.randomKey(),
			SearchType:  "historic",
			Mode:        "purchase",
			PunyCode:    true,
			SearchAfter: searchAfter,
			BasicSearchTerm: &BasicSearchTerm{
				Include: []string{query},
				Exclude: []string{},
			},
		}

		reqBody, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		httpReq, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(httpReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var ro *responseObj
		if err := json.Unmarshal(body, &ro); err != nil {
			return nil, err
		}

		domains = append(domains, ro.DomainsList...)

		if ro.NextPageSearchAfter == "" {
			break
		}
		searchAfter = &ro.NextPageSearchAfter
	}

	return domains, nil
}
	
func (s *Source) RequireAPIKey() bool {
	return true
}

