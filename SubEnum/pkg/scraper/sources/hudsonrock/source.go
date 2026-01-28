package hudsonrock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Source struct{}

type hudsonrockResponse struct {
	Data struct {
		EmployeesUrls []struct {
			URL string `json:"url"`
		} `json:"employees_urls"`
		ClientsUrls []struct {
			URL string `json:"url"`
		} `json:"clients_urls"`
	} `json:"data"`
}

func (s *Source) Name() string {
	return "hudsonrock"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://cavalier.hudsonrock.com/api/json/v2/osint-tools/urls-by-domain?domain=%s", query), nil)
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

	var response hudsonrockResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	allUrls := append(response.Data.EmployeesUrls, response.Data.ClientsUrls...)
	for _, record := range allUrls {
		// Extract domain from URL
		url := record.URL
		if strings.HasPrefix(url, "http://") {
			url = strings.TrimPrefix(url, "http://")
		}
		if strings.HasPrefix(url, "https://") {
			url = strings.TrimPrefix(url, "https://")
		}
		if strings.Contains(url, "/") {
			url = strings.Split(url, "/")[0]
		}
		if url != "" {
			results = append(results, url)
		}
	}

	return results, nil
}

