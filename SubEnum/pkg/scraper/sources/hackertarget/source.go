package hackertarget

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

type Source struct{}

func (s *Source) Name() string {
	return "hackertarget"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", query), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		// Extract subdomain from line (format: subdomain,ip)
		parts := strings.Split(line, ",")
		if len(parts) > 0 && parts[0] != "" {
			results = append(results, parts[0])
		}
	}

	return results, nil
}

