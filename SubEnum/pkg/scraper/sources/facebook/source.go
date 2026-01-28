package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cyinnove/logify"
)

type Source struct {
	apiKeys []apiKey
}

type apiKey struct {
	AppID       string
	Secret      string
	AccessToken string
	Error       error
}

type authResponse struct {
	AccessToken string `json:"access_token"`
}

type response struct {
	Data []struct {
		Domains []string `json:"domains"`
	} `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

const (
	domainsPerPage = "1000"
	authUrl        = "https://graph.facebook.com/oauth/access_token?client_id=%s&client_secret=%s&grant_type=client_credentials"
	domainsUrl     = "https://graph.facebook.com/certificates?fields=domains&access_token=%s&query=%s&limit=" + domainsPerPage
)

func New(apiKeys []string) *Source {
	keys := make([]apiKey, 0)
	for _, key := range apiKeys {
		parts := strings.Split(key, ":")
		if len(parts) == 2 {
			ak := apiKey{AppID: parts[0], Secret: parts[1]}
			ak.FetchAccessToken()
			if ak.Error == nil && ak.AccessToken != "" {
				keys = append(keys, ak)
			}
		}
	}
	return &Source{apiKeys: keys}
}

func (k *apiKey) FetchAccessToken() {
	if k.AppID == "" || k.Secret == "" {
		k.Error = fmt.Errorf("invalid app id or secret")
		return
	}

	resp, err := http.Get(fmt.Sprintf(authUrl, k.AppID, k.Secret))
	if err != nil {
		k.Error = err
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		k.Error = err
		return
	}

	var auth authResponse
	if err := json.Unmarshal(body, &auth); err != nil {
		k.Error = err
		return
	}

	if auth.AccessToken == "" {
		k.Error = fmt.Errorf("invalid response from facebook")
		return
	}

	k.AccessToken = auth.AccessToken
}

func (s *Source) Name() string {
	return "facebook"
}

func (s *Source) RequiresAPIKey() bool {
	return true
}

func (s *Source) randomKey() apiKey {
	if len(s.apiKeys) == 0 {
		return apiKey{}
	}
	return s.apiKeys[rand.Intn(len(s.apiKeys))]
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var allResults []string

	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided for facebook")
	}

	key := s.randomKey()
	if key.AccessToken == "" {
		return nil, fmt.Errorf("no valid API key available")
	}

	// URL encode the query parameter
	domainsURL := fmt.Sprintf(domainsUrl, key.AccessToken, url.QueryEscape(query))

	for {
		req, err := http.NewRequest(http.MethodGet, domainsURL, nil)
		if err != nil {
			return allResults, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return allResults, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return allResults, err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			// Check for specific Facebook API errors
			var fbError struct {
				Error struct {
					Message      string `json:"message"`
					Type         string `json:"type"`
					Code         int    `json:"code"`
					ErrorSubcode int    `json:"error_subcode"`
				} `json:"error"`
			}
			if err := json.Unmarshal(body, &fbError); err == nil && fbError.Error.Code != 0 {
				if fbError.Error.Code == 100 && fbError.Error.ErrorSubcode == 33 {
					return allResults, fmt.Errorf("facebook API error: %s (Code: %d, Subcode: %d). The app may need Certificate Transparency API permissions. Check https://developers.facebook.com/tools/ct", fbError.Error.Message, fbError.Error.Code, fbError.Error.ErrorSubcode)
				}
				return allResults, fmt.Errorf("facebook API error: %s (Code: %d)", fbError.Error.Message, fbError.Error.Code)
			}
			return allResults, fmt.Errorf("api error: status=%d body=%s", resp.StatusCode, string(body))
		}

		var fbResponse response
		if err := json.Unmarshal(body, &fbResponse); err != nil {
			return allResults, err
		}

		for _, v := range fbResponse.Data {
			allResults = append(allResults, v.Domains...)
		}

		if fbResponse.Paging.Next == "" {
			break
		}

		domainsURL = updateParamInURL(fbResponse.Paging.Next, "limit", domainsPerPage)
	}

	return allResults, nil
}

func updateParamInURL(urlStr, param, value string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}
	q := u.Query()
	q.Set(param, value)
	u.RawQuery = q.Encode()
	return u.String()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	source := New([]string{"938193287645247:054ebfbfcc8d60a867fd1edacaf4817b"})
	results, err := source.Search("example.com", &http.Client{})
	if err != nil {
		logify.Errorf("Error searching for domains: %v", err)
		return
	}
	logify.Infof("Found %d domains: %v", len(results), results)
}
