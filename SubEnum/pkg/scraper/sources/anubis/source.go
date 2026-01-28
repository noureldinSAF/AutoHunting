package anubis

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Source struct{}

func (s *Source) Name() string {
	return "anubis"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://jonlu.ca/anubis/subdomains/%s", query), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return results, nil
	}

	var subdomains []string
	err = json.NewDecoder(resp.Body).Decode(&subdomains)
	if err != nil {
		return nil, err
	}

	results = append(results, subdomains...)
	return results, nil
}
