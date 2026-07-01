package git

import (
	"fmt"
	"os"
	"os/exec"
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
