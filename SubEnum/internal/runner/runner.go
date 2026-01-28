package runner

import (
	"subenum/pkg/active/alterx"
	dnsprobe "subenum/pkg/active/dnsbrobe"
	"subenum/pkg/active/zonetransfer"

	"github.com/cyinnove/logify"

	// "subenum/pkg/active/bruteforce"
	"subenum/pkg/scraper"
	"subenum/pkg/scraper/sources"
	"subenum/pkg/utils"
)

func Run(opts *Options) error {

	// Handling UserInput
	// read queries
	opts.queries = utils.ExtractDomainsFromString(opts.Query)
	// check if domains read from a file
	if opts.Query == "" && opts.InputFile == "" {
		logify.Fatalf("No input or file specified")
	}

	if opts.InputFile != "" {
		var err error
		opts.queries, err = utils.ReadInputFromFile(opts.InputFile)
		if err != nil {
			logify.Fatalf("Error reading input file %s: %s", opts.InputFile, err)
		}
	}

	if len(opts.queries) == 0 {
		logify.Fatalf("No input or file specified")
	}

	// read config
	apiKeys, err := scraper.ExtractALLAPIKeys()
	if err != nil {
		logify.Fatalf("Error extracting API keys: %s", err)
	}

	var allResults map[string]bool

	allResults = map[string]bool{}
	uniqueSubdomains := []string{}

	// Get all available sources once
	client := scraper.NewSession(opts.Timeout)
	srcs := sources.GetAllSources(apiKeys)

	logify.Infof("Loaded %d source(s) for enumeration", len(srcs))

	// enumeration for each domain
	for _, q := range opts.queries {
		logify.Infof("Starting passive enumeration for: %s", q)

		queryResults := 0

		for _, src := range srcs {
			domains, err := src.Search(q, client)
			if err != nil {
				// Only log errors in verbose mode
				if opts.Verbose {
					logify.Errorf("Source %s failed: %s", src.Name(), err)
				}
				continue
			}

			for _, domain := range domains {
				if !allResults[domain] {
					logify.Silentf("%s", domain)
					allResults[domain] = true
					queryResults++
				}
			}
		}

		// Zone transfer check during passive enumeration
		if opts.ActiveEnabled {
			if zonetransfer.IsVulnerable(q, opts.Timeout) {
				logify.Infof("Zone transfer vulnerability detected for %s", q)
				zoneSubs, err := zonetransfer.GetSubdomains(q, opts.Timeout)
				if err != nil {
					logify.Errorf("Zone transfer failed for %s: %s", q, err)
				} else {
					for _, s := range zoneSubs {
						if !allResults[s] {
							logify.Silentf("%s", s)
							allResults[s] = true
							uniqueSubdomains = append(uniqueSubdomains, s)
							queryResults++
						}
					}
					logify.Infof("Zone transfer discovered %d subdomain(s) for %s", len(zoneSubs), q)
				}
			}
		}

		logify.Infof("Query %s: Found %d unique subdomain(s) (Total: %d)", q, queryResults, len(allResults))
	}

	for domain := range allResults {
		uniqueSubdomains = append(uniqueSubdomains, domain)
	}

	logify.Infof("Passive enumeration completed: Found %d unique subdomain(s)", len(uniqueSubdomains))

	if opts.ActiveEnabled {
		logify.Infof("Starting active enumeration for %d subdomain(s)", len(uniqueSubdomains))

		// 1. Permutation & Mutation
		mutationSubs, err := alterx.RunMutator(uniqueSubdomains, opts.MaxMutationsSize, opts.Enrich)
		if err != nil {
			logify.Errorf("Mutation generation failed: %s", err)
			return err
		}

		logify.Infof("Generated %d mutation(s)", len(mutationSubs))

		// 2. DNS Probing - Validate the mutated subdomains
		aliveSubdomains := dnsprobe.ProbeSubdomains(mutationSubs, opts.Timeout, opts.Concurrency)

		// 3. Add alive subdomains to results
		newActiveCount := 0
		for _, subdomain := range aliveSubdomains {
			if !allResults[subdomain] {
				logify.Silentf("%s", subdomain)
				allResults[subdomain] = true
				uniqueSubdomains = append(uniqueSubdomains, subdomain)
				newActiveCount++
			}
		}

		logify.Infof("Active enumeration completed: Found %d new subdomain(s) (Total: %d)", newActiveCount, len(allResults))
	}

	if opts.OutputFile != "" {
		if err := utils.WriteOutputToFile(opts.OutputFile, uniqueSubdomains); err != nil {
			logify.Fatalf("Error writing output file: %s", err)
		}
	}

	return nil
}
