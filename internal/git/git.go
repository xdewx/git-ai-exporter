package git

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Runner struct {
	Dir string
}

func NewRunner(dir string) *Runner {
	return &Runner{Dir: dir}
}

func (r *Runner) Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", &GitError{
				Stderr: strings.TrimSpace(string(ee.Stderr)),
				Args:   args,
			}
		}
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return string(out), nil
}

type GitError struct {
	Stderr string
	Args   []string
}

func (e *GitError) Error() string {
	return fmt.Sprintf("git %s: %s", strings.Join(e.Args, " "), e.Stderr)
}

func (r *Runner) CheckDaemon() error {
	_, err := r.Run("show-ref", "refs/notes/ai")
	if err != nil {
		return fmt.Errorf("git-ai daemon not running (no AI notes ref found).\n" +
			"Make sure git-ai is running before committing:\n" +
			"  git-ai bg start")
	}

	if !processRunning("git-ai") {
		return fmt.Errorf("git-ai daemon process not found.\n" +
			"It may have exited unexpectedly. Start it with:\n" +
			"  git-ai bg start\n" +
			"Check status:\n" +
			"  git-ai bg status")
	}

	return nil
}

func processRunning(name string) bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-NoProfile", "-Command",
			"Get-Process '"+name+"' -ErrorAction SilentlyContinue")
		return cmd.Run() == nil
	}
	out, err := exec.Command("pgrep", "-x", name).Output()
	return err == nil && len(out) > 0
}
