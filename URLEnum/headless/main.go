package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/publicsuffix"
)

func main() {
	inFile := flag.String("l", "", "Input file with hosts or URLs (one per line)")
	outFile := flag.String("o", "", "Output file (one URL per line)")
	concurrency := flag.Int("c", 3, "Concurrent targets")
	timeout := flag.Duration("timeout", 30*time.Second, "Timeout per target")
	wait := flag.Duration("wait", 8*time.Second, "Wait after load")
	chromePath := flag.String("chrome", "/usr/bin/google-chrome", "Chrome path")
	flag.Parse()

	if strings.TrimSpace(*inFile) == "" || strings.TrimSpace(*outFile) == "" {
		log.Fatal("use -l input and -o output")
	}

	targets, err := readLines(*inFile)
	if err != nil {
		log.Fatalf("read input: %v", err)
	}

	// Output file
	f, err := os.Create(*outFile)
	if err != nil {
		log.Fatalf("create output: %v", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	// One writer goroutine: dedupe + write (no mutexes)
	results := make(chan string, 2048)
	done := make(chan struct{})
	go func() {
		defer close(done)
		seen := map[string]struct{}{}
		for u := range results {
			if _, ok := seen[u]; ok {
				continue
			}
			seen[u] = struct{}{}
			_, _ = w.WriteString(u + "\n")
		}
		fmt.Printf("Done. Wrote %d unique URL(s) to %s\n", len(seen), *outFile)
	}()

	// One Chrome instance for everything
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(*chromePath),
			chromedp.Flag("headless", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("disable-dev-shm-usage", true),
		)...,
	)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	sem := make(chan struct{}, *concurrency)
	errs := make(chan error, len(targets))

	// Run targets concurrently, limited by sem
	for _, line := range targets {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		targetURL := normalizeTarget(line)
		targetCore := coreDomainFromURL(targetURL)
		if targetURL == "" || targetCore == "" {
			continue
		}

		sem <- struct{}{}
		go func(targetURL, targetCore string) {
			defer func() { <-sem }()

			if err := scanTarget(browserCtx, targetURL, targetCore, *timeout, *wait, func(u string) {
				u = normalizeCapturedURL(u)
				if u == "" {
					return
				}
				if !isInformational(u) {
					return
				}
				if coreDomainFromURL(u) != targetCore {
					return
				}
				results <- u
			}); err != nil {
				errs <- fmt.Errorf("%s -> %w", targetURL, err)
			}
		}(targetURL, targetCore)
	}

	// Wait for all goroutines to release semaphore
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
	// cleanup
	close(results)
	<-done
	close(errs)

	for e := range errs {
		log.Println(e)
	}
}

func scanTarget(browserCtx context.Context, targetURL, targetCore string, perTargetTimeout, wait time.Duration, onURL func(string)) error {
	ctx, cancel := chromedp.NewContext(browserCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, perTargetTimeout)
	defer cancel()

	chromedp.ListenTarget(ctx, func(ev any) {
		if r, ok := ev.(*network.EventResponseReceived); ok && r.Response.URL != "" {
			onURL(r.Response.URL)
		}
	})

	return chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(wait),
	)
}

func normalizeTarget(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return s
	}
	return "https://" + s
}

func normalizeCapturedURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	u.Fragment = ""
	if u.Path != "/" && strings.HasSuffix(u.Path, "/") {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}
	return u.String()
}

func isInformational(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	ext := strings.ToLower(path.Ext(u.Path))
	switch ext {
	case "", ".js", ".html",
		".php", ".phtml", ".php3", ".php4", ".php5", ".phps",
		".asp", ".aspx", ".ashx", ".asmx", ".axd",
		".jsp", ".jspx", ".do", ".action",
		".py", ".rb", ".pl", ".cgi",
		".cfm", ".cfc", ".mjs",
		".lua",
		".go",
		".fcgi", ".ejs", ".erb", ".twig", ".jinja", ".j2",
		".hbs", ".handlebars", ".mustache",
		".liquid", ".ftl", ".vm", ".htm",
		".wasm", ".json", ".xml", ".graphql", ".wsdl", ".yaml", ".yml", ".txt":
		return true
	default:
		return false
	}
}

// coreDomainFromURL returns the registrable-domain label without the public suffix.
// e.g. api.example.co.uk -> "example", example.net -> "example"
func coreDomainFromURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return ""
	}
	host := strings.ToLower(u.Host)
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	etld1, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return ""
	}
	suf, ok := publicsuffix.PublicSuffix(etld1)
    if !ok {
      return ""
     }

	core := strings.TrimSuffix(etld1, "."+suf) // "example" or "a.example"
	core = strings.TrimSuffix(core, ".")
	parts := strings.Split(core, ".")
	return parts[len(parts)-1]
}

func readLines(p string) ([]string, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}
