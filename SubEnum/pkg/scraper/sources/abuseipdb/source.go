package abuseipdb

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Source struct{}

func (s *Source) Name() string {
	return "abuseipdb"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.abuseipdb.com/whois/%s", query), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", "abuseipdb_session=")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract subdomains from <li> tags - pattern matches <li>word...</li>
	re := regexp.MustCompile(`<li>\w.*</li>`)
	matches := re.FindAllString(string(body), -1)

	for _, match := range matches {
		// Remove <li> and </li> tags
		subdomain := regexp.MustCompile(`</?li>`).ReplaceAllString(match, "")
		subdomain = subdomain + "." + query
		results = append(results, subdomain)
	}

	return results, nil
}

