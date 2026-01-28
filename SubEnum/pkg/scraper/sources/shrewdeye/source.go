package shrewdeye

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Source struct{}

func (s *Source) Name() string {
	return "shrewdeye"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://shrewdeye.app/domains/%s.txt", query), nil)
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

	// Parse plain text response - one subdomain per line
	data := strings.TrimSpace(string(body))
	if data == "" {
		return results, nil
	}

	lines := strings.Split(data, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && (strings.HasSuffix(line, "."+query) || line == query) {
			results = append(results, line)
		}
	}

	return results, nil
}

