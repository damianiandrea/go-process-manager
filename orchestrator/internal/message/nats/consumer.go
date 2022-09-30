package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"

	"github.com/damianiandrea/go-process-manager/orchestrator/internal/message"
	"github.com/damianiandrea/go-process-manager/orchestrator/internal/storage"
)

type listRunningProcessesMsgConsumer struct {
	client       *Client
	processStore storage.ProcessStore
}

func NewListRunningProcessesMsgConsumer(client *Client, processStore storage.ProcessStore) *listRunningProcessesMsgConsumer {
	return &listRunningProcessesMsgConsumer{client: client, processStore: processStore}
}

func (c *listRunningProcessesMsgConsumer) Consume(ctx context.Context) error {
	msgCh := make(chan *nats.Msg, 64)
	defer close(msgCh)
	sub, err := c.client.conn.ChanSubscribe("agent.*.processes", msgCh)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("unsubscribing and draining messages: %v", ctx.Err())
			return sub.Drain()
		case msg := <-msgCh:
			data := msg.Data
			log.Printf("received message: %v", string(data))
			processesMsg := &message.RunningProcesses{}
			if err = json.NewDecoder(bytes.NewReader(data)).Decode(processesMsg); err != nil {
				log.Printf("could not decode message: %v", err)
				continue
			}
			processes := make([]storage.Process, 0)
			for _, p := range processesMsg.Processes {
				processes = append(processes, storage.Process{Pid: p.Pid, ProcessUuid: p.ProcessUuid,
					AgentId: processesMsg.AgentId, LastSeen: processesMsg.Timestamp})
			}
			if err = c.processStore.Put(processesMsg.AgentId, processes); err != nil {
				log.Printf("could not store running processes: %v", err)
			}
		}
	}
}

type processOutputMsgConsumer struct {
	client *Client
}

func NewProcessOutputMsgConsumer(client *Client) *processOutputMsgConsumer {
	return &processOutputMsgConsumer{client: client}
}

func (c *processOutputMsgConsumer) ChanConsume(ctx context.Context, processUuid string, outputCh chan []byte) error {
	msgCh := make(chan *nats.Msg, 64)
	subject := fmt.Sprintf("agent.*.processes.%s.output", processUuid)
	sub, err := c.client.js.ChanSubscribe(subject, msgCh, nats.OrderedConsumer())
	if err != nil {
		return err
	}

	go func(*nats.Subscription) {
		for {
			select {
			case <-ctx.Done():
				log.Printf("unsubscribing: %v", ctx.Err())
				if err := sub.Unsubscribe(); err != nil {
					log.Printf("could not unsubscribe: %v", err)
				}
				close(msgCh)
				return
			case msg := <-msgCh:
				data := msg.Data
				log.Printf("received message: %v", string(data))
				outputCh <- data
			}
		}
	}(sub)

	return nil
}
