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
	c.fs.StringVar(&c.list, "l", "", "show only todos from specified list")
	c.fs.StringVar(&c.sort, "s", "", "sort todos by field: priority, due, created, status")
	c.fs.StringVar(&c.status, "S", "", "show only todos with specified status: NEEDS-ACTION, COMPLETED, CANCELLED, ANY")
	return c
}

type ListCmd struct {
	Cmd
	json   bool
	sort   string
	status string
}

func (c *ListCmd) Run() error {
	sb := strings.Builder{}

	collections, err := c.root.Collections()
	if err != nil {
		return err
	}

	if c.list != "" {
		for i, col := range collections {
			if col.Name == c.list {
				items, err := col.Items()
				if err != nil {
					return err
				}
				if len(items) == 0 {
					continue
				}
				if err = writeItems(&sb, items); err != nil {
					return err
				}
				break
			} else if i == len(collections)-1 {
				return fmt.Errorf("List not found: %s", c.list)
			}
		}

	} else {
		for _, col := range collections {
			items, err := col.Items()
			if err != nil {
				return err
			}

			if len(items) == 0 {
				continue
			}

			sb.WriteString(fmt.Sprintf("\n== %s ==\n", col.Name))
			if err = writeItems(&sb, items); err != nil {
				return err
			}
		}
	}

	fmt.Print(sb.String())
	return nil
}

func writeItems(sb *strings.Builder, items []ical.Calendar) error {
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
	return nil
}
