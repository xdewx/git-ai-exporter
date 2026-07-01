# git-ai-exporter Design

## Overview

Go CLI tool that parses git-ai commit notes and outputs structured data, with optional push to git-ai-dashboard's `/api/collect`.

## CLI Interface

```
git-ai-exporter [flags]
```

Flags:

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-r, --repo` | string | `.` | Git repo path |
| `-b, --branch` | string | current | Target branch |
| `-n, --count` | int | 0 (all) | Max commits |
| `--since` | string | | Start date |
| `--until` | string | | End date |
| `--output` | string | `json` | Output format: `json` or `pretty` |
| `--push` | bool | false | Push to dashboard |
| `--url` | string | | Dashboard collect API URL |
| `--token` | string | | API token |
| `--hostname` | string | `os.Hostname()` | Client hostname |

## Architecture

```
main.go → cmd/root.go (cobra)
                ↓
           internal/git/
           ├── git.go       (exec.Command wrapper)
           ├── log.go       (git log --format, parse commit metadata + notes)
           └── numstat.go   (git log --numstat, parse per-commit diff stats)
                ↓
           internal/parser/
           ├── types.go     (CommitData, NoteEntry, SessionData)
           ├── note.go      (parseNote — prologue + JSON sessions)
           └── tool.go      (getToolModel — extract tool/model)
                ↓
           internal/output/formatter.go (JSON / pretty print)
                ↓
           internal/reporter/collect.go (HTTP POST)
```

## Data Flow

1. `git log --format=<fmt>` → parse commit metadata + note content
2. Parse each note → entries (file/line ranges) + sessions (tool/model)
3. Calculate aiAdditions = sum(lineEnd - lineStart + 1)
4. `git log --numstat --format=""` → parse per-commit diff stats
5. Calculate totalAdditions, humanAdditions = total - ai
6. Extract tool/model from sessions → build toolBreakdown
7. Output JSON or push via HTTP

## Parsing Logic

### Note Format

```
file1.rb
  def foo 1-5
  def bar 8-12
file2.py
  class Baz 1-20
---
{"sessions": {"session_id": {"agent_id": {"tool": "openai", "model": "gpt-4"}}}}
```

Prologue is indented lines under a file header. Separator `---`. JSON part contains sessions with tool/model info.

### Git Format

```
▣<sha>▣<author>▣<email>▣<message>▣<date>▣<note>⟐
```

Fields delimited by `▣`, records by `⟐`.

## Cross-Compilation

```makefile
build-all:
  GOOS=linux   GOARCH=amd64 go build -o dist/git-ai-exporter-linux-amd64
  GOOS=darwin  GOARCH=amd64 go build -o dist/git-ai-exporter-darwin-amd64
  GOOS=darwin  GOARCH=arm64 go build -o dist/git-ai-exporter-darwin-arm64
  GOOS=windows GOARCH=amd64 go build -o dist/git-ai-exporter-windows-amd64.exe
```
