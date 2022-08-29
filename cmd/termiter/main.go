package main

import (
	"errors"
	"os"

	"github.com/albertoaer/termiter"
)

func main() {
	if len(os.Args) > 1 {
		path := os.Args[1]
		file, err := os.Open(path)
		termiter.PanicIfError(err)
		tmtf, err := termiter.ReadTermiterFile(file)
		termiter.PanicIfError(err)
		termiter.PanicIfError(tmtf.Verify())
		termiter.PanicIfError(file.Close())
		termiter.PanicIfError(err)
		context, err := termiter.NewExecutionContext(tmtf, os.Args[2:])
		termiter.PanicIfError(err)
		action, err := context.GetStartAction()
		termiter.PanicIfError(err)
		action.Run(context)
	} else {
		termiter.PanicIfError(errors.New("Please provide a Termiter file"))
	}
}
