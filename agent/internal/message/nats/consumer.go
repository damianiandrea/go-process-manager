package nats

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
	"github.com/damianiandrea/go-process-manager/agent/internal/process"
)

type runProcessMsgConsumer struct {
	client  *Client
	decoder message.Decoder

	processManager process.Manager
}

func NewRunProcessMsgConsumer(client *Client, decoder message.Decoder, processManager process.Manager) *runProcessMsgConsumer {
	return &runProcessMsgConsumer{client: client, decoder: decoder, processManager: processManager}
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
			run := &message.RunProcess{}
			if err = c.decoder.Decode(msg.Data, run); err != nil {
				log.Printf("could not decode message: %v", err)
				continue
			}
			if err = c.processManager.Run(ctx, run.ProcessName, run.Args...); err != nil {
				log.Printf("could not run process: %v", err)
			}
		}
	}
}
