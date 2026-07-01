//go:build windows

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/v4/process"
)

const (
	detachedProcess = 0x00000008
	createNoWindow  = 0x08000000
)

func runGuard() error {
	fmt.Fprintln(log.Writer(), "Running git-ai-exporter guard in foreground (Ctrl+C to stop)")
	p := &guardProgram{}
	p.run()
	return nil
}

func doInstallGuard() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	key := `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
	name := guardServiceName()
	value := fmt.Sprintf(`"%s" --guard`, exe)

	cmd := exec.Command("reg", "add", key, "/v", name, "/t", "REG_SZ", "/d", value, "/f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("reg add: %s: %w", string(out), err)
	}
	fmt.Fprintf(log.Writer(), "Guard auto-start added to Registry: HKCU\\...\\Run\\%s\n", name)

	if pids := findGuardProcesses(); len(pids) > 0 {
		fmt.Fprintf(log.Writer(), "Guard already running (PID: %v), skipping\n", pids)
		return nil
	}

	startCmd := exec.Command(exe, "--guard")
	startCmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: detachedProcess | createNoWindow,
	}
	if err := startCmd.Start(); err != nil {
		fmt.Fprintf(log.Writer(), "Warning: failed to start guard process: %v\n", err)
	} else {
		fmt.Fprintf(log.Writer(), "Guard process started (PID: %d)\n", startCmd.Process.Pid)
	}
	return nil
}

func doUninstallGuard() error {
	key := `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
	name := guardServiceName()

	cmd := exec.Command("reg", "delete", key, "/v", name, "/f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("reg delete: %s: %w", string(out), err)
	}
	fmt.Fprintln(log.Writer(), "Guard auto-start removed from Registry")

	killGuardProcesses()
	return nil
}

func findGuardProcesses() []int32 {
	myPid := os.Getpid()
	procs, err := process.Processes()
	if err != nil {
		return nil
	}
	var pids []int32
	for _, p := range procs {
		if p.Pid == int32(myPid) {
			continue
		}
		name, err := p.Name()
		if err != nil {
			continue
		}
		name = strings.TrimSuffix(name, ".exe")
		if name != "git-ai-exporter" {
			continue
		}
		cmdline, err := p.Cmdline()
		if err != nil || !strings.Contains(cmdline, "--guard") {
			continue
		}
		pids = append(pids, p.Pid)
	}
	return pids
}

func killGuardProcesses() {
	for _, pid := range findGuardProcesses() {
		p, err := process.NewProcess(pid)
		if err != nil {
			continue
		}
		if err := p.Kill(); err != nil {
			fmt.Fprintf(log.Writer(), "Warning: failed to kill guard (PID %d): %v\n", pid, err)
		} else {
			fmt.Fprintf(log.Writer(), "Guard process stopped (PID: %d)\n", pid)
		}
	}
}


