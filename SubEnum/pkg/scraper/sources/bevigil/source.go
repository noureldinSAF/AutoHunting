package bevigil

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

type Response struct {
	Domain     string   `json:"domain"`
	Subdomains []string `json:"subdomains"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "bevigil"
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
		return nil, fmt.Errorf("no API keys provided for bevigil")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	getUrl := fmt.Sprintf("https://osint.bevigil.com/api/%s/subdomains/", query)
	req, err := http.NewRequest(http.MethodGet, getUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Access-Token", randomApiKey)
	req.Header.Set("User-Agent", "subfinder")

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

	var response Response
	err = jsoniter.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Subdomains) > 0 {
		results = append(results, response.Subdomains...)
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

