package runner

import (
	"strings"

	"github.com/cyinnove/logify"
	"github.com/noureldinSAF/AutoHunting/DomEnum/pkg/scraper/sources"
	"github.com/noureldinSAF/AutoHunting/URL_Enum_v2/pkg/scraper"
	"github.com/noureldinSAF/AutoHunting/URL_Enum_v2/utils"
)

func Run(opts *Options) error {
	// Core logic for URL enumeration would go here.
	opts.queries = utils.ExtractDomainsFromString(opts.Domain)

	if opts.Domain != "" && opts.Input != "" {
		logify.Fatalf("No input of file specified") 
        }
	if opts.Input != "" {
		logify.Infof("Reading URLs from input file: %s", opts.Input)
	}

	if len(opts.queries) == 0 {
		logify.Fatalf("No domain specified for enumeration")
	}

	apiKeys, err := scraper.ExtractALLAPIKeys()
	if err != nil {
		logify.Fatalf("Failed to read config: %v", err)
	}

	allResults := map[string]bool{}

	for _, q := range opts.queries {
		
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}
		logify.Infof("Starting enumeration for domain: %s via passive sources", q)

		client := scraper.NewSession(opts.Timeout)

		srcs := sources.GetAllSources(apiKeys)

		for _, src := range srcs {
			urls, err := src.Search(q, client)
			if err != nil {
				logify.Errorf("Source %s failed: %v", src.Name(), err)
				continue
			}
			for _, u := range urls {
				u = strings.Trim(u)
				if u == "" {
					continue
				}
				if !allResults[u] {
					allResults[u] = true
					logify.Silentf("Found URL: %s", u)
				}

	}

}

