package test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cyinnove/logify"
	"github.com/noureldinSAF/AutoHunting/URLEnum/pkg/active/crawl"
	"github.com/noureldinSAF/AutoHunting/URLEnum/pkg/active/headless"
	"github.com/noureldinSAF/AutoHunting/URLEnum/pkg/scraper"
	"github.com/noureldinSAF/AutoHunting/URLEnum/pkg/scraper/sources"
	"github.com/noureldinSAF/AutoHunting/URLEnum/pkg/utils"
)

func Run(opts *Options) error {
	opts.queries = utils.ExtractDomainsFromString(opts.Domain)

	if opts.Domain == "" && opts.Input == "" {
		logify.Fatalf("No domain or input file specified")
	}

	if opts.Input != "" {
		var err error
		opts.queries, err = utils.ReadInputFromFile(opts.Input)
		if err != nil {
			logify.Errorf("Error reading queries from file %s: %v", opts.Input, err)
			return err
		}
	}

	if len(opts.queries) == 0 {
		logify.Fatalf("No domain specified for enumeration")
	}

	if opts.Concurrency <= 0 {
		opts.Concurrency = 1
	}

	// Shared results
	allResults := make(map[string]bool)
	var mu sync.Mutex

	addURL := func(u string) {
		u = strings.TrimSpace(u)
		if u == "" {
			return
		}
		mu.Lock()
		if !allResults[u] {
			allResults[u] = true
			//logify.Silentf("Found URL: %s", u)
		}
		mu.Unlock()
	}

	// Build domain list (clean)
	domains := make([]string, 0, len(opts.queries))
	for _, q := range opts.queries {
		q = strings.TrimSpace(q)
		if q == "" || strings.HasPrefix(q, "#") {
			continue
		}
		domains = append(domains, q)
	}
	if len(domains) == 0 {
		logify.Fatalf("No valid domains found")
	}

	// Load sources
	apiKeys, err := scraper.ExtractALLAPIKeys()
	if err != nil {
		logify.Fatalf("Failed to read config: %v", err)
	}
	srcs := sources.GetAllSources(apiKeys)

	// Pick the 2 passive sources explicitly by name
	var webarchiveSrc scraper.Source
	var commoncrawlSrc scraper.Source
	for _, s := range srcs {
		switch strings.ToLower(s.Name()) {
		case "webarchive":
			webarchiveSrc = s
		case "commoncrawl":
			commoncrawlSrc = s
		}
	}
	if webarchiveSrc == nil {
		logify.Errorf("webarchive source not found in sources.GetAllSources")
	}
	if commoncrawlSrc == nil {
		logify.Errorf("commoncrawl source not found in sources.GetAllSources")
	}

	// ---- distribute workers across 4 pipelines ----
	// You can tune these ratios; they must sum <= opts.Concurrency.
	// Example split (26): 6+6+7+7 = 26
	webWorkers := max(1, opts.Concurrency/4)
	ccWorkers := max(1, opts.Concurrency/4)
	crawlWorkers := max(1, (opts.Concurrency-webWorkers-ccWorkers)/2)
	headlessWorkers := opts.Concurrency - webWorkers - ccWorkers - crawlWorkers
	if headlessWorkers < 1 {
		headlessWorkers = 1
		// Adjust down crawlWorkers if needed
		if crawlWorkers > 1 {
			crawlWorkers--
		}
	}

	logify.Infof("Worker split => webarchive=%d, commoncrawl=%d, crawl=%d, headless=%d",
		webWorkers, ccWorkers, crawlWorkers, headlessWorkers,
	)

	// Worker delay (your requirement)
	workerDelay := 10 * time.Second

	// Per-task timeout
	perTaskTimeout := time.Duration(opts.Timeout) * time.Second
	if perTaskTimeout <= 0 {
		perTaskTimeout = 60 * time.Second
	}

	// Domain job fanout: each pipeline needs to process all domains,
	// so we create 4 separate channels (broadcast).
	webJobs := make(chan string, len(domains))
	ccJobs := make(chan string, len(domains))
	crawlJobs := make(chan string, len(domains))
	headJobs := make(chan string, len(domains))

	for _, d := range domains {
		webJobs <- d
		ccJobs <- d
		crawlJobs <- d
		headJobs <- d
	}
	close(webJobs)
	close(ccJobs)
	close(crawlJobs)
	close(headJobs)

	var wg sync.WaitGroup

	// ---- Passive: webarchive ----
	if webarchiveSrc != nil {
		for i := 0; i < webWorkers; i++ {
			wg.Add(1)
			id := i + 1
			go func(workerID int) {
				defer wg.Done()
				client := scraper.NewSession(opts.Timeout)

				for d := range webJobs {
					ctx, cancel := context.WithTimeout(context.Background(), perTaskTimeout)
					urls, err := webarchiveSrc.Search(ctx, d, client)
					cancel()

					if err != nil {
						logify.Errorf("[WEB-W%02d] webarchive failed for %s: %v", workerID, d, err)
					} else {
						for _, u := range urls {
							addURL(u)
						}
					}

					time.Sleep(workerDelay)
				}
			}(id)
		}
	}

	// ---- Passive: commoncrawl ----
	if commoncrawlSrc != nil {
		for i := 0; i < ccWorkers; i++ {
			wg.Add(1)
			id := i + 1
			go func(workerID int) {
				defer wg.Done()
				client := scraper.NewSession(opts.Timeout)

				for d := range ccJobs {
					ctx, cancel := context.WithTimeout(context.Background(), perTaskTimeout)
					urls, err := commoncrawlSrc.Search(ctx, d, client)
					cancel()

					if err != nil {
						logify.Errorf("[CC-W%02d] commoncrawl failed for %s: %v", workerID, d, err)
					} else {
						for _, u := range urls {
							addURL(u)
						}
					}

					time.Sleep(workerDelay)
				}
			}(id)
		}
	}

	// ---- Active: crawl ----
	if opts.ActiveEnabled {
		crawlOpts := &crawl.Options{
			MaxDepth:    2,
			Parallelism: max(1, opts.Concurrency/2), // inside colly itself
			Timeout:     perTaskTimeout,
			AllowQuery:  true,
		}

		for i := 0; i < crawlWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for d := range crawlJobs {
					seed := normalizeSeed(d)

					ctx, cancel := context.WithTimeout(context.Background(), perTaskTimeout)
					found, err := crawl.Enumerate(ctx, seed, opts.IncludeSubdomains, crawlOpts)
					cancel()

					if err != nil {
						logify.Errorf("[CRAWL] failed for %s: %v", seed, err)
					} else {
						for _, u := range found {
							addURL(u)
						}
					}

					time.Sleep(workerDelay)
				}
			}()
		}

		// ---- Active: headless ----
		headOpts := headless.Options{
			Concurrency:   max(1, opts.Concurrency/2),
			Timeout:       perTaskTimeout,
			Wait:          8 * time.Second,
			ChromePath:    "/usr/bin/google-chrome",
			Headless:      true,
			NoSandbox:     true,
			DisableGPU:    true,
			DisableDevShm: true,
		}

		for i := 0; i < headlessWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for d := range headJobs {
					seed := normalizeSeed(d)

					ctx, cancel := context.WithTimeout(context.Background(), perTaskTimeout)
					found, err := headless.Enumerate(ctx, seed, opts.IncludeSubdomains, headOpts)
					cancel()

					if err != nil {
						logify.Errorf("[HEADLESS] failed for %s: %v", seed, err)
					} else {
						for _, u := range found {
							addURL(u)
						}
					}

					time.Sleep(workerDelay)
				}
			}()
		}
	} else {
		// If not active, just drain active job channels quickly
		_ = crawlJobs
		_ = headJobs
	}

	wg.Wait()

	// Build output list
	uniqueURLs := make([]string, 0, len(allResults))
	for u := range allResults {
		uniqueURLs = append(uniqueURLs, u)
	}

	logify.Infof("Done. Total unique URLs: %d", len(uniqueURLs))

	if opts.Output != "" {
		if err := utils.WriteOutputToFile(opts.Output, uniqueURLs); err != nil {
			logify.Fatalf("Error writing output to file %s: %v", opts.Output, err)
		}
	}

	return nil
}

func normalizeSeed(domainOrURL string) string {
	s := strings.TrimSpace(domainOrURL)
	if s == "" {
		return ""
	}
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return s
	}
	return fmt.Sprintf("https://%s", s)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
