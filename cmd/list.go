package cmd

import "flag"

func NewListCmd() *ListCmd {
	c := &ListCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("list", flag.ExitOnError),
		alias:     []string{"ls", "l"},
		shortDesc: "List todos",
		usageLine: "[options]",
	}}
	return c
}

type ListCmd struct {
	Cmd
	json bool
}

func (c *ListCmd) Run() error {
	// decode todos into map

	if c.json {
		// print json
	} else {
		// print raw
	}

	return nil
}
