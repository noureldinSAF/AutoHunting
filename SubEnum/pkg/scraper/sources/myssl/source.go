package myssl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Source struct{}

type myssLResponse struct {
	Data []struct {
		Domain string `json:"domain"`
	} `json:"data"`
}

func (s *Source) Name() string {
	return "myssl"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://myssl.com/api/v1/discover_sub_domain?domain=%s", query), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: status=%d body=%s", resp.StatusCode, string(body))
	}

	var response myssLResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	for _, entry := range response.Data {
		subdomain := entry.Domain
		if subdomain != "" && strings.HasSuffix(subdomain, "."+query) {
			results = append(results, subdomain)
		}
	}

	return results, nil
}

