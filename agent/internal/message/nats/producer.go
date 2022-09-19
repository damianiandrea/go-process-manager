package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
)

type ListRunningProcessesMsgProducer struct {
	client *Client
}

func NewListRunningProcessesMsgProducer(client *Client) *ListRunningProcessesMsgProducer {
	return &ListRunningProcessesMsgProducer{client: client}
}

func (p *ListRunningProcessesMsgProducer) Produce(_ context.Context, processes *message.RunningProcesses) error {
	buffer := bytes.Buffer{}
	if err := json.NewEncoder(&buffer).Encode(processes); err != nil {
		return err
	}
	if err := p.client.conn.Publish(fmt.Sprintf("agent.%s.processes", processes.AgentId), buffer.Bytes()); err != nil {
		log.Printf("could not publish message: %v", err)
		return err
	}
	log.Printf("published message: %v", buffer.String())
	return nil
}