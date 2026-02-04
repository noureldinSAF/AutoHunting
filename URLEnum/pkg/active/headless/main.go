package headless

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/publicsuffix"
)

type Options struct {
	Concurrency     int
	Timeout         time.Duration // per target
	Wait            time.Duration // after load
	ChromePath      string
	Headless        bool
	NoSandbox       bool
	DisableGPU      bool
	DisableDevShm   bool
	ExtraAllocator  []chromedp.ExecAllocatorOption
}

func (o Options) withDefaults() Options {
	if o.Concurrency <= 0 {
		o.Concurrency = 3
	}
	if o.Timeout <= 0 {
		o.Timeout = 30 * time.Second
	}
	if o.Wait <= 0 {
		o.Wait = 8 * time.Second
	}
	if strings.TrimSpace(o.ChromePath) == "" {
		o.ChromePath = "/usr/bin/google-chrome"
	}
	// Default behavior matches original main(): headless true + hardening flags
	if !o.Headless {
		o.Headless = true
	}
	// Keep these enabled by default (same spirit as original)
	if !o.NoSandbox {
		o.NoSandbox = true
	}
	if !o.DisableGPU {
		o.DisableGPU = true
	}
	if !o.DisableDevShm {
		o.DisableDevShm = true
	}
	return o
}

// Enumerate scans a single start target (host or URL) and returns unique informational URLs.
// includeSubdomains controls whether we allow URLs from the same registrable domain (eTLD+1).
// - false: only same hostname as the start URL.
// - true : any subdomain under the same eTLD+1 as the start URL.
func Enumerate(ctx context.Context, start string, includeSubdomains bool, opts Options) ([]string, error) {
	opts = opts.withDefaults()

	startURL := normalizeTarget(start)
	if startURL == "" {
		return nil, errors.New("empty start")
	}

	u0, err := url.Parse(startURL)
	if err != nil || u0.Host == "" {
		return nil, fmt.Errorf("invalid start url: %q", start)
	}

	var allowFn func(candidate string) bool
	if includeSubdomains {
		rootETLD1 := etldPlusOne(u0.Hostname())
		if rootETLD1 == "" {
			return nil, fmt.Errorf("cannot compute eTLD+1 for start host: %s", u0.Hostname())
		}
		allowFn = func(candidate string) bool {
			uu, err := url.Parse(candidate)
			if err != nil || uu.Hostname() == "" {
				return false
			}
			return etldPlusOne(uu.Hostname()) == rootETLD1
		}
	} else {
		rootHost := strings.ToLower(u0.Hostname())
		allowFn = func(candidate string) bool {
			uu, err := url.Parse(candidate)
			if err != nil || uu.Hostname() == "" {
				return false
			}
			return strings.ToLower(uu.Hostname()) == rootHost
		}
	}

	// One Chrome instance (allocator) per Enumerate call
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(opts.ChromePath),
		chromedp.Flag("headless", opts.Headless),
		chromedp.Flag("no-sandbox", opts.NoSandbox),
		chromedp.Flag("disable-gpu", opts.DisableGPU),
		chromedp.Flag("disable-dev-shm-usage", opts.DisableDevShm),
	)
	if len(opts.ExtraAllocator) > 0 {
		allocOpts = append(allocOpts, opts.ExtraAllocator...)
	}

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, allocOpts...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	// Dedupe + collect (like your writer goroutine, but in-memory)
	results := make(chan string, 2048)
	var (
		mu   sync.Mutex
		seen = map[string]struct{}{}
		out  = make([]string, 0, 512)
	)

	var collectorWG sync.WaitGroup
	collectorWG.Add(1)
	go func() {
		defer collectorWG.Done()
		for s := range results {
			mu.Lock()
			if _, ok := seen[s]; ok {
				mu.Unlock()
				continue
			}
			seen[s] = struct{}{}
			out = append(out, s)
			mu.Unlock()
		}
	}()

	sem := make(chan struct{}, opts.Concurrency)
	errCh := make(chan error, 1)

	// This tool version enumerates a single target using concurrency=1 logic,
	// but keeps the same semaphore pattern so you can extend to multi-target easily.
	sem <- struct{}{}
	go func(targetURL string) {
		defer func() { <-sem }()
		err := scanTarget(browserCtx, targetURL, opts.Timeout, opts.Wait, func(raw string) {
			u := normalizeCapturedURL(raw)
			if u == "" {
				return
			}
			if !isInformational(u) {
				return
			}
			if !allowFn(u) {
				return
			}
			select {
			case results <- u:
			case <-ctx.Done():
				return
			}
		})
		if err != nil {
			select {
			case errCh <- fmt.Errorf("%s -> %w", targetURL, err):
			default:
			}
		}
	}(startURL)

	// Wait for goroutine by filling semaphore completely
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}

	close(results)
	collectorWG.Wait()

	select {
	case e := <-errCh:
		// Return partial results + error
		return out, e
	default:
	}

	if ctx.Err() != nil {
		return out, ctx.Err()
	}
	return out, nil
}

func scanTarget(browserCtx context.Context, targetURL string, perTargetTimeout, wait time.Duration, onURL func(string)) error {
	ctx, cancel := chromedp.NewContext(browserCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, perTargetTimeout)
	defer cancel()

	chromedp.ListenTarget(ctx, func(ev any) {
		if r, ok := ev.(*network.EventResponseReceived); ok && r.Response.URL != "" {
			onURL(r.Response.URL)
		}
	})

	var hrefsJSON string

	return chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(wait),
		chromedp.Evaluate("`JSON.stringify(Array.from(document.querySelectorAll(\"a[href]\").map(a => a.href)))`", &hrefsJSON),
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

func etldPlusOne(host string) string {
	host = strings.ToLower(strings.TrimSpace(host))
	if host == "" {
		return ""
	}
	// Strip port if any
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	etld1, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return ""
	}
	return etld1
}
