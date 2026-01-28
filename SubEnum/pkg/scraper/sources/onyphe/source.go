package onyphe

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Source struct {
	apiKeys []string
}

type OnypheResponse struct {
	Error    int      `json:"error"`
	Results  []Result `json:"results"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
	Total    int      `json:"total"`
	MaxPage  int      `json:"max_page"`
}

type Result struct {
	Subdomains []string `json:"subdomains"`
	Hostname   string   `json:"hostname"`
	Forward    string   `json:"forward"`
	Reverse    string   `json:"reverse"`
	Host       string   `json:"host"`
	Domain     string   `json:"domain"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "onyphe"
}

func (s *Source) RequiresAPIKey() bool {
	return true
}

func (s *Source) randomKey() string {
	if len(s.apiKeys) == 0 {
		return ""
	}
	return s.apiKeys[rand.Intn(len(s.apiKeys))]
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var allResults []string
	seen := make(map[string]bool)

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for onyphe")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	headers := map[string]string{"Content-Type": "application/json", "Authorization": "bearer " + randomApiKey}

	page := 1
	pageSize := 1000

	for {
		urlWithQuery := fmt.Sprintf("https://www.onyphe.io/api/v2/search/?q=%s&page=%d&size=%d",
			url.QueryEscape("category:resolver domain:"+query), page, pageSize)

		req, err := http.NewRequest(http.MethodGet, urlWithQuery, nil)
		if err != nil {
			return allResults, err
		}

		for k, v := range headers {
			req.Header.Set(k, v)
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

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return allResults, fmt.Errorf("api error: status=%d body=%s", resp.StatusCode, string(body))
		}

		// Handle flexible JSON parsing
		var respOnyphe OnypheResponse
		err = s.unmarshalOnypheResponse(body, &respOnyphe)
		if err != nil {
			return allResults, err
		}

		for _, record := range respOnyphe.Results {
			for _, subdomain := range record.Subdomains {
				if subdomain != "" && !seen[subdomain] {
					seen[subdomain] = true
					allResults = append(allResults, subdomain)
				}
			}

			if record.Hostname != "" && !seen[record.Hostname] {
				seen[record.Hostname] = true
				allResults = append(allResults, record.Hostname)
			}

			if record.Forward != "" && !seen[record.Forward] {
				seen[record.Forward] = true
				allResults = append(allResults, record.Forward)
			}

			if record.Reverse != "" && !seen[record.Reverse] {
				seen[record.Reverse] = true
				allResults = append(allResults, record.Reverse)
			}
		}

		if len(respOnyphe.Results) == 0 || page >= respOnyphe.MaxPage {
			break
		}

		page++
	}

	return allResults, nil
}

func (s *Source) unmarshalOnypheResponse(data []byte, resp *OnypheResponse) error {
	var raw struct {
		Error    int             `json:"error"`
		Results  json.RawMessage `json:"results"`
		Page     json.RawMessage `json:"page"`
		PageSize json.RawMessage `json:"page_size"`
		Total    json.RawMessage `json:"total"`
		MaxPage  json.RawMessage `json:"max_page"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	resp.Error = raw.Error

	// Parse Results
	var results []Result
	if err := json.Unmarshal(raw.Results, &results); err == nil {
		resp.Results = results
	}

	// Parse Page
	if pageStr := string(raw.Page); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			resp.Page = page
		} else {
			var pageStrQuoted string
			if err := json.Unmarshal(raw.Page, &pageStrQuoted); err == nil {
				if page, err := strconv.Atoi(pageStrQuoted); err == nil {
					resp.Page = page
				}
			}
		}
	}

	// Parse PageSize
	if pageSizeStr := string(raw.PageSize); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			resp.PageSize = pageSize
		} else {
			var pageSizeStrQuoted string
			if err := json.Unmarshal(raw.PageSize, &pageSizeStrQuoted); err == nil {
				if pageSize, err := strconv.Atoi(pageSizeStrQuoted); err == nil {
					resp.PageSize = pageSize
				}
			}
		}
	}

	// Parse Total
	if totalStr := string(raw.Total); totalStr != "" {
		if total, err := strconv.Atoi(totalStr); err == nil {
			resp.Total = total
		} else {
			var totalStrQuoted string
			if err := json.Unmarshal(raw.Total, &totalStrQuoted); err == nil {
				if total, err := strconv.Atoi(totalStrQuoted); err == nil {
					resp.Total = total
				}
			}
		}
	}

	// Parse MaxPage
	if maxPageStr := string(raw.MaxPage); maxPageStr != "" {
		if maxPage, err := strconv.Atoi(maxPageStr); err == nil {
			resp.MaxPage = maxPage
		} else {
			var maxPageStrQuoted string
			if err := json.Unmarshal(raw.MaxPage, &maxPageStrQuoted); err == nil {
				if maxPage, err := strconv.Atoi(maxPageStrQuoted); err == nil {
					resp.MaxPage = maxPage
				}
			}
		}
	}

	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

