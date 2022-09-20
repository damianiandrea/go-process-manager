package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"

	"github.com/damianiandrea/go-process-manager/orchestrator/internal/message"
)

type listRunningProcessesMsgConsumer struct {
	client *Client
}

func NewListRunningProcessesMsgConsumer(client *Client) *listRunningProcessesMsgConsumer {
	return &listRunningProcessesMsgConsumer{client: client}
}

func (c *listRunningProcessesMsgConsumer) Consume(ctx context.Context) error {
	ch := make(chan *nats.Msg, 64)
	defer close(ch)
	sub, err := c.client.conn.ChanSubscribe("agent.*.processes", ch)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("unsubscribing and draining messages: %v", ctx.Err())
			return sub.Drain()
		case msg := <-ch:
			data := msg.Data
			log.Printf("received message: %v", string(data))
			processes := &message.RunningProcesses{}
			if err = json.NewDecoder(bytes.NewReader(data)).Decode(processes); err != nil {
				log.Printf("could not decode message: %v", err)
				continue
			}
		}
	}
}
