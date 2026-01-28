package scraper

import (
	"crypto/tls"
	"net/http"
	"subenum/internal/config"
	"time"
)

type Source interface {
	Name() string
	Search(query string, client *http.Client) ([]string, error)
	RequiresAPIKey() bool
}

func NewSession(timeout int) *http.Client {
	return &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

// ExtractALLAPIKeys reads config and returns the map of service -> api-keys.
// This returns the top-level `config` mapping from your YAML, which is stored
// in cfg.Config (type map[string][]string).
func ExtractALLAPIKeys() (map[string][]string, error) {
	cfg, err := config.ReadConfig()
	if err != nil {
		return nil, err
	}

	return cfg.Config, nil
}
