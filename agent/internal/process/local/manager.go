package local

import (
	"os"
	"os/exec"
	"sync"

	"github.com/damianiandrea/go-process-manager/agent/internal/process"
)

type ProcessManager struct {
	running map[int]struct{}
	mu      sync.RWMutex
}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{running: make(map[int]struct{})}
}

func (m *ProcessManager) ListRunning() ([]*process.Process, error) {
	processes := make([]*process.Process, 0)
	m.mu.RLock()
	for pid := range m.running {
		processes = append(processes, &process.Process{Pid: pid})
	}
	m.mu.RUnlock()
	return processes, nil
}

func (m *ProcessManager) Run(process string, args ...string) error {
	cmd := exec.Command(process, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	m.mu.Lock()
	m.running[cmd.Process.Pid] = struct{}{}
	m.mu.Unlock()
	go m.cleanupOnProcessExit(cmd.Process)
	return nil
}

func (m *ProcessManager) cleanupOnProcessExit(process *os.Process) {
	_, _ = process.Wait()
	m.mu.Lock()
	delete(m.running, process.Pid)
	m.mu.Unlock()
}
