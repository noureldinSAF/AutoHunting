package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Scan downloads the latest version of an npm package, scans its files for secrets, and returns results (deduplicated by secret).
func Scan(pkgName string) (*FinalResult, error) {
	meta, err := fetchPackageMetadata(pkgName)
	if err != nil {
		return nil, err
	}

	version := "latest"
	if v, ok := meta.DistTags["latest"]; ok {
		version = v
	}
	info, ok := meta.Versions[version]
	if !ok {
		return nil, fmt.Errorf("version %q not found for package %s", version, pkgName)
	}

	result, err := scanVersion(info)
	if err != nil {
		return nil, err
	}

	results := []ScanResult{result}
	return resultsBySecret(results), nil
}

func scanVersion(info versionInfo) (ScanResult, error) {
	tmpDir, err := os.MkdirTemp("", "npm-scan")
	if err != nil {
		return ScanResult{}, err
	}
	defer os.RemoveAll(tmpDir)

	resp, err := http.Get(info.Dist.Tarball)
	if err != nil {
		return ScanResult{}, fmt.Errorf("download: %v", err)
	}
	defer resp.Body.Close()

	tgzPath := filepath.Join(tmpDir, "pkg.tgz")
	f, err := os.Create(tgzPath)
	if err != nil {
		return ScanResult{}, err
	}
	_, err = io.Copy(f, resp.Body)
	f.Close()
	if err != nil {
		return ScanResult{}, err
	}

	secrets, err := unpackAndScan(tgzPath)
	if err != nil {
		return ScanResult{}, err
	}

	return ScanResult{
		PackageName: info.Name,
		Version:     info.Version,
		Secrets:     secrets,
	}, nil
}

func (fr *FinalResult) JSON() (string, error) {
	data, err := json.MarshalIndent(fr, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	pkgName := "eslint-plugin-no-secrets"
	if len(os.Args) > 1 {
		pkgName = os.Args[1]
	}

	result, err := Scan(pkgName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	out, err := result.JSON()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.WriteFile("results.json", []byte(out), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Results written to npm.json")
}
