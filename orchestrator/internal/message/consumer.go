package message

import "context"

type ListRunningProcessesMsgConsumer interface {
	Consume(ctx context.Context) error
}

type RunningProcesses struct {
	AgentId   string     `json:"agent_id"`
	Processes []*Process `json:"processes"`
}

type Process struct {
	Pid int `json:"pid"`
}
