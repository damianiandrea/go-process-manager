package message

import "context"

type ExecProcessMsgProducer interface {
	Produce(ctx context.Context, execProcess *ExecProcess) error
}

type ExecProcess struct {
	ProcessName string `json:"process_name"`
}
