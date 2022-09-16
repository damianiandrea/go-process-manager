package message

import "context"

type ExecProcessMsgConsumer interface {
	Consume(ctx context.Context) error
}

type ExecProcess struct {
	ProcessName string `json:"process_name"`
}
