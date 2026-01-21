package sources

import (
	"github.com/noureldinSAF/AutoHunting/pkg/scraper"
	"github.com/noureldinSAF/AutoHunting/pkg/scraper/sources/crtsh"
	"github.com/noureldinSAF/AutoHunting/pkg/scraper/sources/whoisfreaks"
	"github.com/noureldinSAF/AutoHunting/pkg/scraper/sources/whoisxmlapi"
)

var AllSources = [...]scraper.Source{
	&crtsh.Source{},
	&whoisfreaks.Source{},
	&whoisxmlapi.Source{},
}

func GetAllSources(apiKeys map[string][]string) []scraper.Source {

	var sources []scraper.Source

	for _, source := range AllSources {
		sourceName := source.Name()

		if source.RequireAPIKey(){
			if keys, ok := apiKeys[sourceName]; ok && len(keys) > 0 {
				switch sourceName {
				case "whoisxmlapi":
					sources = append(sources, whoisxmlapi.New(keys))
				case "whoisfreaks":
					sources = append(sources, whoisfreaks.New(keys))
				}
			
			}
		} else {
			sources = append(sources, source)
		}
		}

		return sources
}