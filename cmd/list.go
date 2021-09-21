package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vtodo"
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
		for _, item := range items {
			for _, comp := range item.Children {
				if comp.Name == ical.CompToDo {
					t, err := vtodo.Format(comp)
					if err != nil {
						return err
					}
					sb.WriteString(t)
					sb.WriteString("\n")
				}
			}
		}
	}

	fmt.Print(sb.String())

	return nil
}
