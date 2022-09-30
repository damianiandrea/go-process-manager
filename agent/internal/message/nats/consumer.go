package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
	"github.com/damianiandrea/go-process-manager/agent/internal/process"
)

type runProcessMsgConsumer struct {
	client         *Client
	processManager process.Manager
}

func NewRunProcessMsgConsumer(client *Client, processManager process.Manager) *runProcessMsgConsumer {
	return &runProcessMsgConsumer{client: client, processManager: processManager}
}

func (c *runProcessMsgConsumer) Consume(ctx context.Context) error {
	msgCh := make(chan *nats.Msg, 64)
	defer close(msgCh)
	sub, err := c.client.conn.ChanQueueSubscribe("process.run", "agent", msgCh)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("unsubscribing and draining messages: %v", ctx.Err())
			return sub.Drain()
		case msg := <-msgCh:
			data := msg.Data
			log.Printf("received message: %v", string(data))
			run := &message.RunProcess{}
			if err = json.NewDecoder(bytes.NewReader(data)).Decode(run); err != nil {
				log.Printf("could not decode message: %v", err)
				continue
			}
			if err = c.processManager.Run(ctx, run.ProcessName, run.Args...); err != nil {
				log.Printf("could not run process: %v", err)
			}
		}
	}
}
