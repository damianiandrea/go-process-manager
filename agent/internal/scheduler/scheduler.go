package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
	"github.com/damianiandrea/go-process-manager/agent/internal/process"
)

type Heartbeat struct {
	heartRate time.Duration

	processManager process.Manager

	msgProducer message.ListRunningProcessesMsgProducer
}

func NewHeartbeat(heartRate time.Duration, processManager process.Manager,
	msgProducer message.ListRunningProcessesMsgProducer) *Heartbeat {
	return &Heartbeat{heartRate: heartRate, processManager: processManager, msgProducer: msgProducer}
}

func (s *Heartbeat) Beat(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("stopping heartbeat: %v", ctx.Err())
			return ctx.Err()
		case <-time.After(s.heartRate):
			processes, err := s.processManager.ListRunning()
			if err != nil {
				log.Printf("could not get list of running processes: %v", err)
				continue
			}
			running := make([]*message.Process, 0)
			for _, p := range processes {
				running = append(running, &message.Process{Pid: p.Pid, ProcessUuid: p.ProcessUuid})
			}
			now := time.Now().UnixMilli()
			runningProcessesMsg := &message.RunningProcesses{Processes: running, Timestamp: now}
			if err := s.msgProducer.Produce(ctx, runningProcessesMsg); err != nil {
				log.Printf("could not send heartbeat: %v", err)
			}
		}
	}
}
