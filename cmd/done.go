package cmd

import (
	"flag"
	"fmt"
)

func NewDoneCmd() *DoneCmd {
	c := &DoneCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("done", flag.ExitOnError),
		alias:     []string{"do"},
		shortDesc: "Complete todo",
	}}
	return c
}

type DoneCmd struct {
	Cmd
}

func (c *DoneCmd) Run() error {
	fmt.Println(c.fs.Args())

	// check args for uid
	// create a map of all todos
	// find todo by uid
	// set completed prop
	// adjust relevant props accordingly (modified date)
	// encode and write back

	return nil
}
