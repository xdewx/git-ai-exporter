# git-ai-exporter

Parse [git-ai](https://github.com/xdewx/git-ai) commit notes and export structured statistics. Can output JSON or push directly to a [git-ai-dashboard](https://github.com/xdewx/git-ai-dashboard) instance.

## Install

Download the latest binary for your platform from [Releases](../../releases), or build from source:

```bash
go install github.com/xdewx/git-ai-exporter@latest
```

## Usage

```bash
# Parse current repo, output JSON
git-ai-exporter

# Pretty-print the last 50 commits
git-ai-exporter -n 50 --output pretty

# Specify repo path
git-ai-exporter -r /path/to/repo

# Filter by date range
git-ai-exporter --since "2026-01-01" --until "2026-06-30"

# Parse and push to dashboard
git-ai-exporter --push --url https://dash.example.com/api/collect --token your-token
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-r, --repo` | `.` | Git repository path |
| `-b, --branch` | current | Target branch |
| `-n, --count` | 0 (all) | Max commits to parse |
| `--since` | | Start date |
| `--until` | | End date |
| `--output` | `json` | Output format: `json` or `pretty` |
| `--push` | false | Push results to dashboard |
| `--url` | | Dashboard collect API URL |
| `--token` | | API token |
| `--hostname` | hostname | Client hostname identifier |

## Build from source

```bash
git clone git@github.com:xdewx/git-ai-exporter.git
cd git-ai-exporter
go build -o dist/git-ai-exporter .
```

Cross-compile for all platforms:

```bash
make build-all
```

Output:

- `dist/git-ai-exporter-linux-amd64`
- `dist/git-ai-exporter-darwin-amd64`
- `dist/git-ai-exporter-darwin-arm64` (Apple Silicon)
- `dist/git-ai-exporter-windows-amd64.exe`

## How it works

1. Runs `git log --notes=ai` to fetch commits with git-ai notes
2. Parses each note to extract AI-generated code regions (file + line ranges)
3. Calculates AI additions (sum of line ranges) and total additions (via `--numstat`)
4. Extracts tool/model info from note JSON sessions
5. Outputs structured data compatible with `POST /api/collect`

## Output format

```json
{
  "originUrl": "git@github.com:user/repo.git",
  "projectName": "repo",
  "branch": "main",
  "hostname": "my-machine",
  "commits": [
    {
      "sha": "abc123...",
      "author": "Author",
      "authorEmail": "author@example.com",
      "message": "feat: ...",
      "committedAt": "2026-07-01T10:00:00+08:00",
      "humanAdditions": 80,
      "aiAdditions": 20,
      "aiAccepted": 20,
      "diffAddedLines": 100,
      "diffDeletedLines": 5,
      "toolBreakdown": {
        "openai/gpt-4": {
          "aiAdditions": 20,
          "aiAccepted": 20
        }
      }
    }
  ]
}
```
