package git

import (
	"os"
	"testing"

	"github.com/shirou/gopsutil/v4/process"
)

func TestProcessRunning_NonExistent(t *testing.T) {
	if processRunning("this-process-should-not-exist-x7k9m2") {
		t.Fatal("expected false for non-existent process")
	}
}

func TestProcessNamesAreString(t *testing.T) {
	procs, err := process.Processes()
	if err != nil {
		t.Fatal(err)
	}
	if len(procs) == 0 {
		t.Fatal("expected at least one process")
	}

	found := false
	pid := os.Getpid()
	for _, p := range procs {
		if p.Pid == int32(pid) {
			name, err := p.Name()
			if err != nil {
				t.Fatalf("Name() for current process failed: %v", err)
			}
			if name == "" {
				t.Fatal("current process name is empty")
			}
			if !processRunning(name) {
				t.Fatalf("processRunning(%q) should find current process", name)
			}
			found = true
			break
		}
	}
	if !found {
		t.Skip("current process not found in process list")
	}
}
