package sources

import (
	"github.com/noureldinSAF/AutoHunting/URLEnum/pkg/scraper"
	"github.com/noureldinSAF/AutoHunting/URLEnum/pkg/scraper/sources/commoncrawl"
	"github.com/noureldinSAF/AutoHunting/URLEnum/pkg/scraper/sources/urlscan"
	"github.com/noureldinSAF/AutoHunting/URLEnum/pkg/scraper/sources/webarchive"
)

var AllSources = [...]scraper.Source{
	&commoncrawl.Source{},
	&urlscan.Source{},
	&webarchive.Source{},
}

func GetAllSources(apiKeys map[string][]string) []scraper.Source {

	var sources []scraper.Source

	for _, source := range AllSources {
		//sourceName := source.Name()

		if source.RequireAPIKey() {
			// if keys, ok := apiKeys[sourceName]; ok && len(keys) > 0 {
			// 	switch sourceName {
			// 	case "urlscan":
			// 		sources = append(sources, urlscan.New(keys))
			//     }
			continue;
		} else { sources = append(sources, source) }
	    }
    

	return sources
}
 