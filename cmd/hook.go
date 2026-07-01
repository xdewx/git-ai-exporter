package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xdewx/git-ai-exporter/internal/git"
)

//go:embed post-commit.sh
var hookScript string

const hookSig = "git-ai-exporter post-commit hook"

func doInstallHook(r *git.Runner) error {
	hookDir, err := r.Run("rev-parse", "--git-dir")
	if err != nil {
		return fmt.Errorf("not a git repository: %w", err)
	}

	gitDir := filepath.Join(repoDir, trimNewline(hookDir))
	hooksDir := filepath.Join(gitDir, "hooks")
	hookPath := filepath.Join(hooksDir, "post-commit")

	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("create hooks dir: %w", err)
	}

	if existing, err := os.ReadFile(hookPath); err == nil {
		content := string(existing)
		if strings.Contains(content, hookSig) {
			fmt.Fprintln(os.Stderr, "Updating existing git-ai-exporter hook")
		} else {
			chainPath := hookPath + ".local"
			if err := os.WriteFile(chainPath, existing, 0755); err != nil {
				return fmt.Errorf("preserve existing hook: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Existing hook preserved: %s\n", chainPath)
		}
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
