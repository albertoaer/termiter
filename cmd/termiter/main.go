package main

import (
	"fmt"
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
		fmt.Printf("%v\n", tmtf)
		termiter.PanicIfError(file.Close())
	}
}
