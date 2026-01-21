package whoisfreaks

import (
	"encoding/json"
	"fmt"
	"io"

	//"log"
	"math/rand"
	"net/http"
	"net/url"

	"time"

	"github.com/cyinnove/logify"
)

type Source struct {
	apiKeys []string
}

type responseObj struct {
	TotalResult int `json:"total_Result"`
	TotalPages  int `json:"total_Pages"`
	CurrentPage int `json:"current_Page"`
	Domains     []struct {
		DomainName string `json:"domain_name"`
	} `json:"whois_domains_historical"`
}

func New(apiKeys []string) *Source {
	return &Source{apiKeys: apiKeys}
}

func (s *Source) Name() string {
	return "whoisfreaks"
}

func (s *Source) RequireAPIKey() bool {
	return true
}

func (s *Source) randomKey() string {
	return s.apiKeys[rand.Intn(len(s.apiKeys))]
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var domains []string
	if s.apiKeys == nil || len(s.apiKeys) == 0 {
		return []string{}, fmt.Errorf("No API keys provided for whoisfreaks")
	}

	page := 1

	for {

		fmtURL := fmt.Sprintf("https://api.whoisfreaks.com/v1.0/whois?apiKey=%s&whois=reverse&keyword=%s&page=%d", s.randomKey(), url.QueryEscape(query), page)

		req, err := http.NewRequest("GET", fmtURL, nil)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode < 200 || resp.StatusCode > 300 {
			return nil, fmt.Errorf("api error: status code %d, body:%s", resp.StatusCode, string(body))
		}

		var ro responseObj
		if err := json.Unmarshal(body, &ro); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}

		logify.Infof("whoisfreaks: page%d/%d got=%d domains  totoalResults=%d\n", ro.CurrentPage, ro.TotalPages, len(ro.Domains), ro.TotalResult)

		for _, domain := range ro.Domains {
			if domain.DomainName == "" {
				continue
			}
			domains = append(domains, domain.DomainName)
		}

		if ro.CurrentPage >= 5 {
			break
		}
		page++
		time.Sleep(1 * time.Minute)
	}
	return domains, nil
}

// func main() {

// 	logify.MaxLevel = logify.Info

// 	s := New([]string{"446dab4ffd6346168998058953c55852"})

// 	domains, err := s.Search("google", http.DefaultClient)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(domains)
// }
