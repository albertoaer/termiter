package termiter

import (
	"fmt"
	"os"
	"os/exec"
)

func PanicIfError(errs ...error) {
	for _, err := range errs {
		if err != nil {
			fmt.Fprintf(os.Stderr, "[EXIT] due to: %s", err.Error())
			os.Exit(-1)
		}
	}
}

func CheckExecutionError(err error) int {
	if err != nil {
		if exiterr, is := err.(*exec.ExitError); is {
			return exiterr.ExitCode()
		} else {
			PanicIfError(err)
			return -1
		}
	}
	return 0
}
