package message

import "context"

type RunProcessMsgConsumer interface {
	Consume(ctx context.Context) error
}

type RunProcess struct {
	ProcessName string   `json:"process_name"`
	Args        []string `json:"args"`
}
