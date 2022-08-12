package termiter

import (
	"log"
	"os/exec"
)

type OutProxy struct {
	logger *log.Logger
	outErr bool
}

func (o *OutProxy) Write(p []byte) (n int, err error) {
	prev := "[INFO]"
	if o.outErr {
		prev = "[ERR]"
	}
	o.logger.Printf("%s %s", prev, string(p))
	return len(p), nil
}

func RunCommand(name string, args []string) int {
	cmd := exec.Command(name, args...)
	log := log.Default()
	cmd.Stdout = &OutProxy{logger: log, outErr: false}
	cmd.Stderr = &OutProxy{logger: log, outErr: true}
	cmd.Stdin = nil
	err := cmd.Run()
	return CheckExecutionError(err)
}
