//go:build !windows

package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

func doDetach() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable: %w", err)
	}

	var args []string
	for _, a := range os.Args[1:] {
		if a == "--detach" || strings.HasPrefix(a, "--detach=") {
			continue
		}
		args = append(args, a)
	}

	nullIn, err := os.OpenFile(os.DevNull, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("open null device: %w", err)
	}
	defer nullIn.Close()

	attr := &os.ProcAttr{
		Dir:   ".",
		Env:   os.Environ(),
		Files: []*os.File{nullIn, os.Stdout, os.Stderr},
		Sys: &syscall.SysProcAttr{
			Setsid: true,
		},
	}

	proc, err := os.StartProcess(exe, append([]string{exe}, args...), attr)
	if err != nil {
		return fmt.Errorf("start background process: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Exporter running in background (pid: %d)\n", proc.Pid)
	return nil
}
