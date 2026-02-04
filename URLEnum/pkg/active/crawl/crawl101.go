package crawl

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/corpix/uarand"
	"github.com/gocolly/colly/v2"
)

type Options struct {
	MaxDepth     int
	Parallelism  int
	Timeout      time.Duration
	AllowQuery   bool // keep query params (default true if you want)
}

// normalizeURL normalizes a URL by removing fragments and trailing slashes.
// It keeps query params (if opts.AllowQuery == true).
func normalizeURL(u *url.URL, allowQuery bool) string {
	if u == nil {
		return ""
	}

	// Copy to avoid mutating caller
	uu := *u
	uu.Fragment = ""

	path := uu.Path
	// Remove trailing slash except for root path
	if path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	uu.Path = path

	// Build
	normalized := uu.Scheme + "://" + uu.Host + uu.Path
	if allowQuery && uu.RawQuery != "" {
		normalized += "?" + uu.RawQuery
	}
	return normalized
}

func defaultOptions(opts *Options) Options {
	if opts == nil {
		return Options{
			MaxDepth:    2,
			Parallelism: 4,
			Timeout:     15 * time.Second,
			AllowQuery:  true,
		}
	}
	o := *opts
	if o.MaxDepth <= 0 {
		o.MaxDepth = 2
	}
	if o.Parallelism <= 0 {
		o.Parallelism = 4
	}
	if o.Timeout <= 0 {
		o.Timeout = 15 * time.Second
	}
	// AllowQuery: keep whatever user set; if zero-value false is desired, caller can set it.
	// If you want default true always, uncomment:
	// if !opts.AllowQuery { o.AllowQuery = true }
	return o
}

// hostAllowed checks whether hostname is allowed relative to the rootHost.
// If includeSubDomains is true, subdomains of rootHost are allowed.
func hostAllowed(host, rootHost string, includeSubDomains bool) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	rootHost = strings.ToLower(strings.TrimSpace(rootHost))

	if host == "" || rootHost == "" {
		return false
	}
	if host == rootHost {
		return true
	}
	if !includeSubDomains {
		return false
	}
	// Allow *.rootHost
	return strings.HasSuffix(host, "."+rootHost)
}

// Enumerate crawls starting from `start` and returns a unique list of normalized URLs visited/discovered
// (depending on what you emit). This keeps your “emit on request” logic and also visits discovered links.
func Enumerate(ctx context.Context, start string, includeSubDomains bool, opts *Options) ([]string, error) {
	o := defaultOptions(opts)

	if start == "" {
		return nil, errors.New("start URL is empty")
	}
	if !strings.HasPrefix(start, "http://") && !strings.HasPrefix(start, "https://") {
		start = "https://" + start
	}

	startURL, err := url.Parse(start)
	if err != nil {
		return nil, fmt.Errorf("parse start url: %w", err)
	}
	if startURL.Scheme == "" || startURL.Hostname() == "" {
		return nil, fmt.Errorf("invalid start url: %q", start)
	}

	rootHost := strings.ToLower(startURL.Hostname())

	c := colly.NewCollector(
		colly.MaxDepth(o.MaxDepth),
		colly.Async(true),
	)

	c.UserAgent = uarand.GetRandom()
	c.SetRequestTimeout(o.Timeout)

	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: o.Parallelism,
	})

	var (
		mu   sync.Mutex
		seen = make(map[string]struct{})
		out  = make([]string, 0, 256)
	)

	emit := func(u *url.URL) {
		if u == nil {
			return
		}
		// Ensure scheme/host exist for normalize
		if u.Scheme == "" || u.Host == "" {
			return
		}
		// Domain rule
		if !hostAllowed(u.Hostname(), rootHost, includeSubDomains) {
			return
		}

		normalized := normalizeURL(u, o.AllowQuery)
		if normalized == "" {
			return
		}

		mu.Lock()
		if _, exists := seen[normalized]; !exists {
			seen[normalized] = struct{}{}
			out = append(out, normalized)
		}
		mu.Unlock()
	}

	// Respect ctx cancellation
	c.OnRequest(func(r *colly.Request) {
		select {
		case <-ctx.Done():
			r.Abort()
			return
		default:
		}

		emit(r.URL)
	})

	// Discover and follow links from common elements
	c.OnHTML(`a[href], link[href], script[src], iframe[src], form[action]`, func(e *colly.HTMLElement) {
		var raw string
		switch {
		case e.Attr("href") != "":
			raw = e.Attr("href")
		case e.Attr("src") != "":
			raw = e.Attr("src")
		case e.Attr("action") != "":
			raw = e.Attr("action")
		}
		if raw == "" {
			return
		}

		abs := e.Request.AbsoluteURL(raw)
		if abs == "" {
			return
		}

		u, err := url.Parse(abs)
		if err != nil {
			return
		}
		if u.Scheme == "" || u.Hostname() == "" {
			return
		}
		if !hostAllowed(u.Hostname(), rootHost, includeSubDomains) {
			return
		}

		// Emit the discovered URL (optional, but matches your intent of printing uniques)
		emit(u)

		// Visit it (Colly will handle max depth and duplicates at request-time on our side)
		_ = c.Visit(u.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		// Keep silent here; caller can log if needed.
		_ = r
		_ = err
	})

	if err := c.Visit(startURL.String()); err != nil {
		return nil, fmt.Errorf("visit start: %w", err)
	}

	c.Wait()

	// If caller cancelled, return ctx error but still include whatever we collected
	if ctx.Err() != nil {
		return out, ctx.Err()
	}
	return out, nil
}
