package main


type FinalResult struct {
	Results []SecretFinding `json:"results"`
}

// SecretFinding groups a unique secret with all locations where it was found (deduplicated).
type SecretFinding struct {
	Secret string   `json:"secret"`
	Found  []string `json:"found"` // e.g. "package@1.0.0/path/to/file.js"
}

// ScanResult holds the scan results for a package
type ScanResult struct {
	PackageName string   `json:"package_name"`
	Version     string   `json:"version"`
	Secrets     []Secret `json:"secrets"`
}

type Secret struct {
	Path    string `json:"path"`
	Secrets []string
}

type versionInfo struct {
	Name string   `json:"name"`
	Version string `json:"version"`
	Dist   distInfo `json:"dist"`
}

type distInfo struct {
	Tarball string `json:"tarball"`
}

type pkgSearchResponse struct {
	DistTags       map[string]string      `json:"dist-tags"`
	Versions       map[string]versionInfo `json:"versions"`
}