package cmd

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/emersion/go-ical"

	"github.com/fatih/color"
	"github.com/kkga/tdx/vdir"
)

func NewListCmd() *ListCmd {
	c := &ListCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("list", flag.ExitOnError),
		alias:     []string{"ls", "l"},
		short:     "List todos, optionally filtered by query",
		usageLine: "[options] [query]",
	}}
	// TODO handle json flag
	c.fs.BoolVar(&c.json, "json", false, "json output")
	c.fs.BoolVar(&c.byTag, "t", false, "organize by tags")
	c.fs.BoolVar(&c.description, "desc", false, "show todo description in output")
	c.fs.BoolVar(&c.multiline, "2l", false, "use 2-line output for dates and description")
	c.fs.StringVar(&c.listFlag, "l", "", "show only todos from specified `list`")
	c.fs.BoolVar(&c.allLists, "a", false, "show todos from all lists (overrides -l)")
	c.fs.StringVar(&c.sortOption, "s", "", "sort todos by `field`: PRIO, DUE, STATUS, CREATED")
	c.fs.StringVar(&c.statusFilter, "S", "", "show only todos with specified `status`: NEEDS-ACTION, COMPLETED, CANCELLED, ANY")
	return c
}

type ListCmd struct {
	Cmd
	json         bool
	byTag        bool
	multiline    bool
	description  bool
	allLists     bool
	sortOption   string
	statusFilter string
}

type sortOption string

const (
	sortOptionStatus  sortOption = "STATUS"
	sortOptionPrio    sortOption = "PRIO"
	sortOptionDue     sortOption = "DUE"
	sortOptionCreated sortOption = "CREATED"
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
		return fmt.Errorf("Unknown status filter: %q, see %q", c.statusFilter, "tdx list -h")
	}

	// process sort option
	if c.sortOption == "" {
		c.sortOption = c.conf.DefaultSort
	}
	switch sortOption(c.sortOption) {
	case sortOptionStatus, sortOptionPrio, sortOptionDue, sortOptionCreated:
		break
	default:
		return fmt.Errorf("Unknown sort option: %q, see %q", c.sortOption, "tdx list -h")
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

	filterAndSortItems := func(ii []*vdir.Item) (items []*vdir.Item, err error) {
		items, err = filterByStatus(ii, vdir.ToDoStatus(c.statusFilter))
		if err != nil {
			return
		}
		items, err = filterByQuery(items, query)
		if err != nil {
			return
		}

		switch sortOption(c.sortOption) {
		case sortOptionPrio:
			sort.Sort(vdir.ByPriority(items))
		case sortOptionDue:
			sort.Sort(vdir.ByDue(items))
		case sortOptionStatus:
			sort.Sort(vdir.ByStatus(items))
		case sortOptionCreated:
			sort.Sort(vdir.ByCreated(items))
		}

		return
	}

	var m = make(map[string][]*vdir.Item)
	if c.byTag {
		emptyTag := vdir.Tag("[no tags]")

		allItems := []*vdir.Item{}
		for _, items := range collections {
			allItems = append(allItems, items...)
		}

		allItems, err := filterAndSortItems(allItems)
		if err != nil {
			return err
		}

		for _, item := range allItems {
			tags, err := item.Tags()
			if err != nil {
				return err
			}
			if len(tags) > 0 {
				for _, tag := range tags {
					m[tag.String()] = append(m[tag.String()], item)
				}
			} else {
				m[emptyTag.String()] = append(m[emptyTag.String()], item)
			}
		}
	} else {
		for col, items := range collections {
			items, err := filterAndSortItems(items)
			if err != nil {
				return err
			}

			for _, item := range items {
				m[col.String()] = append(m[col.String()], item)
			}
		}

	}

	if len(m) == 0 {
		return fmt.Errorf("No todos found")
	}

	// sort map keys and prepare output
	keys := []string{}
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	col := color.New(color.Bold, color.FgYellow).SprintFunc()

	var sb = strings.Builder{}
	for _, key := range keys {
		sb.WriteString(col(fmt.Sprintf("== %s\n", key)))
		for _, i := range m[key] {
			if err := writeItem(&sb, *c, *i); err != nil {
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
