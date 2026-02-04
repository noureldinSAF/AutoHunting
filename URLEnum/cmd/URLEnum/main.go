// this is the script of url enumeration tool
package main
import (
	"flag"
	"github.com/cyinnove/logify"
	"github.com/noureldinSAF/AutoHunting/URLEnum/internal/runner"
	//"github.com/noureldinSAF/AutoHunting/URLEnum/test"
)

func main() {
	//logify.MaxLevel = logify.Error
	options := &runner.Options{}
	flag.StringVar(&options.Domain, "d", "", "Domain to enumerate URLs for")
	flag.IntVar(&options.Timeout, "timeout", 60, "Timeout for each request in seconds")
	flag.IntVar(&options.Concurrency, "c", 5, "Number of concurrent requests")
	flag.StringVar(&options.Input, "i", "", "Input file with list of URLs")
	flag.StringVar(&options.Output, "o", "results.txt", "Output file for results")
	flag.BoolVar(&options.ActiveEnabled, "active", false, "Enable active scanning mode (crawl + headless)")
	flag.Parse()

	if err := runner.Run(options); err != nil {
		logify.Errorf("Error running URL enumeration: %v", err)
		return
	}
}
	
