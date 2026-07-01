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
	url := originURL
	url = strings.TrimSuffix(url, ".git")

	if idx := strings.LastIndexByte(url, '/'); idx >= 0 {
		url = url[idx+1:]
	}
	if idx := strings.LastIndexByte(url, ':'); idx >= 0 {
		url = url[idx+1:]
	}

	return url
}
