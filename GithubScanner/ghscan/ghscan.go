package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cyinnove/logify"
)

// NewScan initializes a new scan process.
// API keys are required when scanning an org (-org) so the GitHub API can list repos.
// When only scanning specific repos (-repo), keys are optional: public repos clone without auth.
func NewScan(apiKeys []string, target []string) (*Scan, error) {
	var validKeys []string
	for _, key := range apiKeys {
		if key == "" {
			continue
		}
		if isAliveKey(key) {
			validKeys = append(validKeys, key)
		} else {
			fmt.Printf("key is not valid => %s\n", key)
		}
	}

	// When scanning only specific repos (target nil), no API is used — token optional for public repos.
	if len(validKeys) == 0 && len(target) == 0 {
		return &Scan{
			APIKeys:        nil,
			Target:         nil,
			Errorfs:        make(chan error),
		}, nil
	}

	if len(validKeys) == 0 {
		return nil, errors.New("all the api keys are invalid, please add valid ones (required for -org)")
	}

	return &Scan{
		APIKeys:        validKeys,
		Target:         target,
		CurrentUsedKey: validKeys[0],
		Errorfs:        make(chan error),
	}, nil
}

// scanRepo scans a single repo and returns its report (if any matches found).
// Duplicates (same rule + matched content) are filtered within the same repo.
func scanRepo(repo RepoInfo, patterns []CompiledPattern) *RepoReport {
	dir, err := os.MkdirTemp("", "ghscan-*")
	if err != nil {
		logify.Errorf("temp dir: %v", err)
		return nil
	}

	logify.Infof("[+] Cloning %s", repo.FullName)
	if err := CloneRepo(repo.CloneURL, dir); err != nil {
		logify.Warningf("clone %s: %v", repo.FullName, err)
		_ = RemoveRepo(dir)
		return nil
	}

	patches, err := GetCommitPatches(dir)
	if err != nil {
		logify.Warningf("commit patches %s: %v", repo.FullName, err)
		_ = RemoveRepo(dir)
		return nil
	}

	// Deduplicate matches at repo level: rule + matched content -> first commit URL
	seen := make(map[string]map[string]string) // rule -> matched text -> commit URL
	var allMatches []MatchResult

	for _, cp := range patches {
		commitURL := repo.HTMLURL + "/commit/" + cp.SHA
		matches := ScanContent(cp.Patch, patterns, commitURL)

		for _, m := range matches {
			if seen[m.Rule] == nil {
				seen[m.Rule] = make(map[string]string)
			}
			if _, exists := seen[m.Rule][m.Matched]; exists {
				// Duplicate: skip this match (first occurrence kept)
				continue
			}
			seen[m.Rule][m.Matched] = m.CommitURL
			allMatches = append(allMatches, m)
		}
	}

	if err := RemoveRepo(dir); err != nil {
		logify.Warningf("remove clone %s: %v", repo.FullName, err)
	}
	logify.Infof("[+] Done %s (deleted clone)", repo.FullName)

	if len(allMatches) == 0 {
		return nil
	}

	return &RepoReport{
		RepoURL: repo.HTMLURL,
		Matches: allMatches,
	}
}

// RunScan lists repos for each org, then scans each repo.
// Returns reports grouped by organization.
func (s *Scan) RunScan(patterns []CompiledPattern) ([]OrgReport, error) {
	var reports []OrgReport

	for _, org := range s.Target {
		logify.Infof("[+] Listing repos for org %s", org)
		orgRepos, err := s.ListOrgRepos(org)
		if err != nil {
			logify.Errorf("list repos for %s: %v", org, err)
			continue
		}
		logify.Infof("[+] Found %d repos in %s", len(orgRepos), org)

		var repoReports []RepoReport
		for _, repo := range orgRepos {
			if report := scanRepo(repo, patterns); report != nil {
				repoReports = append(repoReports, *report)
			}
		}

		if len(repoReports) > 0 {
			reports = append(reports, OrgReport{
				Org:   org,
				Repos: repoReports,
			})
		}
	}

	return reports, nil
}

// RunScanRepos scans only the given repos (no org listing).
// Groups results by repo owner.
func (s *Scan) RunScanRepos(repos []RepoInfo, patterns []CompiledPattern) ([]OrgReport, error) {
	// Group repos by owner (org)
	ownerRepos := make(map[string][]RepoInfo)
	for _, repo := range repos {
		parts := strings.Split(repo.FullName, "/")
		if len(parts) < 2 {
			continue
		}
		owner := parts[0]
		ownerRepos[owner] = append(ownerRepos[owner], repo)
	}

	var reports []OrgReport
	for owner, ownerRepoList := range ownerRepos {
		var repoReports []RepoReport
		for _, repo := range ownerRepoList {
			if report := scanRepo(repo, patterns); report != nil {
				repoReports = append(repoReports, *report)
			}
		}
		if len(repoReports) > 0 {
			reports = append(reports, OrgReport{
				Org:   owner,
				Repos: repoReports,
			})
		}
	}

	return reports, nil
}

// ScanContent runs all compiled patterns on content and returns matches.
func ScanContent(content string, patterns []CompiledPattern, commitURL string) []MatchResult {
	var matches []MatchResult
	seen := make(map[string]map[string]bool) // rule -> matched text -> true

	for _, p := range patterns {
		for _, sub := range p.Regex.FindAllString(content, -1) {
			if len(strings.TrimSpace(sub)) < 20 {
				continue
			}
			if seen[p.Name] == nil {
				seen[p.Name] = make(map[string]bool)
			}
			if seen[p.Name][sub] {
				continue
			}
			seen[p.Name][sub] = true
			matches = append(matches, MatchResult{
				Rule:      p.Name,
				Matched:   sub,
				CommitURL: commitURL,
			})
		}
	}
	return matches
}

func main() {
	var (
		orgFlag     = flag.String("org", "", "Organization(s) to scan (comma-separated). Example: -org=myorg or -org=org1,org2")
		repoFlag    = flag.String("repo", "", "Specific repo(s) to scan as owner/name (comma-separated). Overrides -org. Example: -repo=owner/repo or -repo=org/r1,org/r2")
		tokenFlag   = flag.String("token", "", "GitHub token (or GITHUB_TOKEN). Required for -org; optional for -repo (public repos clone without auth).")
		regexesFlag = flag.String("regexes", "regexes.yaml", "Path to regexes YAML file")
		outFlag     = flag.String("out", "report.json", "Output report JSON path")
	)
	flag.Parse()

	token := *tokenFlag
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	patterns, err := LoadRegexes(*regexesFlag)
	if err != nil {
		logify.Errorf("load regexes: %v", err)
		os.Exit(1)
	}
	
	logify.Infof("Loaded %d regex patterns from %s", len(patterns), *regexesFlag)

	// -repo takes precedence: scan only those repos (token optional for public repos)
	if *repoFlag != "" {
		runRepoScan(*repoFlag, token, patterns, *outFlag)
		return
	}

	// -org: list all repos in org(s) then scan (token required for API)
	if *orgFlag == "" {
		logify.Errorf("either -org or -repo is required. Example: -org=myorg or -repo=owner/repo")
		os.Exit(1)
	}
	if token == "" {
		logify.Errorf("token required for -org: use -token=TOKEN or set GITHUB_TOKEN (needed to list org repos)")
		os.Exit(1)
	}

	runOrgScan(*orgFlag, token, patterns, *outFlag)
}

func runRepoScan(repoFlag, token string, patterns []CompiledPattern, outPath string) {
	repos, err := parseRepoSlugs(repoFlag)
	if err != nil {
		logify.Errorf("parse -repo: %v", err)
		os.Exit(1)
	}

	logify.Infof("Scanning %d repo(s)", len(repos))
	if token == "" {
		logify.Infof("No token set; public repos will clone without auth (private repos will fail)")
	}

	keys := []string{}
	if token != "" {
		keys = []string{token}
	}

	scan, err := NewScan(keys, nil)
	if err != nil {
		logify.Errorf("%v", err)
		os.Exit(1)
	}

	report, err := scan.RunScanRepos(repos, patterns)
	if err != nil {
		logify.Errorf("scan: %v", err)
		os.Exit(1)
	}

	writeReport(report, outPath)
}

func runOrgScan(orgFlag, token string, patterns []CompiledPattern, outPath string) {
	orgs := parseList(orgFlag)
	logify.Infof("Scanning org(s): %s", strings.Join(orgs, ", "))

	scan, err := NewScan([]string{token}, orgs)
	if err != nil {
		logify.Errorf("%v", err)
		os.Exit(1)
	}

	report, err := scan.RunScan(patterns)
	if err != nil {
		logify.Errorf("scan: %v", err)
		os.Exit(1)
	}

	writeReport(report, outPath)
}

func parseList(s string) []string {
	var out []string
	for _, v := range strings.Split(s, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func parseRepoSlugs(s string) ([]RepoInfo, error) {
	var repos []RepoInfo
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		info, err := RepoFromSlug(part)
		if err != nil {
			return nil, err
		}
		repos = append(repos, info)
	}
	return repos, nil
}

func writeReport(report []OrgReport, path string) {
	totalMatches := 0
	for _, org := range report {
		for _, repo := range org.Repos {
			totalMatches += len(repo.Matches)
		}
	}

	logify.Infof("Found %d matches across %d org(s)", totalMatches, len(report))
	if totalMatches == 0 {
		logify.Infof("No matches. Report not written.")
		return
	}

	f, err := os.Create(path)
	if err != nil {
		logify.Errorf("create report: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		logify.Errorf("write report: %v", err)
		os.Exit(1)
	}

	logify.Infof("Report written to %s", path)
}
