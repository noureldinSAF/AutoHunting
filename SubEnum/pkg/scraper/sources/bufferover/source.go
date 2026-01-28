package bufferover

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Source struct {
	apiKeys []string
}

type response struct {
	Meta struct {
		Errors []string `json:"Errors"`
	} `json:"Meta"`
	FDNSA   []string `json:"FDNS_A"`
	RDNS    []string `json:"RDNS"`
	Results []string `json:"Results"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "bufferover"
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
		return nil, fmt.Errorf("no API keys provided for bufferover")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	sourceURL := fmt.Sprintf("https://tls.bufferover.run/dns?q=.%s", query)
	req, err := http.NewRequest(http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", randomApiKey)

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

	var bufforesponse response
	err = jsoniter.Unmarshal(body, &bufforesponse)
	if err != nil {
		return nil, err
	}

	metaErrors := bufforesponse.Meta.Errors
	if len(metaErrors) > 0 {
		return nil, fmt.Errorf("%s", strings.Join(metaErrors, ", "))
	}

	var subdomains []string

	if len(bufforesponse.FDNSA) > 0 {
		subdomains = bufforesponse.FDNSA
		subdomains = append(subdomains, bufforesponse.RDNS...)
	} else if len(bufforesponse.Results) > 0 {
		subdomains = bufforesponse.Results
	}

	results = append(results, subdomains...)
	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

