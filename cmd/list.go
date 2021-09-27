package cmd

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/emersion/go-ical"

	"github.com/fatih/color"
	"github.com/hako/durafmt"
	"github.com/kkga/tdx/vdir"
)

func NewListCmd() *ListCmd {
	c := &ListCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("list", flag.ExitOnError),
		alias:     []string{"ls", "l"},
		shortDesc: "List todos, optionally filtered by query",
		usageLine: "[options] [query]",
	}}
	c.fs.BoolVar(&c.json, "json", false, "json output")
	c.fs.BoolVar(&c.description, "desc", false, "show todo description in output")
	c.fs.BoolVar(&c.multiline, "ml", false, "use 2-line output for dates and description")
	c.fs.StringVar(&c.listFlag, "l", "", "show only todos from specified `list`")
	c.fs.BoolVar(&c.allLists, "a", false, "show todos from all lists (overrides -l)")
	c.fs.StringVar(&c.sortOption, "s", "", "sort todos by `field`: PRIO, DUE, STATUS")
	c.fs.StringVar(&c.statusFilter, "S", "", "show only todos with specified `status`: NEEDS-ACTION, COMPLETED, CANCELLED, ANY")
	return c
}

type ListCmd struct {
	Cmd
	json         bool
	multiline    bool
	description  bool
	allLists     bool
	sortOption   string
	statusFilter string
}

type sortOption string

const (
	sortOptionStatus sortOption = "STATUS"
	sortOptionPrio   sortOption = "PRIO"
	sortOptionDue    sortOption = "DUE"
)

func (c *ListCmd) Run() error {
	var query string

	if len(c.fs.Args()) > 0 {
		query = strings.Join(c.fs.Args(), "")
	}

	// process status filter
	if c.statusFilter == "" {
		c.statusFilter = c.conf.DefaultStatus
	}
	switch vdir.ToDoStatus(c.statusFilter) {
	case vdir.StatusNeedsAction, vdir.StatusCompleted, vdir.StatusCancelled, vdir.StatusAny:
		break
	default:
		return fmt.Errorf("Incorrect status filter: %s. See: tdx list -h", c.statusFilter)
	}

	// process sort option
	if c.sortOption == "" {
		c.sortOption = c.conf.DefaultSort
	}
	switch sortOption(c.sortOption) {
	case sortOptionStatus, sortOptionPrio, sortOptionDue:
		break
	default:
		return fmt.Errorf("Incorrect sort option: %s. See: tdx list -h", c.sortOption)
	}

	// if cmd has collection specified via flag, delete other collections from map
	collections := c.vdir
	if c.collection != nil && c.allLists == false {
		for col := range collections {
			if col != c.collection {
				delete(collections, col)
			}
		}
	}

	// filter and sort items
	var m = make(map[vdir.Collection][]vdir.Item)
	for k, v := range collections {
		items, err := filterByStatus(v, vdir.ToDoStatus(c.statusFilter))
		if err != nil {
			return err
		}
		items, err = filterByQuery(items, query)
		if err != nil {
			return err
		}

		switch sortOption(c.sortOption) {
		case sortOptionPrio:
			sort.Sort(vdir.ByPriority(items))
		case sortOptionDue:
			sort.Sort(vdir.ByDue(items))
		// TODO: implement due and status sorting
		case sortOptionStatus:
			// sort.Sort(vdir.ByStatus(items))
		}

		for _, item := range items {
			m[*k] = append(m[*k], *item)
		}
	}

	// prepare output
	var sb = strings.Builder{}
	for col, items := range m {
		if len(m) > 1 {
			colorList := color.New(color.Bold, color.FgYellow).SprintFunc()
			sb.WriteString(colorList(fmt.Sprintf("== %s (%d) ==\n", col.Name, len(items))))
		}

		for _, i := range items {
			if err := writeItem(&sb, *c, i); err != nil {
				return err
			}
		}
	}

	fmt.Print(sb.String())
	return nil
}

func parseDueDate(dur time.Duration) (duration string, err error) {
	duration = durafmt.Parse(dur).String()
	return
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

func filterByQuery(items []*vdir.Item, query string) (filtered []*vdir.Item, err error) {
	if query == "" {
		return items, nil
	}

	for _, i := range items {
		for _, comp := range i.Ical.Children {
			if comp.Name == ical.CompToDo {
				summary, propErr := comp.Props.Text(ical.PropSummary)
				if propErr != nil {
					return nil, propErr
				}
				if strings.Contains(strings.ToLower(summary), strings.ToLower(query)) {
					filtered = append(filtered, i)
				}
			}
		}
	}
	return
}

func writeItem(sb *strings.Builder, c ListCmd, item vdir.Item) error {
	opts := []vdir.FormatOption{}
	if c.multiline {
		opts = append(opts, vdir.FormatMultiline)
	}
	if c.description {
		opts = append(opts, vdir.FormatDescription)
	}

	s, err := item.Format(opts...)
	if err != nil {
		return err
	}

	sb.WriteString(s)

	return nil
}
