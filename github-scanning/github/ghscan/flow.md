## Flow 

#### Git commit reading phase

The tool does **not** scan the current files on disk. It scans **every commit’s diff** in the repo so it can find secrets that were ever added (even if later removed). This is implemented in `GetCommitPatches(repoDir)` in `git.go`.

| Step | Command / action | Result |
|------|------------------|--------|
| 1 | `git -C <dir> rev-list --all` | List of all commit SHAs (every commit reachable from any branch or tag). |
| 2 | For each SHA: `git -C <dir> show <sha> -p --no-color` | Full **patch** (diff) for that commit: added/removed lines, no ANSI colors. |
| 3 | Store | `[]CommitPatch{SHA, Patch}` — one struct per commit. |
| 4 | (Later) For each patch | `ScanContent(cp.Patch, patterns, commitURL)` runs regexes on the patch text; matches get the commit URL `repo.HTMLURL/commit/<SHA>`. |

- **`rev-list --all`**: enumerates which commits exist (all branches).
- **`git show <sha> -p`**: outputs that commit’s patch (the diff introduced by that commit). `--no-color` keeps output plain text.
- Commits that fail `git show` (e.g. no diff, binary) are skipped; the rest are scanned.
- Each match is tied to a **commit URL** (e.g. `https://github.com/owner/repo/commit/abc123`), not to a single file path.

So in the **git commit reading** phase the tool:

1. **Reads** which commits exist via `git rev-list --all`.
2. **Reads** each commit’s content by getting its **patch** via `git show <sha> -p`.
3. **Scans** that patch text with the loaded regexes and attaches the commit URL to any matches.

---

## Mermaid flowchart

```mermaid
flowchart TD
    A[User: ./ghscan -org=... or -repo=...] --> B[Parse flags]
    B --> C[Resolve token: -token or GITHUB_TOKEN]
    C --> D[Load regexes from YAML → compiled patterns]
    D --> E{-repo set?}
    E -->|Yes| F[Parse -repo slugs → RepoInfo list]
    F --> G[NewScan keys?, nil]
    G --> H[RunScanRepos → scanRepos]
    E -->|No| I{-org set?}
    I -->|No| J[Exit: need -org or -repo]
    I -->|Yes| K[Token required?]
    K -->|Missing| L[Exit: token required for -org]
    K -->|OK| M[Parse -org list]
    M --> N[NewScan token, orgs]
    N --> O[RunScan: for each org ListOrgRepos]
    O --> P[Collect all RepoInfo]
    P --> H
    H --> Q[For each repo:]
    Q --> R[MkdirTemp → CloneRepo]
    R --> S[GetCommitPatches: rev-list + git show -p]
    S --> T[For each commit: ScanContent patch with patterns]
    T --> U[Append ReportEntry if matches]
    U --> V[RemoveRepo clone]
    V --> W{More repos?}
    W -->|Yes| Q
    W -->|No| X[writeReport → JSON file]
    X --> Y[Result: report.json]
```

