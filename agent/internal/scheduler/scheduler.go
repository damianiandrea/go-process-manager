package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
	"github.com/damianiandrea/go-process-manager/agent/internal/process"
)

type HeartbeatScheduler struct {
	agentId        string
	heartRate      time.Duration
	processManager process.Manager
	msgProducer    message.ListRunningProcessesMsgProducer
}

func NewHeartbeatScheduler(agentId string, heartRate time.Duration, processManager process.Manager,
	msgProducer message.ListRunningProcessesMsgProducer) *HeartbeatScheduler {
	return &HeartbeatScheduler{agentId: agentId, heartRate: heartRate, processManager: processManager,
		msgProducer: msgProducer}
}

func (s *HeartbeatScheduler) Beat(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("stopping heartbeat: %v", ctx.Err())
			return ctx.Err()
		case <-time.After(s.heartRate):
			processes := s.processManager.ListRunning()
			running := make([]*message.Process, 0)
			for _, p := range processes {
				running = append(running, &message.Process{Pid: p.Pid})
			}
			runningProcessesMsg := &message.RunningProcesses{AgentId: s.agentId, Processes: running}
			if err := s.msgProducer.Produce(ctx, runningProcessesMsg); err != nil {
				log.Printf("could not send heartbeat: %v", err)
			}
		}
	}
}
