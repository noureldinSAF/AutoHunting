package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var blacklistedFiles = []string{
	"node_modules",
	"package.json",
	"package-lock.json",
	"yarn.lock",
	"LICENSE",
}

func isBlacklisted(blacklist []string, value string) bool {
	value = strings.TrimSpace(strings.ToLower(value))

	for _, item := range blacklist {
		if value == strings.ToLower(strings.TrimSpace(item)) {
			return true
		}
	}

	return false
}
// isPathBlacklisted returns true if the full path contains any blacklisted file or directory (any path segment).
func isPathBlacklisted(blacklist []string, fullPath string) bool {
	segments := strings.Split(filepath.ToSlash(fullPath), "/")
	for _, seg := range segments {
		if seg == "" {
			continue
		}
		if isBlacklisted(blacklist, seg) {
			return true
		}
	}
	return false
}

// fetchPackageMetadata fetches the metadata of an npm package.
func fetchPackageMetadata(pkgName string) (*pkgSearchResponse, error) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s", pkgName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch package metadata: %v", err)
	}
	defer resp.Body.Close()

	var pkgData pkgSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&pkgData); err != nil {
		return nil, fmt.Errorf("failed to decode package metadata: %v", err)
	}

	return &pkgData, nil
}

// unpackAndScan unpacks a .tgz file and scans the files for secrets using the provided regex patterns.
func unpackAndScan(tgzFile string) ([]Secret, error) {
	file, err := os.Open(tgzFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open .tgz file: %v", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	secrets := make([]Secret, 0)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar header: %v", err)
		}

		if isPathBlacklisted(blacklistedFiles, header.Name) {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// skip
		case tar.TypeReg:
			fileSecrets, err := scanTarFile(tr)
			if err != nil {
				return nil, err
			}
			if len(fileSecrets) == 0 {
				continue
			}

			filteredSecrets := make([]string, 0)
			for _, match := range fileSecrets {
				if match == "" || len(match) < 25 || len(match) > 1500 {
					continue
				}
				if FalsePositiveFilter.MatchString(match) {
					continue
				}
				filteredSecrets = append(filteredSecrets, match)
			}

			if len(filteredSecrets) == 0 {
				continue
			}

			secrets = append(secrets, Secret{
				Path:    header.Name,
				Secrets: filteredSecrets,
			})
		}
	}

	return secrets, nil
}

// scanTarFile scans a file in the tar archive for secrets using the provided regex patterns.
func scanTarFile(r io.Reader) ([]string, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
	}
	matches := secretPatterns.FindAllString(string(content), -1)
	return matches, nil
}

// resultsBySecret converts ScanResults into a deduplicated list keyed by secret value, with locations in "found".
func resultsBySecret(results []ScanResult) *FinalResult {
	bySecret := make(map[string]map[string]struct{})
	for _, r := range results {
		locPrefix := r.PackageName + "@" + r.Version
		for _, s := range r.Secrets {
			for _, secretVal := range s.Secrets {
				if secretVal == "" {
					continue
				}
				if bySecret[secretVal] == nil {
					bySecret[secretVal] = make(map[string]struct{})
				}
				bySecret[secretVal][locPrefix+"/"+s.Path] = struct{}{}
			}
		}
	}
	out := make([]SecretFinding, 0, len(bySecret))
	for secret, locSet := range bySecret {
		found := make([]string, 0, len(locSet))
		for loc := range locSet {
			found = append(found, loc)
		}
		out = append(out, SecretFinding{Secret: secret, Found: found})
	}
	return &FinalResult{Results: out}
}
