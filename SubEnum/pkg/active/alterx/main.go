package alterx

import (
	"context"

	"github.com/projectdiscovery/alterx"
)

func RunMutator(subdomains []string, maxSize int, enrich bool) ([]string, error) {
	result := []string{}

	opts := &alterx.Options{
		Domains: subdomains,
		MaxSize: maxSize,
		Enrich:  enrich,
	}

	m, err := alterx.New(opts)
	if err != nil {
		return nil, err
	}

	subdomainsCh := m.Execute(context.Background())

	// Apply count-based limit if maxSize > 0
	// Note: alterx's MaxSize limits data size in bytes, not count
	// So we need to limit by count here
	count := 0
	for subdomain := range subdomainsCh {
		// If maxSize > 0, limit the number of results
		if maxSize > 0 && count >= maxSize {
			break
		}
		result = append(result, subdomain)
		count++
	}

	return result, nil

}
