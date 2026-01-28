package windvane

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Source struct {
	apiKeys []string
}

type response struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data responseData `json:"data"`
}

type responseData struct {
	List         []domainEntry `json:"list"`
	PageResponse pageInfo      `json:"page_response"`
}

type domainEntry struct {
	Domain string `json:"domain"`
}

type pageInfo struct {
	Total     string `json:"total"`
	Count     string `json:"count"`
	TotalPage string `json:"total_page"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "windvane"
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

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for windvane")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	headers := map[string]string{"Content-Type": "application/json", "X-Api-Key": randomApiKey}

	page := 1
	count := 1000
	for {
		requestBody, _ := json.Marshal(map[string]interface{}{"domain": query, "page_request": map[string]int{"page": page, "count": count}})
		
		req, err := http.NewRequest(http.MethodPost, "https://windvane.lichoin.com/trpc.backendhub.public.WindvaneService/ListSubDomain", bytes.NewReader(requestBody))
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

		var windvaneResponse response
		err = json.Unmarshal(body, &windvaneResponse)
		if err != nil {
			return allResults, err
		}

		for _, record := range windvaneResponse.Data.List {
			allResults = append(allResults, record.Domain)
		}

		pageInfo := windvaneResponse.Data.PageResponse
		var totalRecords, recordsPerPage int

		if totalRecords, err = strconv.Atoi(pageInfo.Total); err != nil {
			break
		}
		if recordsPerPage, err = strconv.Atoi(pageInfo.Count); err != nil {
			break
		}

		if (page-1)*recordsPerPage >= totalRecords {
			break
		}

		page++
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

