package urlscan

import (
	"encoding/json"
	"fmt"

	//"log"
	"math/rand"
	"net/http"
	"net/url"

	"time"
	"context"
	"strings"
	

	//"github.com/cyinnove/logify"
)

type Source struct { apiKeys []string }

func (s *Source) Name() string {
	return "urlscan"
}

func (s *Source) RequireAPIKey() bool {
	return true
}


// type responseObj struct {
// 	TotalResult int `json:"total_Result"`
// 	TotalPages  int `json:"total_Pages"`
// 	CurrentPage int `json:"current_Page"`
// 	Domains     []struct {
// 		DomainName string `json:"domain_name"`
// 	} `json:"whois_domains_historical"`
// }

func New(apiKeys []string) *Source {
	return &Source{apiKeys: apiKeys}
}



func (s *Source) randomKey() string {
	return s.apiKeys[rand.Intn(len(s.apiKeys))]
}

func (s *Source) Search(ctx context.Context, query string, client *http.Client) ([]string, error) {

	var urls []string

	if query == "" {
		return urls, nil
	}

	qs := url.Values{}
	qs.Set("q", fmt.Sprintf("domain:%s", query))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://urlscan.io/api/v1/search/?q="+qs.Encode(), nil)
	if err != nil {
		return urls, err
	}
	req.Header.Set("User-Agent", "uarand.GetRandom()") 

	resp, err := client.Do(req)
	if err != nil {
		return urls, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("webarchive returned non-200 status code: %d", resp.StatusCode)
	}

	var parsed struct {
		Results []struct {
			Page struct {
				URL string `json:"url"`
			} `json:"page"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	out := make([]string, 0, len(parsed.Results))
	for _, r := range parsed.Results {
		out = append(out, r.Page.URL)

		u := strings.TrimSpace(r.Page.URL)
		if u == "" {
			continue
		}
		
		out = append(out, u)
	}

	return out, nil
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
