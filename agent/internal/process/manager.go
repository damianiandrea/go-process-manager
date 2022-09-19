package process

type Manager interface {
	ListRunning() ([]*Process, error)
	Run(process string, args ...string) error
}

type Process struct {
	Pid int
}
