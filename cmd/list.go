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
	c.fs.StringVar(&c.list, "l", "", "show only todos from specified list")
	c.fs.StringVar(&c.sort, "s", "", "sort todos by field: priority, due, created, status")
	c.fs.StringVar(&c.status, "S", "NEEDS-ACTION", "show only todos with specified status: NEEDS-ACTION, COMPLETED, CANCELLED, ANY")
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

	// check status flag
	if c.status != "" {
		s := vdir.ToDoStatus(c.status)
		switch {
		case s == vdir.StatusNeedsAction || s == vdir.StatusCompleted || s == vdir.StatusCancelled || s == vdir.StatusAny:
			break
		default:
			return fmt.Errorf("Incorrect status filter: %s. See: tdx list -h.", c.status)
		}
	}

	filtered := make(map[vdir.Collection][]vdir.Item)

	for col, items := range collections {
		if c.list != "" {
			if col.Name == c.list {
				items, err = filterByStatus(items, vdir.ToDoStatus(c.status))
				if err != nil {
					return err
				}
				for _, i := range items {
					filtered[*col] = append(filtered[*col], *i)
				}
			}
		} else {
			items, err = filterByStatus(items, vdir.ToDoStatus(c.status))
			if err != nil {
				return err
			}
			for _, i := range items {
				filtered[*col] = append(filtered[*col], *i)
			}
		}
	}

	for col, items := range filtered {
		if c.list == "" {
			sb.WriteString(fmt.Sprintf("== %s (%d) ==\n", col.Name, len(items)))
		}
		for _, i := range items {
			if err = writeItem(&sb, i); err != nil {
				return err
			}
		}
	}

	fmt.Print(sb.String())
	return nil
}

func filterByStatus(items []*vdir.Item, status vdir.ToDoStatus) (filtered []*vdir.Item, err error) {
	if status == vdir.StatusAny {
		return items, nil
	}

	for _, i := range items {
		for _, comp := range i.Ical.Children {
			if comp.Name == ical.CompToDo {
				s, propErr := comp.Props.Text(ical.PropStatus)
				if propErr != nil {
					return nil, propErr
				}
				if s == string(status) {
					filtered = append(filtered, i)
				}
			}
		}
	}
	return
}

func writeItem(sb *strings.Builder, item vdir.Item) error {
	for _, comp := range item.Ical.Children {
		if comp.Name == ical.CompToDo {
			t, err := item.Format()
			if err != nil {
				return err
			}
			sb.WriteString(t)
			sb.WriteString("\n")
		}
	}
	return nil
}
