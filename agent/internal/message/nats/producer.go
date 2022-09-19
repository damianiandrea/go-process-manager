package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
)

type ListRunningProcessesMsgProducer struct {
	client *Client
}

func NewListRunningProcessesMsgProducer(client *Client) *ListRunningProcessesMsgProducer {
	return &ListRunningProcessesMsgProducer{client: client}
}

func (p *ListRunningProcessesMsgProducer) Produce(_ context.Context, running *message.RunningProcesses) error {
	buffer := bytes.Buffer{}
	if err := json.NewEncoder(&buffer).Encode(running); err != nil {
		return err
	}
	if err := p.client.conn.Publish("agent.x.process", buffer.Bytes()); err != nil { // TODO: x = agent_id
		log.Printf("could not publish message: %v", err)
		return err
	}
	log.Printf("published message: %v", buffer.String())
	return nil
}
