package racent

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Source struct{}

type racentResponse struct {
	Data struct {
		List []struct {
			DNSNames []string `json:"dnsnames"`
		} `json:"list"`
	} `json:"data"`
}

func (s *Source) Name() string {
	return "racent"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://face.racent.com/tool/query_ctlog?keyword=%s", query), nil)
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

	var response racentResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	for _, item := range response.Data.List {
		results = append(results, item.DNSNames...)
	}

	return results, nil
}

