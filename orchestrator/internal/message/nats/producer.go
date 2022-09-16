package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/damianiandrea/go-process-manager/orchestrator/internal/message"
)

type ExecProcessMsgProducer struct {
	client *Client
}

func NewExecProcessMsgProducer(client *Client) *ExecProcessMsgProducer {
	return &ExecProcessMsgProducer{client: client}
}

func (p *ExecProcessMsgProducer) Produce(_ context.Context, execProcess *message.ExecProcess) error {
	buffer := bytes.Buffer{}
	if err := json.NewEncoder(&buffer).Encode(execProcess); err != nil {
		return err
	}
	if err := p.client.conn.Publish("process.exec", buffer.Bytes()); err != nil {
		log.Printf("could not publish message: %v", err)
		return err
	}
	log.Printf("published message: %v", buffer.String())
	return nil
}
