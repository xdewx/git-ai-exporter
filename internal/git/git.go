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
	return RunCmd(r.Dir, "git", args...)
}

func RunCmd(dir, bin string, args ...string) (string, error) {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", &CmdError{
				Stderr: strings.TrimSpace(string(ee.Stderr)),
				Bin:    bin,
				Args:   args,
			}
		}
		return "", fmt.Errorf("%s %s: %w", bin, strings.Join(args, " "), err)
	}
	return string(out), nil
}

type CmdError struct {
	Stderr string
	Bin    string
	Args   []string
}

func (e *CmdError) Error() string {
	return fmt.Sprintf("%s %s: %s", e.Bin, strings.Join(e.Args, " "), e.Stderr)
}

func (r *Runner) CheckDaemon() error {
	_, err := r.Run("show-ref", "refs/notes/ai")
	if err != nil {
		return fmt.Errorf("git-ai daemon not started yet (no AI notes found).\n" +
			"Start it with:\n" +
			"  git-ai bg start")
	}

	if processRunning("git-ai") {
		return nil
	}

	fmt.Fprintln(os.Stderr, "git-ai daemon not running, starting...")
	if _, err := RunCmd(r.Dir, "git-ai", "bg", "start"); err != nil {
		return fmt.Errorf("failed to start git-ai daemon:\n" +
			"  %s\n" +
			"Try manually:\n" +
			"  git-ai bg start", strings.TrimSpace(err.Error()))
	}

	if !processRunning("git-ai") {
		return fmt.Errorf("git-ai daemon failed to start.\n" +
			"Try manually:\n" +
			"  git-ai bg start\n" +
			"Check status:\n" +
			"  git-ai bg status")
	}

	fmt.Fprintln(os.Stderr, "git-ai daemon started")
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
