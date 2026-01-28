package fofa

import (
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Source struct {
	apiKeys []apiKey
}

type apiKey struct {
	username string
	secret   string
}

type fofaResponse struct {
	Error   bool     `json:"error"`
	ErrMsg  string   `json:"errmsg"`
	Size    int      `json:"size"`
	Results []string `json:"results"`
}

func New(apiKeys []string) *Source {
	keys := make([]apiKey, 0)
	for _, key := range apiKeys {
		parts := strings.Split(key, ":")
		if len(parts) == 2 {
			keys = append(keys, apiKey{username: parts[0], secret: parts[1]})
		}
	}
	return &Source{apiKeys: keys}
}

func (s *Source) Name() string {
	return "fofa"
}

func (s *Source) RequiresAPIKey() bool {
	return true
}

func (s *Source) randomKey() apiKey {
	if len(s.apiKeys) == 0 {
		return apiKey{}
	}
	return s.apiKeys[rand.Intn(len(s.apiKeys))]
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for fofa")
	}

	randomApiKey := s.randomKey()
	if randomApiKey.username == "" || randomApiKey.secret == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	// fofa api doc https://fofa.info/static_pages/api_help
	qbase64 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("domain=\"%s\"", query)))
	url := fmt.Sprintf("https://fofa.info/api/v1/search/all?full=true&fields=host&page=1&size=10000&email=%s&key=%s&qbase64=%s", randomApiKey.username, randomApiKey.secret, qbase64)
	
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	var response fofaResponse
	err = jsoniter.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Error {
		return nil, fmt.Errorf("%s", response.ErrMsg)
	}

	if response.Size > 0 {
		for _, subdomain := range response.Results {
			if strings.HasPrefix(strings.ToLower(subdomain), "http://") || strings.HasPrefix(strings.ToLower(subdomain), "https://") {
				subdomain = subdomain[strings.Index(subdomain, "//")+2:]
			}
			re := regexp.MustCompile(`:\d+$`)
			if re.MatchString(subdomain) {
				subdomain = re.ReplaceAllString(subdomain, "")
			}
			results = append(results, subdomain)
		}
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

