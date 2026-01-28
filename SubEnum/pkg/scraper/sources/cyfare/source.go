package cyfare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Source struct{}

type cyfareRequest struct {
	Domain string `json:"domain"`
}

type cyfareResponse struct {
	Subdomains []string `json:"subdomains"`
}

func (s *Source) Name() string {
	return "cyfare"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	requestBody := cyfareRequest{Domain: query}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://cyfare.net/apps/VulnerabilityStudio/subfind/query.php", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Origin", "https://cyfare.net")
	req.Header.Set("Content-Type", "application/json")

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

	var response cyfareResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	results = append(results, response.Subdomains...)
	return results, nil
}

