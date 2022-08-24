package termiter

import (
	"log"
	"os/exec"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

type OutProxy struct {
	logger *log.Logger
	outErr bool
}

func (o *OutProxy) Write(p []byte) (n int, err error) {
	prev := "[INFO]"
	paint := color.New(color.BgBlue, color.FgWhite)
	if o.outErr {
		prev = "[ERR]"
		paint = color.New(color.BgRed, color.FgWhite)
	}
	output := paint.Sprintf("%s %s", prev, string(p))
	o.logger.Printf(output)
	return len(p), nil
}

func RunCommand(name string, args []string, env []string) int {
	cmd := exec.Command(name, args...)
	cmd.Env = env
	log := log.Default()
	log.SetOutput(colorable.NewColorableStdout())
	cmd.Stdout = &OutProxy{logger: log, outErr: false}
	cmd.Stderr = &OutProxy{logger: log, outErr: true}
	cmd.Stdin = nil
	err := cmd.Run()
	return CheckExecutionError(err)
}
