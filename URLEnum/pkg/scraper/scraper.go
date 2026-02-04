package scraper

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"
	"github.com/noureldinSAF/AutoHunting/URLEnum/internal/config"
    "net"
)

type Source interface {
	Name() string
	Search(ctx context.Context, query string, client *http.Client) ([]string, error)
	RequireAPIKey() bool
}

func NewSession(timeout int) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},

		// Connection pool tuning
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     90 * time.Second,

		// Dialer & timeouts
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second, // TCP dial timeout
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second,
	}

	return &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: tr,
	}
}


func ExtractALLAPIKeys() (map[string][]string, error) {
	cfg, err := config.ReadConfig()
	if err != nil {
		return nil, err
	}

	return cfg.Config, nil
}