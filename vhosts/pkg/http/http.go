package http

import (
	"io"
	"net/http"
)

type Response struct {
	StatusCode    int
	ContentLength int64
}

func GetResponse(timeout int, host, ip string) (*Response, error) {

	req, err := http.NewRequest(http.MethodGet, ip, nil)
	if err != nil {
		return nil, err
	}

	req.Host = host

	resp, err := getClient(timeout).Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.ContentLength == -1 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return &Response{ContentLength: int64(len(body)), StatusCode: resp.StatusCode}, nil
	}

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		ContentLength: resp.ContentLength,
		StatusCode:    resp.StatusCode,
	}, nil
}
