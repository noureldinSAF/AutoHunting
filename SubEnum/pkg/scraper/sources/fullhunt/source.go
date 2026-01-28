package fullhunt

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

type fullHuntResponse struct {
	Hosts   []string `json:"hosts"`
	Message string   `json:"message"`
	Status  int      `json:"status"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "fullhunt"
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
		return nil, fmt.Errorf("no API keys provided for fullhunt")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://fullhunt.io/api/v1/domain/%s/subdomains", query), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", randomApiKey)

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

	var response fullHuntResponse
	err = jsoniter.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	results = append(results, response.Hosts...)
	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

