package runner

import (
	"context"
	"strings"
	"sync"
	"time"
    //"math/rand"
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

	apiKeys, err := scraper.ExtractALLAPIKeys()
	if err != nil {
		logify.Fatalf("Failed to read config: %v", err)
	}

	// Shared results across all workers (dedupe)
	allResults := make(map[string]bool)
	var resultsMu sync.Mutex

	addResult := func(u string) {
		u = strings.TrimSpace(u)
		if u == "" {
			return
		}
		resultsMu.Lock()
		if !allResults[u] {
			allResults[u] = true
			//logify.Infof("Found URL: %s", u)
		}
		resultsMu.Unlock()
	}

	// =========================
	// Passive Enumeration (concurrent over queries)
	// =========================
	
	srcs := sources.GetAllSources(apiKeys)

	jobs := make(chan string, len(opts.queries))
	var wg sync.WaitGroup

	workerCount := opts.Concurrency
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		workerID := i + 1

		go func(id int) {
			defer wg.Done()

			// One client per worker (safe + avoids shared state surprises)
			client := scraper.NewSession(opts.Timeout)

			for q := range jobs {
				q = strings.TrimSpace(q)
				if q == "" {
					continue
				}

				logify.Infof("[W%02d] Starting enumeration for domain: %s via passive sources", id, q)

				for _, src := range srcs {
					// Per-source timeout controlled here (NO defer inside loop)
					ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opts.Timeout)*time.Second)
					urls, err := src.Search(ctx, q, client)
					cancel()

					if err != nil {
						if src.Name() == "commoncrawl" {
							continue
						}
						logify.Errorf("[W%02d] Source %s failed for %s: %v", id, src.Name(), q, err)
						continue
					}

					for _, u := range urls {
						addResult(u)
					}
				}

				resultsMu.Lock()
				total := len(allResults)
				resultsMu.Unlock()

				logify.Infof("[W%02d] Passive enumeration done for %s. Total unique URLs so far: %d", id, q, total)
				time.Sleep(10 * time.Second) // brief pause between queries
			}
		}(workerID)
	}

	// Enqueue all queries
	for _, q := range opts.queries {
		q = strings.TrimSpace(q)
		if q == "" || strings.HasPrefix(q, "#") {
			continue
		}
		jobs <- q
	}
	close(jobs)

	wg.Wait()

	// Collect unique URLs after passive stage
	uniqueURLs := make([]string, 0, len(allResults))
	for u := range allResults {
		uniqueURLs = append(uniqueURLs, u)
	}

	// =========================
	// Active Enumeration (concurrent over seeds)
	// =========================
	if opts.ActiveEnabled {
		logify.Infof("Started Active Enumeration for %d seed URL(s)", len(uniqueURLs))

		perTargetTimeout := time.Duration(opts.Timeout) * time.Second
		if perTargetTimeout <= 0 {
			perTargetTimeout = 30 * time.Second
		}

		ctx, cancel := context.WithTimeout(context.Background(), perTargetTimeout)
		defer cancel()

		// 1) Crawl tool (Colly)
		crawlOpts := &crawl.Options{
			MaxDepth:    2,
			Parallelism: opts.Concurrency,
			Timeout:     perTargetTimeout,
			AllowQuery:  true,
		}

		// Run crawl concurrently over seeds (bounded)
		{
			seeds := append([]string(nil), uniqueURLs...)
			seedJobs := make(chan string, len(seeds))
			var seedWG sync.WaitGroup

			for i := 0; i < opts.Concurrency; i++ {
				seedWG.Add(1)
				go func() {
					defer seedWG.Done()
					for seed := range seedJobs {
						select {
						case <-ctx.Done():
							return
						default:
						}
						found, err := crawl.Enumerate(ctx, seed, opts.IncludeSubdomains, crawlOpts)
						if err != nil {
							logify.Errorf("crawl failed for %s: %v", seed, err)
							continue
						}
						for _, u := range found {
							addResult(u)
						}
					}
				}()
			}

			for _, s := range seeds {
				seedJobs <- s
			}
			close(seedJobs)
			seedWG.Wait()
		}

		// Refresh seeds after crawl added more
		resultsMu.Lock()
		uniqueURLs = uniqueURLs[:0]
		for u := range allResults {
			uniqueURLs = append(uniqueURLs, u)
		}
		resultsMu.Unlock()

		// 2) Headless tool (chromedp)
		headlessOpts := headless.Options{
			Concurrency:   opts.Concurrency,
			Timeout:       perTargetTimeout,
			Wait:          8 * time.Second,
			ChromePath:    "/usr/bin/google-chrome",
			Headless:      true,
			NoSandbox:     true,
			DisableGPU:    true,
			DisableDevShm: true,
		}

		// Run headless concurrently over seeds (bounded)
		{
			seeds := append([]string(nil), uniqueURLs...)
			seedJobs := make(chan string, len(seeds))
			var seedWG sync.WaitGroup

			for i := 0; i < opts.Concurrency; i++ {
				seedWG.Add(1)
				go func() {
					defer seedWG.Done()
					for seed := range seedJobs {
						select {
						case <-ctx.Done():
							return
						default:
						}
						found, err := headless.Enumerate(ctx, seed, opts.IncludeSubdomains, headlessOpts)
						if err != nil {
							logify.Errorf("headless failed for %s: %v", seed, err)
							continue
						}
						for _, u := range found {
							addResult(u)
						}
					}
				}()
			}

			for _, s := range seeds {
				seedJobs <- s
			}
			close(seedJobs)
			seedWG.Wait()
		}

		resultsMu.Lock()
		total := len(allResults)
		resultsMu.Unlock()

		logify.Infof("Active Enumeration finished. Total unique URLs: %d", total)

		// Final unique list from map
		uniqueURLs = make([]string, 0, len(allResults))
		for u := range allResults {
			uniqueURLs = append(uniqueURLs, u)
		}
	}

	// Output
	if opts.Output != "" {
		if err := utils.WriteOutputToFile(opts.Output, uniqueURLs); err != nil {
			logify.Fatalf("Error writing output to file %s: %v", opts.Output, err)
		}
	}

	return nil
}
