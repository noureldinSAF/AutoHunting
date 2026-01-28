package threatminer

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Source struct{}

type response struct {
	StatusCode    string   `json:"status_code"`
	StatusMessage string   `json:"status_message"`
	Results       []string `json:"results"`
}

func (s *Source) Name() string {
	return "threatminer"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.threatminer.org/v2/domain.php?q=%s&rt=5", query), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	results = append(results, data.Results...)
	return results, nil
}

