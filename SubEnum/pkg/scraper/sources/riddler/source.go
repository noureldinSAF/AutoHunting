package riddler

import (
	"fmt"
	"io"
	"net/http"

	"subenum/pkg/utils"
)

type Source struct{}

func (s *Source) Name() string {
	return "riddler"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://riddler.io/search?q=pld:%s&view_type=data_table", query), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read entire body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	// Create subdomain extractor for this domain
	extractor, err := utils.NewSubdomainExtractor(query)
	if err != nil {
		return []string{}, err
	}

	// Extract subdomains from the response
	results := extractor.Extract(string(body))

	return results, nil
}

