package robtex

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Source struct {
	apiKeys []string
}

const (
	addrRecord     = "A"
	iPv6AddrRecord = "AAAA"
	baseURL        = "https://proapi.robtex.com/pdns"
)

type result struct {
	Rrname string `json:"rrname"`
	Rrdata string `json:"rrdata"`
	Rrtype string `json:"rrtype"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "robtex"
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
	var allResults []string

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for robtex")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	headers := map[string]string{"Content-Type": "application/x-ndjson"}

	ips, err := s.enumerate(client, fmt.Sprintf("%s/forward/%s?key=%s", baseURL, query, randomApiKey), headers)
	if err != nil {
		return nil, err
	}

	for _, result := range ips {
		if result.Rrtype == addrRecord || result.Rrtype == iPv6AddrRecord {
			domains, err := s.enumerate(client, fmt.Sprintf("%s/reverse/%s?key=%s", baseURL, result.Rrdata, randomApiKey), headers)
			if err != nil {
				continue
			}
			for _, result := range domains {
				allResults = append(allResults, result.Rrdata)
			}
		}
	}

	return allResults, nil
}

func (s *Source) enumerate(client *http.Client, targetURL string, headers map[string]string) ([]result, error) {
	var results []result

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return results, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var response result
		err = json.NewDecoder(bytes.NewBufferString(line)).Decode(&response)
		if err != nil {
			continue
		}

		results = append(results, response)
	}

	return results, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

