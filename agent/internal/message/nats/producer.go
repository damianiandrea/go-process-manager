package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
)

type listRunningProcessesMsgProducer struct {
	agentId string

	client *Client
}

func NewListRunningProcessesMsgProducer(agentId string, client *Client) *listRunningProcessesMsgProducer {
	return &listRunningProcessesMsgProducer{agentId: agentId, client: client}
}

func (p *listRunningProcessesMsgProducer) Produce(_ context.Context, processes *message.RunningProcesses) error {
	processes.AgentId = p.agentId
	buffer := bytes.Buffer{}
	if err := json.NewEncoder(&buffer).Encode(processes); err != nil {
		return err
	}
	subject := fmt.Sprintf("agent.%s.processes", p.agentId)
	if err := p.client.conn.Publish(subject, buffer.Bytes()); err != nil {
		log.Printf("could not publish message: %v", err)
		return err
	}
	log.Printf("published message: %v", buffer.String())
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
