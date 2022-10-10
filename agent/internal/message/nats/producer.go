package nats

import (
	"context"
	"fmt"
	"log"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
)

type listRunningProcessesMsgProducer struct {
	agentId string

	client  *Client
	encoder message.Encoder
}

func NewListRunningProcessesMsgProducer(agentId string, client *Client, encoder message.Encoder,
) *listRunningProcessesMsgProducer {
	return &listRunningProcessesMsgProducer{agentId: agentId, client: client, encoder: encoder}
}

func (p *listRunningProcessesMsgProducer) Produce(_ context.Context, processes *message.RunningProcesses) error {
	processes.AgentId = p.agentId
	bytes, err := p.encoder.Encode(processes)
	if err != nil {
		return err
	}
	subject := fmt.Sprintf("agent.%s.processes", p.agentId)
	if err := p.client.conn.Publish(subject, bytes); err != nil {
		log.Printf("could not publish message: %v", err)
		return err
	}
	return nil
}

type processOutputMsgProducer struct {
	agentId string

	client *Client
}

func NewProcessOutputMsgProducer(agentId string, client *Client) *processOutputMsgProducer {
	return &processOutputMsgProducer{agentId: agentId, client: client}
}

func (p *processOutputMsgProducer) Produce(_ context.Context, output *message.ProcessOutput) error {
	subject := fmt.Sprintf("agent.%s.processes.%s.output", p.agentId, output.ProcessUuid)
	if _, err := p.client.js.Publish(subject, output.Data); err != nil {
		log.Printf("could not publish message: %v", err)
		return err
	}
	return nil
}
