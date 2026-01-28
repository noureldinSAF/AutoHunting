package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/corpix/uarand"
	"github.com/gocolly/colly/v2"
)


// normalizeURL normalizes a URL by removing fragments and trailing slashes, but keeps query params
func normalizeURL(u *url.URL) string {
	path := u.Path
	// Remove trailing slash except for root path
	if path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	normalized := u.Scheme + "://" + u.Host + path
	if u.RawQuery != "" {
		normalized += "?" + u.RawQuery
	}
	return normalized
}

func main() {
	// Define command-line flags
	depth := flag.Int("depth", 2, "Maximum depth to crawl (default: 2)")
	timeout := flag.Duration("timeout", 30*time.Second, "Request timeout (default: 30s)")
	threads := flag.Int("threads", 5, "Number of concurrent threads (default: 5)")
	verbose := flag.Bool("verbose", false, "Show verbose output including errors")
	outputFile := flag.String("output", "", "Output file path (default: stdout)")
	flag.Parse()

	// Open output file if specified
	var output *os.File = os.Stdout
	var err error
	if *outputFile != "" {
		output, err = os.Create(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		targetURL := strings.TrimSpace(scanner.Text())
		if targetURL == "" {
			continue
		}

		parsedURL, err := url.Parse(targetURL)
		if err != nil {
			if *verbose {
				fmt.Fprintf(os.Stderr, "Error parsing URL %s: %v\n", targetURL, err)
			}
			continue
		}

		// Create a new collector with configurable depth and async crawling
		c := colly.NewCollector(
			colly.MaxDepth(*depth),
			colly.Async(true),
			colly.AllowedDomains(parsedURL.Hostname()),
		)

		// Set user agent
		c.UserAgent = uarand.GetRandom()

		// Set timeout
		c.SetRequestTimeout(*timeout)

		// Set thread limit
		c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: *threads,
		})

		// Track seen URLs to avoid duplicates
		seen := make(map[string]bool)

		allowedHostname := parsedURL.Hostname()

		// Extract and print all links from HTML elements
		c.OnHTML("a[href], link[href], script[src], iframe[src], form[action]", func(e *colly.HTMLElement) {
			var linkURL string
			switch {
			case e.Attr("href") != "":
				linkURL = e.Attr("href")
			case e.Attr("src") != "":
				linkURL = e.Attr("src")
			case e.Attr("action") != "":
				linkURL = e.Attr("action")
			}

			if linkURL != "" {
			
				if absoluteURL := e.Request.AbsoluteURL(linkURL); absoluteURL != "" {
					printUnique(seen, absoluteURL, output, allowedHostname)
				}
			}
		})

		// Print URL when visiting and skip static assets
		c.OnRequest(func(r *colly.Request) {
			printUnique(seen, r.URL.String(), output, allowedHostname)			

		})

		// Handle errors
		c.OnError(func(r *colly.Response, err error) {
			if *verbose {
				fmt.Fprintf(os.Stderr, "Error visiting %s: %v\n", r.Request.URL, err)
			}
		})

		// Start crawling
		if err := c.Visit(targetURL); err != nil {
			if *verbose {
				fmt.Fprintf(os.Stderr, "Error visiting URL %s: %v\n", targetURL, err)
			}
			continue
		}

		c.Wait()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}

// printUnique prints a URL if it hasn't been seen before and is in the allowed domain
func printUnique(seen map[string]bool, urlStr string, output *os.File, allowedHostname string) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	// Check if URL is in allowed domain BEFORE processing
	if parsed.Hostname() != allowedHostname {
		return
	}

	normalized := normalizeURL(parsed)

	// Print if we haven't seen this URL before
	if !seen[normalized] {
		seen[normalized] = true
		fmt.Fprintln(output, normalized)
	}
}
