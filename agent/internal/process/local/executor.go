package local

import (
	"os"
	"os/exec"

	"github.com/damianiandrea/go-process-manager/agent/internal/process"
)

type ProcessExecutor struct {
}

func (e *ProcessExecutor) Exec(process *process.Process) error {
	cmd := exec.Command(process.Name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
