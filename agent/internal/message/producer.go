package message

import "context"

type ListRunningProcessesMsgProducer interface {
	Produce(ctx context.Context, running *RunningProcesses) error
}

type RunningProcesses struct {
	Running []*Process `json:"running"`
}

type Process struct {
	Pid int `json:"pid"`
}
