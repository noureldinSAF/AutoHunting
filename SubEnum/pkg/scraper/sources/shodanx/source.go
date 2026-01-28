package shodanx

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Source struct{}

func (s *Source) Name() string {
	return "shodanx"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	url := fmt.Sprintf("https://www.shodan.io/domain/%s", query)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return results, nil
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	// Find the <ul id="subdomains"> element
	var subdomainsList *html.Node
	var findSubdomainsList func(*html.Node)
	findSubdomainsList = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "ul" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == "subdomains" {
					subdomainsList = n
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findSubdomainsList(c)
		}
	}
	findSubdomainsList(doc)

	if subdomainsList == nil {
		return results, nil
	}

	// Extract all <li> elements from the subdomains list
	var extractLiText func(*html.Node)
	extractLiText = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			var getText func(*html.Node) string
			getText = func(node *html.Node) string {
				if node.Type == html.TextNode {
					return node.Data
				}
				var text string
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					text += getText(c)
				}
				return text
			}

			text := strings.TrimSpace(getText(n))
			if text != "" {
				subdomain := text + "." + query
				results = append(results, subdomain)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractLiText(c)
		}
	}
	extractLiText(subdomainsList)

	return results, nil
}

