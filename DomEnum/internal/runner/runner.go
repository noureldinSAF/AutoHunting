package runner

import (
	"github.com/noureldinSAF/AutoHunting/DomEnum/pkg/active/bruteforce"
	"github.com/noureldinSAF/AutoHunting/DomEnum/pkg/active/dnsprobe"
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

	cleanQ := make([]string, 0, len(opts.queries))
    for _, q := range opts.queries {
	nq, ok := utils.NormalizeDomain(q)
	if ok {
		cleanQ = append(cleanQ, nq)
	 }
    }  
    opts.queries = cleanQ


	apiKeys, err := scraper.ExtractAllAPIKeys()
	if err != nil {
		logify.Fatalf("Error reading API keys: %v", err)
	}

	allResults := map[string]bool{}


	for _, q := range opts.queries {

		logify.Infof("Started Enumeration for query: %s", q)
		client := scraper.NewSession(opts.Timeout)
		srcs := sources.GetAllSources(apiKeys)
		for _, src := range srcs {
			domains, err := src.Search(q, client)
			if err != nil {
				logify.Errorf("Error searching %s: %s", q, err)
			}
			for _, domain := range domains {
	           nd, ok := utils.NormalizeDomain(domain)
	        if !ok {
		       continue
	        }
	        if !allResults[nd] {
		    logify.Silentf("%s", nd)
		    allResults[nd] = true
	        }
}

		}
		logify.Infof("Passive Enumeration found %d domain(s)", len(allResults))
	}

	uniqueDomains := []string{}


	for domain := range allResults {
		uniqueDomains = append(uniqueDomains, domain)
	}

	if opts.ActiveEnabled {

		logify.Infof("Started Active Enumeration for %d domain(s)", len(uniqueDomains))
		allWordList, err := bruteforce.GenerateWordList(uniqueDomains)
		if err != nil {
			logify.Fatalf("Error generating bruteforce wordlist: %s", err)
		}
		aliveDomains := dnsprobe.CheckDomainsWithConcurrency(allWordList, opts.Concurrency)

		for _,domain := range aliveDomains {
			if !allResults[domain] {
				allResults[domain] = true
				uniqueDomains = append(uniqueDomains, domain )
			}
		}
		logify.Infof("Ative Enumeration found %d", len(uniqueDomains))
	}


	if opts.OutputFile != "" {
		if err := utils.WriteOutputToFile(opts.OutputFile, uniqueDomains); err != nil {
			logify.Fatalf("Error writing output to file %s: %v", opts.OutputFile, err)
		}
	}

	return nil
}
