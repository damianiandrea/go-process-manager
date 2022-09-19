package message

import "context"

type RunProcessMsgProducer interface {
	Produce(ctx context.Context, run *RunProcess) error
}

type RunProcess struct {
	ProcessName string   `json:"process_name"`
	Args        []string `json:"args,omitempty"`
}
