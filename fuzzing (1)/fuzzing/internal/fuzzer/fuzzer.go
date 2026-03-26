package fuzzer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"fuzzing/internal/fingerprint"
	"fuzzing/internal/utils"
)

// Config holds fuzzer options.
type Config struct {
	Target          string
	Wordlist        string
	DelayMs         int
	Timeout         int
	BaselineDelayMs int // Delay in ms between the 3 baseline requests.
	Method          string
	StatusFilter    []int // If non-empty, only report findings with these status codes.
	StatusFilterStr string
}

// Match represents an interesting response (differs from baseline).
type Match struct {
	URL           string
	StatusCode    int
	ContentLength int64
}

// Run loads the wordlist, captures the baseline fingerprint, fuzzes each path,
// compares to baseline, and prints findings.
func Run(cfg Config) error {
	words, err := utils.LoadWordlist(cfg.Wordlist)
	if err != nil {
		return fmt.Errorf("wordlist: %w", err)
	}


	cfg.StatusFilter , _ = utils.ParseStatusFilter(cfg.StatusFilterStr)
	
	client := &http.Client{
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
		Transport: &http.Transport{MaxIdleConnsPerHost: 10},
	}

	baseURL := strings.TrimSuffix(cfg.Target, "/")
	baseline, err := fingerprint.CaptureRange(client, baseURL, cfg.Method, cfg.BaselineDelayMs)
	if err != nil {
		return fmt.Errorf("baseline fingerprint: %w", err)
	}

	fmt.Printf("[*] Baseline range: status %d-%d, length %d-%d bytes (3 samples)\n\n",
		baseline.StatusCodeMin, baseline.StatusCodeMax,
		baseline.ContentLengthMin, baseline.ContentLengthMax)

	var findings []*Match
	for i, path := range words {
		url := utils.BuildURL(baseURL, path)
		if url == "" {
			continue
		}

		req, err := http.NewRequest(cfg.Method, url, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] %s: %v\n", path, err)
			continue
		}
		
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] %s: %v\n", path, err)
			continue
		}

		fp := fingerprint.FromResponse(resp)
		if fp.ContentLength == -1 {
			body, _ := io.ReadAll(resp.Body)
			fp.ContentLength = int64(len(body))
		}
		defer resp.Body.Close()

		if !baseline.InRange(fp) && cfg.statusAllowed(fp.StatusCode) {

			findings = append(findings, &Match{
				URL:           url,
				StatusCode:    fp.StatusCode,
				ContentLength: fp.ContentLength,
			})
			fmt.Println(utils.FormatNucleiStyle(url, fp.StatusCode, fp.Headers["Content-Type"], fp.ContentLength))
		}

		if cfg.DelayMs > 0 && i < len(words)-1 {
			time.Sleep(time.Duration(cfg.DelayMs) * time.Millisecond)
		}
	}

	fmt.Printf("\n[*] Done. %d interesting, %d total\n", len(findings), len(words))
	return nil
}

// statusAllowed returns true if no filter is set or status is in the allowed list.
func (cfg *Config) statusAllowed(status int) bool {
	if len(cfg.StatusFilter) == 0 {
		return true
	}
	for _, s := range cfg.StatusFilter {
		if s == status {
			return true
		}
	}
	return false
}

