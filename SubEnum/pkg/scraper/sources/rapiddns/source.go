package rapiddns

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"subenum/pkg/utils"
)

type Source struct{}

var pagePattern = regexp.MustCompile(`class="page-link" href="/subdomain/[^"]+\?page=(\d+)">`)

func (s *Source) Name() string {
	return "rapiddns"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var allResults []string
	seen := make(map[string]bool)

	// Create subdomain extractor for this domain
	extractor, err := utils.NewSubdomainExtractor(query)
	if err != nil {
		return allResults, err
	}

	page := 1
	maxPages := 1
	for {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://rapiddns.io/subdomain/%s?page=%d&full=1", query, page), nil)
		if err != nil {
			return allResults, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return allResults, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return allResults, err
		}

		src := string(body)
		
		// Extract subdomains from the HTML
		extracted := extractor.Extract(src)
		for _, subdomain := range extracted {
			cleaned := strings.TrimSpace(subdomain)
			if cleaned != "" && !seen[cleaned] {
				seen[cleaned] = true
				allResults = append(allResults, cleaned)
			}
		}

		if maxPages == 1 {
			matches := pagePattern.FindAllStringSubmatch(src, -1)
			if len(matches) > 0 {
				lastMatch := matches[len(matches)-1]
				if len(lastMatch) > 1 {
					maxPages, _ = strconv.Atoi(lastMatch[1])
				}
			}
		}

		if page >= maxPages {
			break
		}
		page++
	}

	return allResults, nil
}

