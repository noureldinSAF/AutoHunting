package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)


func LoadWordlist(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}

func TrimContentType(s string) string {
	if idx := strings.Index(s, ";"); idx >= 0 {
		return strings.TrimSpace(s[:idx])
	}
	return strings.TrimSpace(s)
}

// formatNucleiStyle returns a finding line in Nuclei-style format:
// https://example.com/path [cl:200] [ct:text/html] [cl:2132]
func FormatNucleiStyle(url string, status int, contentType string, contentLength int64) string {
	return fmt.Sprintf("%s [cl:%d] [ct:%s] [cl:%d]", url, status, contentType, contentLength)
}

// parseStatusFilter parses a comma-separated list of status codes (e.g. "200,301,302").
// Returns nil for empty string (no filter).
func ParseStatusFilter(s string) ([]int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	var out []int
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid status code %q: %w", p, err)
		}
		out = append(out, n)
	}
	return out, nil
}


func BuildURL(baseURL, path string) string {
	baseURL = strings.TrimSpace(baseURL)
	path = strings.TrimSpace(path)
	if baseURL == "" {
		return ""
	}
	if path == "" {
		return strings.TrimSuffix(baseURL, "/")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return strings.TrimSuffix(baseURL, "/") + path
}