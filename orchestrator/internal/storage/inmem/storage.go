package inmem

import "github.com/damianiandrea/go-process-manager/orchestrator/internal/storage"

type ProcessStore struct {
	processesByAgent map[string][]storage.Process
}

func NewProcessStore() *ProcessStore {
	return &ProcessStore{processesByAgent: make(map[string][]storage.Process)}
}

func (s *ProcessStore) GetAll() ([]storage.Process, error) {
	processes := make([]storage.Process, 0)
	for _, agentProcesses := range s.processesByAgent {
		processes = append(processes, agentProcesses...)
	}
	return processes, nil
}

func (s *ProcessStore) Put(agentId string, processes []storage.Process) error {
	s.processesByAgent[agentId] = processes
	return nil
}
