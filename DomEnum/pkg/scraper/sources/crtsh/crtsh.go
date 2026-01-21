package crtsh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cyinnove/tldify"
)

type Source struct{}

func New() *Source {
	return &Source{}
}

type response struct {
	CommonName string `json:"common_name"`
}

func (s *Source) Name() string {
	return "crtsh"
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	// 1. Send HTTP Request
	fmtURL := fmt.Sprintf("https://crt.sh/?q=%s&output=json", query)
	req, err := http.NewRequest("GET", fmtURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 2. Parse Response
	var r []*response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	// print(len(r))

	// 3. Normalize and collect ALL results
	for _, cn := range r {
		cn.CommonName = normalizeDomain(cn.CommonName)
		if cn.CommonName != "" {
			results = append(results, cn.CommonName) // Keep adding to the list
		}
	}

	// 4. Return the full list after the loop finishes
	return results, nil
}

func (s *Source) RequireAPIKey() bool {
	return false
}

func normalizeDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	domain = strings.Trim(domain, "*.")

	parsedDomain, err := tldify.Parse(domain)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(parsedDomain.Domain + "." + parsedDomain.TLD)

}