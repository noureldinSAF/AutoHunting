package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cyinnove/logify"
)

// CommitPatch holds the SHA and full patch (diff) of a single commit.
type CommitPatch struct {
	SHA   string
	Patch string
}

// CloneRepo clones the repository into dir.
func CloneRepo(cloneURL, dir string) error {
	cmd := exec.Command("git", "clone", "--quiet", cloneURL, dir)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone: %w", err)
	}
	return nil
}

// GetCommitPatches returns all commit patches (diffs) in the repo, including added and deleted lines.
func GetCommitPatches(repoDir string) ([]CommitPatch, error) {
	// Get all commit SHAs (all branches)
	revList := exec.Command("git", "-C", repoDir, "rev-list", "--all")
	revList.Stderr = nil
	out, err := revList.Output()
	if err != nil {
		return nil, fmt.Errorf("git rev-list: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	shas := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			shas = append(shas, line)
		}
	}

	if len(shas) == 0 {
		return nil, nil
	}

	patches := make([]CommitPatch, 0, len(shas))
	for _, sha := range shas {
		show := exec.Command("git", "-C", repoDir, "show", sha, "-p", "--no-color")
		show.Stderr = nil
		patch, err := show.Output()
		if err != nil {
			// Some commits (e.g. no diff, or binary) may fail; skip
			continue
		}
		patches = append(patches, CommitPatch{SHA: sha, Patch: string(patch)})
	}

	return patches, nil
}

// RemoveRepo deletes the cloned repo directory.
func RemoveRepo(dir string) error {
	if dir == "" {
		return nil
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(abs); err != nil {
		logify.Warningf("remove repo dir %s: %v", abs, err)
		return err
	}
	return nil
}
