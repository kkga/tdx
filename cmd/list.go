package cmd

import (
	"flag"
	"fmt"
)

func NewListCmd() *ListCmd {
	c := &ListCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("list", flag.ExitOnError),
		alias:     []string{"ls", "l"},
		shortDesc: "List todos",
		usageLine: "[options]",
	}}
	c.fs.BoolVar(&c.json, "json", false, "json output")
	return c
}

type ListCmd struct {
	Cmd
	json bool
}

func (c *ListCmd) Run() error {
	// decode todos into map

	if c.json {
		fmt.Println("list -json called")
	} else {
		fmt.Println("list called")
	}

	return nil
}
