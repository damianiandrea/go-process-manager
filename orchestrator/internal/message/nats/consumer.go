package nats

import (
	"bytes"
	"context"
	"encoding/json"
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
	ch := make(chan *nats.Msg, 64)
	defer close(ch)
	sub, err := c.client.conn.ChanSubscribe("agent.*.processes", ch)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("unsubscribing and draining messages: %v", ctx.Err())
			return sub.Drain()
		case msg := <-ch:
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
