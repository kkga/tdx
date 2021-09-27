package main

import (
	"log"
	"os"

	"github.com/kkga/tdx/cmd"
)

var version = "dev"

func main() {
	log.SetFlags(0)

	if err := cmd.Root(os.Args[1:], version); err != nil {
		log.Fatal(err)
	}
}
