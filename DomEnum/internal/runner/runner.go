package runner

import (
	"github.com/noureldinSAF/AutoHunting/DomEnum/pkg/scraper"
	"github.com/noureldinSAF/AutoHunting/DomEnum/pkg/scraper/sources"
	"github.com/noureldinSAF/AutoHunting/DomEnum/pkg/utils"

	"github.com/cyinnove/logify"
)

func Run(opts *Options) error {

	opts.queries = utils.ExtractDomainsFromString(opts.Query)

	if opts.InputFile == "" && opts.Query == "" {
		logify.Fatalf("Please specify  a domain or an input file.")
	}

	if opts.InputFile != "" {
		var err error
		opts.queries, err = utils.ReadInputFromFile(opts.InputFile)
		if err != nil {
			logify.Errorf("Error reading queries from file %s: %v", opts.InputFile, err)
			return err
		}
	}

	if len(opts.queries) == 0 {
		logify.Fatalf("No domains specified")
	}
	apiKeys, err := scraper.ExtractAllAPIKeys()
	if err != nil {
		logify.Fatalf("Error reading API keys: %v", err)
	}

	var allResults map[string]bool
	allResults = make(map[string]bool)

	for _, q := range opts.queries {
		client := scraper.NewSession(opts.Timeout)
		srcs := sources.GetAllSources(apiKeys)
		for _, src := range srcs {
			domains, err := src.Search(q, client)
			if err != nil {
				logify.Errorf("Error searching %s: %s", q, err)
			}
			for _, domain := range domains {
				if !allResults[domain] {
					logify.Silentf("%s", domain)
					allResults[domain] = true
				}
			}

		}
	}

	uniqueDomains := []string{}
	for domain := range allResults {
		uniqueDomains = append(uniqueDomains, domain)
	}

	if opts.OutputFile != "" {
		if err := utils.WriteOutputToFile(opts.OutputFile, uniqueDomains); err != nil {
			logify.Fatalf("Error writing output to file %s: %v", opts.OutputFile, err)
		}
	}

	return nil
}
