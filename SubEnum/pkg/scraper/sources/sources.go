package sources

import (
	"subenum/pkg/scraper"
	"subenum/pkg/scraper/sources/abuseipdb"
	"subenum/pkg/scraper/sources/alienvault"
	"subenum/pkg/scraper/sources/anubis"
	"subenum/pkg/scraper/sources/arpsyndicate"
	"subenum/pkg/scraper/sources/bevigil"
	"subenum/pkg/scraper/sources/binaryedge"
	"subenum/pkg/scraper/sources/bufferover"
	"subenum/pkg/scraper/sources/builtwith"
	"subenum/pkg/scraper/sources/c99"
	"subenum/pkg/scraper/sources/censys"
	"subenum/pkg/scraper/sources/certspotter"
	"subenum/pkg/scraper/sources/chaos"
	"subenum/pkg/scraper/sources/chinaz"
	"subenum/pkg/scraper/sources/coderog"
	"subenum/pkg/scraper/sources/commoncrawl"
	"subenum/pkg/scraper/sources/crtsh"
	"subenum/pkg/scraper/sources/cyfare"
	"subenum/pkg/scraper/sources/digitalyama"
	"subenum/pkg/scraper/sources/dnsdb"
	"subenum/pkg/scraper/sources/dnsdumpster"
	"subenum/pkg/scraper/sources/dnsrepo"
	"subenum/pkg/scraper/sources/driftnet"
	"subenum/pkg/scraper/sources/fofa"
	"subenum/pkg/scraper/sources/fullhunt"
	"subenum/pkg/scraper/sources/github"
	"subenum/pkg/scraper/sources/gitlab"
	"subenum/pkg/scraper/sources/google"
	"subenum/pkg/scraper/sources/hackertarget"
	"subenum/pkg/scraper/sources/hudsonrock"
	"subenum/pkg/scraper/sources/huntermap"
	"subenum/pkg/scraper/sources/intelx"
	"subenum/pkg/scraper/sources/leakix"
	"subenum/pkg/scraper/sources/merklemap"
	"subenum/pkg/scraper/sources/myssl"
	"subenum/pkg/scraper/sources/netlas"
	"subenum/pkg/scraper/sources/odin"
	"subenum/pkg/scraper/sources/onyphe"
	"subenum/pkg/scraper/sources/pugrecon"
	"subenum/pkg/scraper/sources/quake"
	"subenum/pkg/scraper/sources/racent"
	"subenum/pkg/scraper/sources/rapidapi"
	"subenum/pkg/scraper/sources/rapiddns"
	"subenum/pkg/scraper/sources/rapidfinder"
	"subenum/pkg/scraper/sources/rapidscan"
	"subenum/pkg/scraper/sources/reconcloud"
	"subenum/pkg/scraper/sources/redhuntlabs"
	"subenum/pkg/scraper/sources/riddler"
	"subenum/pkg/scraper/sources/robtex"
	"subenum/pkg/scraper/sources/rsecloud"
	"subenum/pkg/scraper/sources/securitytrails"
	"subenum/pkg/scraper/sources/shodan"
	"subenum/pkg/scraper/sources/shodanx"
	"subenum/pkg/scraper/sources/shrewdeye"
	"subenum/pkg/scraper/sources/sitedossier"
	"subenum/pkg/scraper/sources/threatbook"
	"subenum/pkg/scraper/sources/threatcrowd"
	"subenum/pkg/scraper/sources/threatminer"
	"subenum/pkg/scraper/sources/urlscan"
	"subenum/pkg/scraper/sources/virustotal"
	"subenum/pkg/scraper/sources/waybackarchive"
	"subenum/pkg/scraper/sources/whoisxmlapi"
	"subenum/pkg/scraper/sources/windvane"
	"subenum/pkg/scraper/sources/zoomeyeapi"
)

var AllSources = [...]scraper.Source{
	&abuseipdb.Source{},
	&alienvault.Source{},
	&anubis.Source{},
	&arpsyndicate.Source{},
	&bevigil.Source{},
	&binaryedge.Source{},
	&bufferover.Source{},
	&builtwith.Source{},
	&c99.Source{},
	&censys.Source{},
	&certspotter.Source{},
	&chaos.Source{},
	&chinaz.Source{},
	&coderog.Source{},
	&commoncrawl.Source{},
	&crtsh.Source{},
	&cyfare.Source{},
	&digitalyama.Source{},
	&dnsdb.Source{},
	&dnsdumpster.Source{},
	&dnsrepo.Source{},
	&driftnet.Source{},
	&fofa.Source{},
	&fullhunt.Source{},
	&github.Source{},
	&gitlab.Source{},
	&google.Source{},
	&hackertarget.Source{},
	&hudsonrock.Source{},
	&huntermap.Source{},
	&intelx.Source{},
	&leakix.Source{},
	&merklemap.Source{},
	&myssl.Source{},
	&netlas.Source{},
	&odin.Source{},
	&onyphe.Source{},
	&pugrecon.Source{},
	&quake.Source{},
	&racent.Source{},
	&rapidapi.Source{},
	&rapiddns.Source{},
	&rapidfinder.Source{},
	&rapidscan.Source{},
	&reconcloud.Source{},
	&redhuntlabs.Source{},
	&riddler.Source{},
	&robtex.Source{},
	&rsecloud.Source{},
	&securitytrails.Source{},
	&shodan.Source{},
	&shodanx.Source{},
	&shrewdeye.Source{},
	&sitedossier.Source{},
	&threatbook.Source{},
	&threatcrowd.Source{},
	&threatminer.Source{},
	&urlscan.Source{},
	&virustotal.Source{},
	&waybackarchive.Source{},
	&whoisxmlapi.Source{},
	&windvane.Source{},
	&zoomeyeapi.Source{},
}

// GetAllSources returns all available sources based on provided API keys
func GetAllSources(apiKeys map[string][]string) []scraper.Source {
	var sources []scraper.Source

	// Loop through all registered sources
	for _, source := range AllSources {
		sourceName := source.Name()

		// If source requires API key, check if we have keys for it
		if source.RequiresAPIKey() {
			if keys, ok := apiKeys[sourceName]; ok && len(keys) > 0 {
				switch sourceName {
				case "arpsyndicate":
					sources = append(sources, arpsyndicate.New(keys))
				case "bevigil":
					sources = append(sources, bevigil.New(keys))
				case "binaryedge":
					sources = append(sources, binaryedge.New(keys))
				case "bufferover":
					sources = append(sources, bufferover.New(keys))
				case "builtwith":
					sources = append(sources, builtwith.New(keys))
				case "c99":
					sources = append(sources, c99.New(keys))
				case "censys":
					sources = append(sources, censys.New(keys))
				case "chaos":
					sources = append(sources, chaos.New(keys))
				case "chinaz":
					sources = append(sources, chinaz.New(keys))
				case "coderog":
					sources = append(sources, coderog.New(keys))
				case "digitalyama":
					sources = append(sources, digitalyama.New(keys))
				case "dnsdb":
					sources = append(sources, dnsdb.New(keys))
				case "dnsrepo":
					sources = append(sources, dnsrepo.New(keys))
				case "driftnet":
					sources = append(sources, driftnet.New(keys))
				case "fofa":
					sources = append(sources, fofa.New(keys))
				case "fullhunt":
					sources = append(sources, fullhunt.New(keys))
				case "github":
					sources = append(sources, github.New(keys))
				case "gitlab":
					sources = append(sources, gitlab.New(keys))
				case "google":
					sources = append(sources, google.New(keys))
				case "huntermap":
					sources = append(sources, huntermap.New(keys))
				case "intelx":
					sources = append(sources, intelx.New(keys))
				case "leakix":
					sources = append(sources, leakix.New(keys))
				case "merklemap":
					sources = append(sources, merklemap.New(keys))
				case "netlas":
					sources = append(sources, netlas.New(keys))
				case "odin":
					sources = append(sources, odin.New(keys))
				case "onyphe":
					sources = append(sources, onyphe.New(keys))
				case "pugrecon":
					sources = append(sources, pugrecon.New(keys))
				case "quake":
					sources = append(sources, quake.New(keys))
				case "rapidapi":
					sources = append(sources, rapidapi.New(keys))
				case "rapidfinder":
					sources = append(sources, rapidfinder.New(keys))
				case "rapidscan":
					sources = append(sources, rapidscan.New(keys))
				case "redhuntlabs":
					sources = append(sources, redhuntlabs.New(keys))
				case "robtex":
					sources = append(sources, robtex.New(keys))
				case "rsecloud":
					sources = append(sources, rsecloud.New(keys))
				case "securitytrails":
					sources = append(sources, securitytrails.New(keys))
				case "shodan":
					sources = append(sources, shodan.New(keys))
				case "threatbook":
					sources = append(sources, threatbook.New(keys))
				case "virustotal":
					sources = append(sources, virustotal.New(keys))
				case "whoisxmlapi":
					sources = append(sources, whoisxmlapi.New(keys))
				case "windvane":
					sources = append(sources, windvane.New(keys))
				case "zoomeyeapi":
					sources = append(sources, zoomeyeapi.New(keys))
				}
			}
		} else {
			// Source doesn't require API key, add it directly
			sources = append(sources, source)
		}
	}

	return sources
}
