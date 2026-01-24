package main

import (
	"github.com/noureldinSAF/AutoHunting/DomEnum/internal/runner"
	"flag"
	"github.com/cyinnove/logify"
)

var options = &runner.Options{}

// create a main function for flags
func main() {
	flag.StringVar(&options.Query, "q", "", "Domain to enumerate")
	flag.IntVar(&options.Timeout, "t", 10, "Timeout for enumeration ( seconds )")
	flag.IntVar(&options.Concurrency, "c", 5, "Concurrency level for enumeration")
	flag.StringVar(&options.InputFile, "i", "", "input file")
	flag.StringVar(&options.OutputFile, "o", "", "output file")
	flag.BoolVar(&options.ActiveEnabled, "active", false, "Active Enumeration" )

	flag.Parse()

	if err := runner.Run(options); err != nil {
		logify.Errorf("Error Running the tool %v", err)
	}

}