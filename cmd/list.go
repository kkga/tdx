package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/emersion/go-ical"
	"github.com/fatih/color"
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
	c.fs.BoolVar(&c.allLists, "a", false, "show todos from all lists (overrides -l)")
	c.fs.StringVar(&c.sort, "s", "", "sort todos by field: priority, due, created, status")
	c.fs.StringVar(&c.status, "S", "NEEDS-ACTION", "show only todos with specified status: NEEDS-ACTION, COMPLETED, CANCELLED, ANY")
	return c
}

type ListCmd struct {
	Cmd
	json     bool
	allLists bool
	sort     string
	status   string
}

func (c *ListCmd) Run() error {
	// query := c.fs.Args()

	collections := c.allCollections

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

	// if cmd has collection specified via flag, delete other collections from map
	if c.collection != nil && c.allLists == false {
		for col := range collections {
			if col != c.collection {
				delete(collections, col)
			}
		}
	}

	var filtered = make(map[vdir.Collection][]vdir.Item)
	for col, items := range collections {
		items, err := filterByStatus(items, vdir.ToDoStatus(c.status))
		if err != nil {
			return err
		}
		for _, i := range items {
			filtered[*col] = append(filtered[*col], *i)
		}
	}

	sb := strings.Builder{}
	for col, items := range filtered {
		// if len(filtered) > 1 {
		colorList := color.New(color.Bold, color.FgYellow).SprintFunc()
		sb.WriteString(colorList(fmt.Sprintf("== %s (%d) ==\n", col.Name, len(items))))
		// }
		for _, i := range items {
			if err := writeItem(&sb, i); err != nil {
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
