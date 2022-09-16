package process

type Executor interface {
	Exec(process *Process) error
}

type Process struct {
	Name string
}
