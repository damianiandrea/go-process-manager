package process

import "context"

type Manager interface {
	ListRunning() ([]*Process, error)
	Run(ctx context.Context, processName string, args ...string) error
}

type Process struct {
	Pid         int
	ProcessUuid string
}
