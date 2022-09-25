package local

import (
	"context"
	"os"
	"os/exec"
	"sync"

	"github.com/google/uuid"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
	"github.com/damianiandrea/go-process-manager/agent/internal/process"
)

type processManager struct {
	msgProducer message.ProcessOutputMsgProducer

	running map[int]*process.Process
	mu      sync.RWMutex
}

func NewProcessManager(msgProducer message.ProcessOutputMsgProducer) *processManager {
	return &processManager{msgProducer: msgProducer, running: make(map[int]*process.Process)}
}

func (m *processManager) ListRunning() ([]*process.Process, error) {
	processes := make([]*process.Process, 0)
	m.mu.RLock()
	for _, p := range m.running {
		processes = append(processes, &process.Process{Pid: p.Pid, ProcessUuid: p.ProcessUuid})
	}
	m.mu.RUnlock()
	return processes, nil
}

func (m *processManager) Run(ctx context.Context, processName string, args ...string) error {
	processUuid := uuid.NewString()
	writer := process.NewOutputWriter(ctx, processUuid, m.msgProducer)
	cmd := exec.Command(processName, args...)
	cmd.Stdout = writer
	cmd.Stderr = writer
	if err := cmd.Start(); err != nil {
		return err
	}
	m.mu.Lock()
	m.running[cmd.Process.Pid] = &process.Process{Pid: cmd.Process.Pid, ProcessUuid: processUuid}
	m.mu.Unlock()
	go m.cleanupOnProcessExit(cmd.Process)
	return nil
}

func (m *processManager) cleanupOnProcessExit(process *os.Process) {
	_, _ = process.Wait()
	m.mu.Lock()
	delete(m.running, process.Pid)
	m.mu.Unlock()
}
