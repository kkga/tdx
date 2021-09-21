package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/emersion/go-ical"
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

	sb := strings.Builder{}

	collections, err := c.root.Collections()
	if err != nil {
		return err
	}

	for _, c := range collections {
		sb.WriteString(fmt.Sprintf("\n== %s ==\n", c.Name))
		items, err := c.Items()
		if err != nil {
			return err
		}
		for i, item := range items {
			for _, comp := range item.Children {
				if comp.Name == ical.CompToDo {
					summary := comp.Props.Get(ical.PropSummary)
					sb.WriteString(fmt.Sprintf("%d: %s\n", i, summary.Value))
				}
			}
		}
	}

	fmt.Print(sb.String())

	return nil
}
