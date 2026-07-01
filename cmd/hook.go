package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xdewx/git-ai-exporter/internal/git"
)

//go:embed post-commit.sh
var hookScript string

func doInstallHook(r *git.Runner) error {
	hookDir, err := r.Run("rev-parse", "--git-dir")
	if err != nil {
		return fmt.Errorf("not a git repository: %w", err)
	}

	gitDir := filepath.Join(repoDir, trimNewline(hookDir))
	hookPath := filepath.Join(gitDir, "hooks", "post-commit")

	if err := os.MkdirAll(filepath.Dir(hookPath), 0755); err != nil {
		return fmt.Errorf("create hooks dir: %w", err)
	}

	if err := os.WriteFile(hookPath, []byte(hookScript), 0755); err != nil {
		return fmt.Errorf("write hook: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Hook installed: %s\n\n", hookPath)
	fmt.Fprintln(os.Stderr, "Configure your dashboard:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  git config hooks.ai-exporter-url https://your-dashboard.com/api/collect")
	fmt.Fprintln(os.Stderr, "  git config hooks.ai-exporter-token your-api-token")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Optional:")
	fmt.Fprintln(os.Stderr, "  git config hooks.ai-exporter-count 1     # commits per push (default: 1)")
	fmt.Fprintln(os.Stderr, "  git config hooks.ai-exporter-hostname    # override hostname")

	return nil
}

func trimNewline(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return s
}
