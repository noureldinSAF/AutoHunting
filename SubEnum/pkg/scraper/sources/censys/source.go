package censys

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Source struct {
	apiKeys []apiKey
}

type apiKey struct {
	token  string
	secret string
}

type response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Result result `json:"result"`
}

type result struct {
	Query      string  `json:"query"`
	Total      float64 `json:"total"`
	DurationMS int     `json:"duration_ms"`
	Hits       []hit   `json:"hits"`
	Links      links   `json:"links"`
}

type hit struct {
	Parsed            parsed   `json:"parsed"`
	Names             []string `json:"names"`
	FingerprintSha256 string   `json:"fingerprint_sha256"`
}

type parsed struct {
	ValidityPeriod validityPeriod `json:"validity_period"`
	SubjectDN      string         `json:"subject_dn"`
	IssuerDN       string         `json:"issuer_dn"`
}

type validityPeriod struct {
	NotAfter  string `json:"not_after"`
	NotBefore string `json:"not_before"`
}

type links struct {
	Next string `json:"next"`
	Prev string `json:"prev"`
}

const (
	maxCensysPages = 10
	maxPerPage     = 100
)

func New(apiKeys []string) *Source {
	keys := make([]apiKey, 0)
	for _, key := range apiKeys {
		parts := strings.Split(key, ":")
		if len(parts) == 2 {
			keys = append(keys, apiKey{token: parts[0], secret: parts[1]})
		}
	}
	return &Source{apiKeys: keys}
}

func (s *Source) Name() string {
	return "censys"
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
	var allResults []string

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for censys")
	}

	randomApiKey := s.randomKey()
	if randomApiKey.token == "" || randomApiKey.secret == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	certSearchEndpoint := "https://search.censys.io/api/v2/certificates/search"
	cursor := ""
	currentPage := 1

	for {
		u, err := url.Parse(certSearchEndpoint)
		if err != nil {
			return allResults, err
		}

		q := u.Query()
		q.Set("q", query)
		q.Set("per_page", strconv.Itoa(maxPerPage))
		if cursor != "" {
			q.Set("cursor", cursor)
		}
		u.RawQuery = q.Encode()

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return allResults, err
		}

		// Basic auth
		auth := base64.StdEncoding.EncodeToString([]byte(randomApiKey.token + ":" + randomApiKey.secret))
		req.Header.Set("Authorization", "Basic "+auth)

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

		var censysResponse response
		err = json.Unmarshal(body, &censysResponse)
		if err != nil {
			return allResults, err
		}

		for _, hit := range censysResponse.Result.Hits {
			allResults = append(allResults, hit.Names...)
		}

		// Exit the censys enumeration if last page is reached
		cursor = censysResponse.Result.Links.Next
		if cursor == "" || currentPage >= maxCensysPages {
			break
		}
		currentPage++
	}

	return allResults, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

