package dnsdb

import (
	"bufio"
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
	apiKeys []string
}

const urlBase string = "https://api.dnsdb.info/dnsdb/v2"

type rateResponse struct {
	Rate rate
}

type rate struct {
	OffsetMax json.Number `json:"offset_max"`
}

type safResponse struct {
	Condition string   `json:"cond"`
	Obj       dnsdbObj `json:"obj"`
	Msg       string   `json:"msg"`
}

type dnsdbObj struct {
	Name string `json:"rrname"`
}

func New(apiKeys []string) *Source {
	return &Source{
		apiKeys: apiKeys,
	}
}

func (s *Source) Name() string {
	return "dnsdb"
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
	seen := make(map[string]bool)

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for dnsdb")
	}

	randomApiKey := s.randomKey()
	if randomApiKey == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	headers := map[string]string{
		"X-API-KEY": randomApiKey,
		"Accept":    "application/x-ndjson",
	}

	offsetMax, err := s.getMaxOffset(client, headers)
	if err != nil {
		return allResults, err
	}

	path := fmt.Sprintf("lookup/rrset/name/*.%s", query)
	urlTemplate := fmt.Sprintf("%s/%s?", urlBase, path)
	queryParams := url.Values{}
	queryParams.Add("limit", "0")
	queryParams.Add("swclient", "subfinder")

	for {
		url := urlTemplate + queryParams.Encode()

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return allResults, err
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			return allResults, err
		}

		var respCond string
		reader := bufio.NewReader(resp.Body)
		for {
			n, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				resp.Body.Close()
				return allResults, err
			}

			var response safResponse
			err = json.Unmarshal(n, &response)
			if err != nil {
				resp.Body.Close()
				return allResults, err
			}

			respCond = response.Condition
			if respCond == "" || respCond == "ongoing" {
				if response.Obj.Name != "" {
					name := strings.TrimSuffix(response.Obj.Name, ".")
					if !seen[name] {
						seen[name] = true
						allResults = append(allResults, name)
					}
				}
			} else if respCond != "begin" {
				break
			}
		}

		resp.Body.Close()

		if respCond == "limited" {
			if offsetMax != 0 && len(allResults) <= int(offsetMax) {
				queryParams.Set("offset", strconv.Itoa(len(allResults)))
				continue
			}
		} else if respCond != "succeeded" {
			return allResults, fmt.Errorf("dnsdb terminated with condition: %s", respCond)
		}

		break
	}

	return allResults, nil
}

func (s *Source) getMaxOffset(client *http.Client, headers map[string]string) (uint64, error) {
	var offsetMax uint64
	url := fmt.Sprintf("%s/rate_limit", urlBase)
	
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return offsetMax, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return offsetMax, err
	}
	defer resp.Body.Close()

	var rateResp rateResponse
	err = json.NewDecoder(resp.Body).Decode(&rateResp)
	if err != nil {
		return offsetMax, err
	}

	if offsetMaxNum, err := rateResp.Rate.OffsetMax.Int64(); err == nil {
		offsetMax = uint64(offsetMaxNum)
	}

	return offsetMax, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

