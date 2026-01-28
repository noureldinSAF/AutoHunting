package http

import (
	"crypto/tls"
	"net/http"
	"time"
)

func getClient(timeout int) *http.Client {
	return &http.Client{
		Timeout: time.Second * time.Duration(timeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
