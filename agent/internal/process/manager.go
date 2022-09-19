package process

type Manager interface {
	ListRunning() []*Process
	Run(process string, args ...string) error
}

type Process struct {
	Pid int
}
