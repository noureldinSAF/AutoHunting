package runner

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/miekg/dns"

	"github.com/zomaxsec/dnsenum/pkg/dnsprobe"	
	"bufio"
	"sort"
)

func convertRecordsToUint16(records []string) []uint16 {
	if len(records) == 0 {
		return []uint16{}
	}

	for _, r := range records {
		r = strings.ToUpper(r)
	}

	recordMap := map[string]uint16{
		"A":     dns.TypeA,
		"AAAA":  dns.TypeAAAA,
		"CNAME": dns.TypeCNAME,
		"MX":    dns.TypeMX,
		"NS":    dns.TypeNS,
		"TXT":   dns.TypeTXT,
		"PTR":   dns.TypePTR,
		"SRV":   dns.TypeSRV,
		"SOA":   dns.TypeSOA,
	}

	var result []uint16
	for _, r := range records {
		upperR := strings.ToUpper(r)
		if rt, ok := recordMap[upperR]; ok {
			result = append(result, rt)
		}
	}

	return result
}

type OutputJSON struct {
	Subdomains []*SubdomainResult `json:"subdomains"`
}

type SubdomainResult struct {
	Host             string              `json:"host"`
	Records          map[string][]string `json:"records"`
	ResolversChecked int                 `json:"resolvers_checked"`
}

func writeJSONOutput(filename string, results []*dnsprobe.ProbeResult) error {
	// Define the initial output structure.
	output := OutputJSON{
		Subdomains: []*SubdomainResult{},
	}

	for _, result := range results {
		if len(result.Results) == 0 {
			continue
		}

		// Map to collect non-empty answers with status "NOERROR", key is the record type.
		records := make(map[string][]string)
		for _, r := range result.Results {
			// Only include records with successful "NOERROR" DNS responses and responses with values.
			if r.DNStatus == "NOERROR" && len(r.Values) > 0 {
				records[r.RecordName] = r.Values
			}
		}

		// Only include hosts that have at least one valid record.
		if len(records) > 0 {
			output.Subdomains = append(output.Subdomains, &SubdomainResult{
				Host:             result.Host, // Hostname or subdomain probed
				Records:          records,     // Map of record type to values
				ResolversChecked: 0,           // Placeholder (update if used)
			})
		}
	}

	// Convert to pretty indented JSON.
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON result to the target file with standard permissions.
	return os.WriteFile(filename, data, 0644)
}

func writeLinesOutput(path string, results []*dnsprobe.ProbeResult) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()

    seen := make(map[string]struct{})
    out := make([]string, 0, len(results))

    for _, pr := range results {
        if pr == nil {
            continue
        }

        // ✅ Adjust this line depending on your ProbeResult struct:
        // Common names might be: pr.Target, pr.Domain, pr.Hostname, pr.Input, etc.
        target := pr.Host

        if target == "" {
            continue
        }

        valid := false
        for _, r := range pr.Results {
            if r.DNStatus == "NOERROR" && len(r.Values) > 0 {
                valid = true
                break
            }
        }

        if !valid {
            continue
        }

        if _, ok := seen[target]; ok {
            continue
        }
        seen[target] = struct{}{}
        out = append(out, target)
    }

    sort.Strings(out)

    w := bufio.NewWriter(f)
    for _, t := range out {
        if _, err := w.WriteString(t + "\n"); err != nil {
            return err
        }
    }
    return w.Flush()
}
