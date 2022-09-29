package storage

type ProcessStore interface {
	GetAll() ([]Process, error)
	Put(agentId string, processes []Process) error
}

type Process struct {
	Pid         int
	ProcessUuid string
	AgentId     string
	LastSeen    int64
}
