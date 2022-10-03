package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/damianiandrea/go-process-manager/orchestrator/internal/message"
)

type runProcessMsgProducer struct {
	client *Client
}

func NewRunProcessMsgProducer(client *Client) *runProcessMsgProducer {
	return &runProcessMsgProducer{client: client}
}

func (p *runProcessMsgProducer) Produce(_ context.Context, run *message.RunProcess) error {
	buffer := bytes.Buffer{}
	if err := json.NewEncoder(&buffer).Encode(run); err != nil {
		return err
	}
	if err := p.client.conn.Publish("process.run", buffer.Bytes()); err != nil {
		log.Printf("could not publish message: %v", err)
		return err
	}
	return nil
}
