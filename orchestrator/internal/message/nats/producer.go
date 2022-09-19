package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/damianiandrea/go-process-manager/orchestrator/internal/message"
)

type RunProcessMsgProducer struct {
	client *Client
}

func NewRunProcessMsgProducer(client *Client) *RunProcessMsgProducer {
	return &RunProcessMsgProducer{client: client}
}

func (p *RunProcessMsgProducer) Produce(_ context.Context, run *message.RunProcess) error {
	buffer := bytes.Buffer{}
	if err := json.NewEncoder(&buffer).Encode(run); err != nil {
		return err
	}
	if err := p.client.conn.Publish("process.run", buffer.Bytes()); err != nil {
		log.Printf("could not publish message: %v", err)
		return err
	}
	log.Printf("published message: %v", buffer.String())
	return nil
}
