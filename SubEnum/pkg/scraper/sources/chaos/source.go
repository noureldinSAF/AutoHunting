package chaos

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/projectdiscovery/chaos-client/pkg/chaos"
)

type Source struct {
	apiKeys []string
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "chaos"
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
		return nil, fmt.Errorf("no API keys provided for chaos")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	chaosClient := chaos.New(randomApiKey)
	for result := range chaosClient.GetSubdomains(&chaos.SubdomainsRequest{
		Domain: query,
	}) {
		if result.Error != nil {
			return results, result.Error
		}
		results = append(results, fmt.Sprintf("%s.%s", result.Subdomain, query))
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

