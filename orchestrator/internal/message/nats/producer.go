package nats

import (
	"context"
	"log"

	"github.com/damianiandrea/go-process-manager/orchestrator/internal/message"
)

type runProcessMsgProducer struct {
	client  *Client
	encoder message.Encoder
}

func NewRunProcessMsgProducer(client *Client, encoder message.Encoder) *runProcessMsgProducer {
	return &runProcessMsgProducer{client: client, encoder: encoder}
}

func (p *runProcessMsgProducer) Produce(_ context.Context, run *message.RunProcess) error {
	bytes, err := p.encoder.Encode(run)
	if err != nil {
		return err
	}
	if err := p.client.conn.Publish("process.run", bytes); err != nil {
		log.Printf("could not publish message: %v", err)
		return err
	}
	return nil
}
