## Overview

`dockerscan` (also called DockerScan in the code) is a command‑line tool that:

- **Connects to your local Docker daemon**
- **Pulls or inspects Docker images**
- **Starts a temporary container from each image tag**
- **Streams the filesystem from inside the container**
- **Searches text files for secrets using a broad regex pattern**
- **Outputs all findings as JSON**

This file walks learners through that flow step by step, from the CLI to the scanning engine.

---

## High‑level execution flow

1. **User runs the CLI**
   - Examples:
     - `dockerscan -i nginx`
     - `dockerscan -f images.txt`
   - The entrypoint is `cmd/dockerhunter/main.go`.

2. **Cobra parses flags and builds options**
   - Flags defined in `main.go`:
     - **`-i, --image`**: single Docker image name to scan (e.g. `nginx` or `nginx:latest`).
     - **`-f, --input-file`**: path to a file with one image per line.
     - **`-o, --output`**: JSON output file path. If omitted, a name is auto‑generated.
     - **`--debug` / `--verbose`**: control log verbosity via `logify`.
   - Values are stored in a `runner.Options` struct.

3. **Runner validates input and decides the mode**
   - The main orchestration logic lives in `pkg/runner/runner.go`.
   - `runner.Run(opts)` checks:
     - If **both** `ImageName` and `ImagesInputFile` are empty → return an error.
   - Then:
     - If `ImagesInputFile` is **set** → call `processImagesFromFile`.
     - Else → call `processSingleImage`.

4. **Runner creates a `DockerScan` client**
   - For each image, `processSingleImage` calls:
     - `client.NewDockerScan(imageName, opts.RegexesFile)`
   - `NewDockerScan` (in `pkg/client/client.go`):
     - Splits `imageName` into:
       - **`imageName`** (e.g. `nginx`)
       - **`version/tag`** (e.g. `latest`) if present.
     - Creates a Docker client using environment variables:
       - `docker.NewClientFromEnv()`
     - Builds a `DockerScan` struct that keeps:
       - Docker client
       - Image name & version
       - A slice of matches to be filled during scanning.

5. **DockerScan performs the scan**
   - The core scanning logic is in `DockerScan.Scan()`:
     1. **Determine which tags to scan**
        - If the user passed `image:tag`, only that tag is scanned (`singleVersionScan`).
        - Otherwise, it calls `fetchTags(imageName)` (`pkg/client/api.go`) to get tags from Docker Hub:
          - Makes an HTTP request to `https://registry.hub.docker.com/v2/repositories/<image>/tags`.
          - Parses the JSON response into a list of tag names.
     2. **For each tag**:
        - Build `fullImageName = "<image>:<tag>"`.
        - Try `InspectImage` to see if it already exists locally.
        - If not found, call `PullImage` to download it.
        - Re‑inspect the pulled image to get its config (especially `WorkingDir`).

6. **Container creation and scan root**
   - For each tag:
     - `createContainer(fullImageName)` makes a container that runs:
       - `/bin/sh -c "while true; do sleep 10; done"`
       - This keeps the container alive during scanning.
     - The image’s `Config.WorkingDir` is used as the base directory when available.
     - Otherwise, `scanRoot(containerID)` (in `helpers.go`) chooses a sensible default:
       - Tries `/app`, `/code`, `/usr/src/app`, `/` (in that order) inside the container.
       - It runs a small `/bin/sh` command to check which directory actually exists.

7. **Streaming files out of the container**
   - `scanContainerFiles(container.ID)` is responsible for getting file contents:
     - Uses `DownloadFromContainer` with the chosen root path to stream a **tar archive** of files.
     - Reads the tar stream entry by entry:
       - Skip non‑regular files (directories, symlinks, etc.).
       - Normalize to an absolute path like `/path/inside/container`.
       - Check against blacklists:
         - **Blacklisted directories**: e.g. `node_modules`, `.git`, build outputs, large data directories.
         - **Blacklisted extensions**: images, archives, binaries (e.g. `.png`, `.zip`, `.pdf`).
       - Load the file contents into memory.
       - Skip clearly binary content and invalid UTF‑8 to avoid noise and performance issues.

8. **Secret detection**
   - For each text file that passes filters, `scanContent(filePath, content)` is called.
   - It uses a pre‑compiled regex called `GeneralPattern` (`client/types.go`):
     - This is a **single, broad regex** designed to catch many credential‑like strings:
       - Looks for keywords like `key`, `pass`, `token`, `secret`, `admin`, `ssh`, `api`, etc.
       - Looks for common assignment patterns (`=`, `:`, `=>`, etc.).
       - Tries to filter out very short or low‑value matches (e.g. `len(match) < 20` is skipped).
   - Every match is stored as a `SecretMatch`:
     - `Secret` → the matched string.
     - `FilePath` → where it was found.

9. **Deduplicating and formatting output**
   - After scanning all relevant files and tags:
     - `RemoveDuplicateMatches()` is called to drop duplicate `(filePath, secret)` pairs.
     - `GenerateJSONOutput()` wraps all matches in a top‑level JSON object:
       - `{ "secrets": [ { "secret": "...", "file_path": "..." }, ... ] }`
       - Uses `json.MarshalIndent` to make the output human‑readable.

10. **Writing results to disk**
    - Back in the runner (`processSingleImage`):
      - If the user set `-o/--output`, that filename is used (ensuring `.json` extension).
      - Otherwise, a filename is auto‑generated based on the image and current timestamp, for example:
        - `nginx_latest_2004-01-02_15-04-05.json` (slashes are replaced with `-`).
      - The helper `saveResults` handles creating/truncating the file and writing the JSON string.

11. **Processing multiple images from a file**
    - If `-f/--input-file` is used:
      - `processImagesFromFile` reads each non‑empty line as an image name.
      - For each image:
        - It clones the `Options` struct, sets `ImageName` to that line, and calls `processSingleImage`.
      - This allows batch scanning many images in one command.

---

## Key structs and concepts (for learners)

- **`runner.Options`** (`pkg/runner/options.go`):
  - Holds all user‑supplied configuration:
    - `ImageName`, `ImagesInputFile`, `OutputFile`, `RegexesFile`, `IsFileInput`, `AllVersions`.

- **`client.DockerScan`** (`pkg/client/types.go` & `client.go`):
  - Encapsulates:
    - Docker client connection.
    - Target image name and tag.
    - Working directory inside the container.
    - In‑memory list of all secret matches.

- **`SecretMatch`**:
  - Very simple data structure:
    - `Secret` (the actual matched string).
    - `FilePath` (where it lives).

