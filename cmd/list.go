package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
)

func NewListCmd() *ListCmd {
	c := &ListCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("list", flag.ExitOnError),
		alias:     []string{"ls", "l"},
		shortDesc: "List todos",
		usageLine: "[options]",
	}}
	c.fs.BoolVar(&c.json, "json", false, "json output")
	c.fs.StringVar(&c.listFlag, "l", "", "show only todos from specified list")
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

	if c.listFlag != "" {
		for col, items := range collections {
			if col.Name == c.listFlag {
				if len(items) == 0 {
					continue
				}
				sb.WriteString(fmt.Sprintf("== %s (%d) ==\n", col.Name, len(items)))
				for _, i := range collections[col] {
					if err = writeItem(&sb, i.Id, *i); err != nil {
						return err
					}
				}
				break
			}
		}
	} else {
		for col, items := range collections {
			if len(items) == 0 {
				continue
			}
			sb.WriteString(fmt.Sprintf("== %s (%d) ==\n", col.Name, len(items)))
			for _, i := range collections[col] {
				if err = writeItem(&sb, i.Id, *i); err != nil {
					return err
				}
			}
		}
	}

	fmt.Print(sb.String())
	return nil
}

func writeItem(sb *strings.Builder, id int, item vdir.Item) error {
	for _, comp := range item.Ical.Children {
		if comp.Name == ical.CompToDo {
			t, err := item.Format()
			if err != nil {
				return err
			}
			sb.WriteString(fmt.Sprintf("%2d ", id))
			sb.WriteString(t)
			sb.WriteString("\n")
		}
	}
	return nil
}
