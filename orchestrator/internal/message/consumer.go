package message

import "context"

type ListRunningProcessesMsgConsumer interface {
	Consume(ctx context.Context) error
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

type ProcessOutputMsgConsumer interface {
	ChanConsume(ctx context.Context, processUuid string, outputCh chan []byte) error
}
