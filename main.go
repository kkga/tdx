package main

import (
	"log"

	"github.com/kkga/tdx/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
