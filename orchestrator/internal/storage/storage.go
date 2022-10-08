package storage

import "context"

type ProcessStore interface {
	GetAll(ctx context.Context) ([]Process, error)
	Put(ctx context.Context, agentId string, processes []Process) error
}

type Process struct {
	Pid         int
	ProcessUuid string
	AgentId     string
	LastSeen    int64
}
