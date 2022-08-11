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
		if err != nil {
			panic(err)
		}
		tmtf, err := termiter.ReadTermiterFile(file)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v\n", tmtf)
		file.Close()
	}
}
