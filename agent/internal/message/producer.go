package message

import "context"

type ListRunningProcessesMsgProducer interface {
	Produce(ctx context.Context, processes *RunningProcesses) error
}

type RunningProcesses struct {
	AgentId   string     `json:"agent_id"`
	Processes []*Process `json:"processes"`
	Timestamp int64      `json:"timestamp"`
}

type Process struct {
	Pid         int    `json:"pid"`
	ProcessUuid string `json:"process_uuid"`
}
