package scraper

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"
	"github.com/noureldinSAF/AutoHunting/URLEnum/internal/config"
    //"net"
)

type Source interface {
	Name() string
	Search(ctx context.Context, query string, client *http.Client) ([]string, error)
	RequireAPIKey() bool
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


func ExtractALLAPIKeys() (map[string][]string, error) {
	cfg, err := config.ReadConfig()
	if err != nil {
		return nil, err
	}

	return cfg.Config, nil
}