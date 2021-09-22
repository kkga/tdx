package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
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

	// if c.listFlag != "" {
	// 	// TODO this is now handled in Init
	// 	for _, col := range collections {
	// 		if col.Name == c.listFlag {
	// 			items, err := col.Items()
	// 			if err != nil {
	// 				return err
	// 			}
	// 			if len(items) == 0 {
	// 				continue
	// 			}
	// 			if err = writeItems(&sb, items); err != nil {
	// 				return err
	// 			}
	// 			break
	// 		}
	// 	}
	// } else {
	for c, items := range collections {
		if len(items) == 0 {
			continue
		}
		fmt.Println(items)
		sb.WriteString(fmt.Sprintf("== %s (%d) ==\n", c.Name, len(items)))
		for _, i := range collections[c] {
			if err = writeItem(&sb, i.Id, *i); err != nil {
				return err
			}
		}
	}
	// }

	fmt.Print(sb.String())
	return nil
}

func writeItem(sb *strings.Builder, id int, item vdir.Item) error {
	for _, comp := range item.Ical.Children {
		if comp.Name == ical.CompToDo {
			t, err := vtodo.Format(comp)
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
