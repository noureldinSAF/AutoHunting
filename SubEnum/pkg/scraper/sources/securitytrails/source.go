package securitytrails

import (
	"bytes"
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

type response struct {
	Meta struct {
		ScrollID string `json:"scroll_id"`
	} `json:"meta"`
	Records []struct {
		Hostname string `json:"hostname"`
	} `json:"records"`
	Subdomains []string `json:"subdomains"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "securitytrails"
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
		return nil, fmt.Errorf("no API keys provided for securitytrails")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	headers := map[string]string{"Content-Type": "application/json", "APIKEY": randomApiKey}
	var scrollId string

	for {
		var resp *http.Response
		var err error

		if scrollId == "" {
			requestBody := []byte(fmt.Sprintf(`{"query":"apex_domain='%s'"}`, query))
			req, err := http.NewRequest(http.MethodPost, "https://api.securitytrails.com/v1/domains/list?include_ips=false&scroll=true", bytes.NewReader(requestBody))
			if err != nil {
				return allResults, err
			}

			for k, v := range headers {
				req.Header.Set(k, v)
			}

			resp, err = client.Do(req)
			if err != nil {
				// Fallback to subdomains endpoint if scroll fails with 403
				if resp != nil && resp.StatusCode == 403 {
					req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.securitytrails.com/v1/domain/%s/subdomains", query), nil)
					if err != nil {
						return allResults, err
					}

					for k, v := range headers {
						req.Header.Set(k, v)
					}

					resp, err = client.Do(req)
					if err != nil {
						return allResults, err
					}
				} else {
					return allResults, err
				}
			}
		} else {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.securitytrails.com/v1/scroll/%s", scrollId), nil)
			if err != nil {
				return allResults, err
			}

			for k, v := range headers {
				req.Header.Set(k, v)
			}

			resp, err = client.Do(req)
			if err != nil {
				return allResults, err
			}
		}

		var securityTrailsResponse response
		err = json.NewDecoder(resp.Body).Decode(&securityTrailsResponse)
		resp.Body.Close()
		if err != nil {
			return allResults, err
		}

		for _, record := range securityTrailsResponse.Records {
			allResults = append(allResults, record.Hostname)
		}

		for _, subdomain := range securityTrailsResponse.Subdomains {
			if strings.HasSuffix(subdomain, ".") {
				subdomain += query
			} else {
				subdomain = subdomain + "." + query
			}
			allResults = append(allResults, subdomain)
		}

		scrollId = securityTrailsResponse.Meta.ScrollID
		if scrollId == "" {
			break
		}
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

