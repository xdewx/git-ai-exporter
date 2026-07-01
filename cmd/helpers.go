package cmd

import (
	"path/filepath"
	"strings"

	"github.com/xdewx/git-ai-exporter/internal/git"
)

func resolvePath(p string) (string, error) {
	if filepath.IsAbs(p) {
		return p, nil
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func getCurrentBranch(r *git.Runner) (string, error) {
	out, err := r.Run("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func getOriginURL(r *git.Runner) (string, error) {
	out, err := r.Run("config", "--get", "remote.origin.url")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func extractProjectName(originURL string) string {
	url := strings.TrimSuffix(originURL, ".git")

	// URL with scheme: https://, http://, ssh://, git://
	if idx := strings.Index(url, "://"); idx >= 0 {
		url = url[idx+3:]
		if atIdx := strings.IndexByte(url, '@'); atIdx >= 0 {
			url = url[atIdx+1:]
		}
		return url
	}

	// SCP-style: [user@]host:path → host/path
	if parts := strings.SplitN(url, ":", 2); len(parts) == 2 {
		host := parts[0]
		if atIdx := strings.IndexByte(host, '@'); atIdx >= 0 {
			host = host[atIdx+1:]
		}
		return host + "/" + parts[1]
	}

	return url
}
