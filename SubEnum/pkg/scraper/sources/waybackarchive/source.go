package waybackarchive

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"subenum/pkg/utils"
)

type Source struct{}

func (s *Source) Name() string {
	return "waybackarchive"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var allResults []string
	seen := make(map[string]bool)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=txt&fl=original&collapse=urlkey", query), nil)
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
		return allResults, err
	}

	bodyStr := string(body)
	// Decode URL-encoded content
	bodyStr, _ = url.QueryUnescape(bodyStr)

	// Create subdomain extractor for this domain
	extractor, err := utils.NewSubdomainExtractor(query)
	if err != nil {
		return allResults, err
	}

	// Extract subdomains from the response
	results := extractor.Extract(bodyStr)

	// Deduplicate and clean results
	for _, subdomain := range results {
		// Fix for triple encoded URL artifacts
		cleaned := strings.ToLower(subdomain)
		cleaned = strings.TrimPrefix(cleaned, "25")
		cleaned = strings.TrimPrefix(cleaned, "2f")
		if cleaned != "" && !seen[cleaned] {
			seen[cleaned] = true
			allResults = append(allResults, cleaned)
		}
	}

	return allResults, nil
}

