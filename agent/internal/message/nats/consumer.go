package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
	"github.com/damianiandrea/go-process-manager/agent/internal/process"
	"github.com/nats-io/nats.go"
)

type ExecProcessMsgConsumer struct {
	client   *Client
	executor process.Executor
}

func NewExecProcessMsgConsumer(client *Client, executor process.Executor) *ExecProcessMsgConsumer {
	return &ExecProcessMsgConsumer{client: client, executor: executor}
}

func (c *ExecProcessMsgConsumer) Consume(ctx context.Context) error {
	ch := make(chan *nats.Msg, 64)
	defer close(ch)
	sub, err := c.client.conn.ChanQueueSubscribe("process.exec", "agent", ch)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return sub.Drain()
		case execProcessMsg := <-ch:
			data := execProcessMsg.Data
			log.Printf("received message: %v", string(data))
			execProcess := &message.ExecProcess{}
			if err = json.NewDecoder(bytes.NewReader(data)).Decode(execProcess); err != nil {
				log.Printf("could not decode message: %v", err)
				continue
			}
			if err = c.executor.Exec(&process.Process{Name: execProcess.ProcessName}); err != nil {
				log.Printf("could not execute process: %v", err)
			}
		}
	}
}
