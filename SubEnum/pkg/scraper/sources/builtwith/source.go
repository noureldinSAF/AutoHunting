package builtwith

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Source struct {
	apiKeys []string
}

type response struct {
	Results []resultItem `json:"Results"`
}

type resultItem struct {
	Result result `json:"Result"`
}

type result struct {
	Paths []path `json:"Paths"`
}

type path struct {
	Domain    string `json:"Domain"`
	Url       string `json:"Url"`
	SubDomain string `json:"SubDomain"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "builtwith"
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
	var results []string

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for builtwith")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.builtwith.com/v21/api.json?KEY=%s&HIDETEXT=yes&HIDEDL=yes&NOLIVE=yes&NOMETA=yes&NOPII=yes&NOATTR=yes&LOOKUP=%s", randomApiKey, query), nil)
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

	var data response
	err = jsoniter.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	for _, result := range data.Results {
		for _, path := range result.Result.Paths {
			results = append(results, fmt.Sprintf("%s.%s", path.SubDomain, path.Domain))
		}
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

