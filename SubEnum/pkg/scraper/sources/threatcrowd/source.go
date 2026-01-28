package threatcrowd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Source struct{}

type threatCrowdResponse struct {
	ResponseCode string   `json:"response_code"`
	Subdomains   []string `json:"subdomains"`
	Undercount   string   `json:"undercount"`
}

func (s *Source) Name() string {
	return "threatcrowd"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	url := fmt.Sprintf("http://ci-www.threatcrowd.org/searchApi/v2/domain/report/?domain=%s", query)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tcResponse threatCrowdResponse
	if err := json.Unmarshal(body, &tcResponse); err != nil {
		return nil, err
	}

	for _, subdomain := range tcResponse.Subdomains {
		if subdomain != "" {
			results = append(results, subdomain)
		}
	}

	return results, nil
}
