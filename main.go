package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kkga/ctodo/cmd"
)

var version = "dev"

func main() {
	// log.SetFlags(0)

	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Printf("kks %s\n", version)
		os.Exit(0)
	}

	if err := cmd.Root(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
