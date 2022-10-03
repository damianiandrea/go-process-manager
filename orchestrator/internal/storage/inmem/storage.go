package inmem

import (
	"sync"

	"github.com/damianiandrea/go-process-manager/orchestrator/internal/storage"
)

type ProcessStore struct {
	processesByAgent map[string][]storage.Process
	mu               sync.RWMutex
}

func NewProcessStore() *ProcessStore {
	return &ProcessStore{processesByAgent: make(map[string][]storage.Process)}
}

func (s *ProcessStore) GetAll() ([]storage.Process, error) {
	s.mu.RLock()
	processes := make([]storage.Process, 0)
	for _, agentProcesses := range s.processesByAgent {
		processes = append(processes, agentProcesses...)
	}
	s.mu.RUnlock()
	return processes, nil
}

func (s *ProcessStore) Put(agentId string, processes []storage.Process) error {
	s.mu.Lock()
	s.processesByAgent[agentId] = processes
	s.mu.Unlock()
	return nil
}
