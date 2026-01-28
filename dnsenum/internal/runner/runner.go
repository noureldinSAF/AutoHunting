package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cyinnove/logify"

	"github.com/zomaxsec/dnsenum/pkg/dnsprobe"
	"github.com/zomaxsec/dnsenum/pkg/utils"
)

func Run(opts *Options) error {

	// Read Targets from user input

	if opts.TargetsFile != "" {
		var err error
		opts.Targets, err = utils.ReadInputFromFile(opts.TargetsFile)
		if err != nil {
			logify.Fatalf("Failed to read targets file: %v", err)
		}
	}

	if utils.IsStdin() {
		var err error
		opts.Targets, err = utils.ReadInputFromStdin()
		if err != nil {
			logify.Fatalf("Failed to read targets file: %v", err)
		}

	}

	if len(opts.Targets) == 0 {
		logify.Fatalf("No targets specified")
	}

	// Prepare custom rsolvers list if exist

	if opts.CustomResolversFile != "" {
		var err error
		opts.customResolvers, err = utils.ReadInputFromFile(opts.CustomResolversFile)
		if err != nil {
			logify.Fatalf("Failed to read custom resolvers file: %v", err)
		}
	}

	records := convertRecordsToUint16(opts.Records)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logify.Infof("Interrupted, shutting down gracefully...")
		cancel()
	}()

	results := []*dnsprobe.ProbeResult{}

	// Concurrency control	
	clevel := make(chan struct{}, opts.Concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var processed int64
	total := len(opts.Targets)

	logify.Infof("Starting DNS enumeration for %d targets", total)

	for _, target := range opts.Targets {
		select {
		case <-ctx.Done():
			goto done
		default:
		}

		wg.Add(1)

		go func(target string) {
			defer wg.Done()

			clevel <- struct{}{}
			defer func() { <-clevel }()

			result := dnsprobe.RunProbe(opts.Timeout, target, records, opts.customResolvers, opts.Strategy)
			mu.Lock()
			results = append(results, result)
			processed++
		for _, r := range result.Results {
			if r.DNStatus == "NOERROR" && len(r.Values) > 0 {
				logify.Silentf("%s", target)
			}
		}
			mu.Unlock()
		}(target)
	}

done:
	wg.Wait()

	validCount := 0
	for _, r := range results {
		if len(r.Results) > 0 {
			validCount++
		}
	}

	logify.Infof("Completed: %d/%d targets have valid DNS records", validCount, total)

	if opts.OutputFile != "" {
    var err error
    switch opts.OutputFormat {
    case "json":
        err = writeJSONOutput(opts.OutputFile, results)
    case "lines":
        err = writeLinesOutput(opts.OutputFile, results)
    default:
        return fmt.Errorf("invalid output format: %s", opts.OutputFormat)
    }

    if err != nil {
        return fmt.Errorf("failed to write output file: %w", err)
    }
    logify.Infof("Results written to %s", opts.OutputFile)
}


	return nil
}
