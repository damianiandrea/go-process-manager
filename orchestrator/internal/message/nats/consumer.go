package nats

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"

	"github.com/damianiandrea/go-process-manager/orchestrator/internal/message"
	"github.com/damianiandrea/go-process-manager/orchestrator/internal/storage"
)

type listRunningProcessesMsgConsumer struct {
	client  *Client
	decoder message.Decoder

	processStore storage.ProcessStore
}

func NewListRunningProcessesMsgConsumer(client *Client, decoder message.Decoder, processStore storage.ProcessStore,
) *listRunningProcessesMsgConsumer {
	return &listRunningProcessesMsgConsumer{client: client, decoder: decoder, processStore: processStore}
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
			processesMsg := &message.RunningProcesses{}
			if err = c.decoder.Decode(msg.Data, processesMsg); err != nil {
				log.Printf("could not decode message: %v", err)
				continue
			}
			processes := make([]storage.Process, 0)
			for _, p := range processesMsg.Processes {
				processes = append(processes, storage.Process{Pid: p.Pid, ProcessUuid: p.ProcessUuid,
					AgentId: processesMsg.AgentId, LastSeen: processesMsg.Timestamp})
			}
			if err = c.processStore.Put(ctx, processesMsg.AgentId, processes); err != nil {
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

	go func(sub *nats.Subscription) {
		for {
			select {
			case <-ctx.Done():
				if err := sub.Unsubscribe(); err != nil {
					log.Printf("could not unsubscribe: %v", err)
				}
				close(msgCh)
				return
			case msg := <-msgCh:
				outputCh <- msg.Data
			}
		}
	}(sub)

	return nil
}
