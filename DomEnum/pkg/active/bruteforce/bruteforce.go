// bruteforce/bruteforce.go
package bruteforce

import (
	"fmt"
	"strings"
	"github.com/cyinnove/tldify"
)

var TLDs = []string{
	"com","net","org","io","ai","co","me","dev","app","xyz",
	"uk","us","de","fr","ca","au","in","jp","cn","sa","ae",
}

func GenerateWordList(domains []string) ([]string, error) {
	final := make([]string, 0, len(domains)*len(TLDs))
	seen := make(map[string]struct{}, len(domains)*len(TLDs))

	for _, d := range domains {
		d = strings.TrimSpace(strings.TrimSuffix(d, "."))
		if d == "" || d == "." {
			continue
		}

		pd, err := tldify.Parse(d)
		if err != nil || pd.Domain == "" {
			// don’t fail the whole run; just skip bad inputs
			continue
		}

		for _, t := range TLDs {
			cand := fmt.Sprintf("%s.%s", pd.Domain, t)
			if _, ok := seen[cand]; ok {
				continue
			}
			seen[cand] = struct{}{}
			final = append(final, cand)
		}
	}

	return final, nil
}