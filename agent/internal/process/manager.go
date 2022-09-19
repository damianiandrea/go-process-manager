package process

type Manager interface {
	ListRunning() ([]*Process, error)
	Run(processName string, args ...string) error
}

type Process struct {
	Pid         int
	ProcessUuid string
}
