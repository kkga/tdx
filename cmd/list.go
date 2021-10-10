package cmd

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/kkga/tdx/vdir"
)

func NewListCmd() *ListCmd {
	c := &ListCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("list", flag.ExitOnError),
		alias:     []string{"ls", "l"},
		usageLine: "[options] [query]",
		short:     "List todos, optionally filtered by query",
		long: `ENVIRONMENT VARIABLES
  TDX_LIST_OPTS
        default options for <list> command;
        example: filter by due date in next 2 days, from 'myList', organized by tags...
            TDX_LIST_OPTS='-d=2 -l=myList -t'`,
	}}
	// TODO handle json flag
	// c.fs.BoolVar(&c.json, "json", false, "json output")
	// c.fs.BoolVar(&c.byTag, "t", false, "organize by tags")
	c.fs.BoolVar(&c.description, "desc", false, "show description in output")
	c.fs.BoolVar(&c.multiline, "2l", false, "use 2-line output for dates and description")
	c.fs.StringVar(&c.listFilter, "l", "", "filter by `list`")
	c.fs.BoolVar(&c.allLists, "a", false, "show all lists (overrides -l)")
	c.fs.StringVar(&c.sortOption, "s", "prio", "sort by `field`: prio, due, status, created")
	c.fs.StringVar(&c.statusFilter, "S", "needs-action", "filter by `status`: needs-action, completed, cancelled, any")
	c.fs.IntVar(&c.dueFilter, "d", 0, "filter by due date in next N `days`")
	c.fs.StringVar(&c.tagFilter, "t", "", "filter todos by given tags")
	c.fs.StringVar(&c.tagExcludeFilter, "T", "", "exclude todos with given tags")
	return c
}

type ListCmd struct {
	Cmd
	// json         bool
	byTag            bool
	multiline        bool
	description      bool
	allLists         bool
	tagFilter        string
	tagExcludeFilter string
	dueFilter        int
	listFilter       string
	statusFilter     string
	sortOption       string
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

	if len(c.conf.ListOpts) > 0 {
		c.fs.Parse(strings.Split(c.conf.ListOpts, " "))
	}

	if err := c.fs.Parse(c.args); err != nil {
		return err
	}

	if err := checkStatusFlag(c.statusFilter); err != nil {
		return err
	}

	if err := checkSortFlag(c.sortOption); err != nil {
		return err
	}

	// if list flag set, delete other collections from vdir
	vd := c.vdir
	if c.listFilter != "" && c.allLists == false {
		if err := c.checkListFlag(c.listFilter, false, c); err != nil {
			return err
		}
		for col := range vd {
			if col.Name != c.listFilter {
				delete(vd, col)
			}
		}
	}

	filterItems := func(items []*vdir.Item) (filtered []*vdir.Item, err error) {
		filtered = items

		filtered, err = vdir.Filter(vdir.ByStatus(filtered), vdir.ToDoStatus(c.statusFilter))
		if err != nil {
			return
		}
		filtered, err = vdir.Filter(vdir.ByTag(filtered), vdir.Tag(c.tagFilter))
		if err != nil {
			return
		}
		filtered, err = vdir.Filter(vdir.ByDue(filtered), c.dueFilter)
		if err != nil {
			return
		}
		filtered, err = vdir.Filter(vdir.ByText(filtered), query)
		if err != nil {
			return
		}

		return

	}

	sortItems := func(items []*vdir.Item) {
		switch sortOption(strings.ToUpper(c.sortOption)) {
		case sortOptionPrio:
			sort.Sort(vdir.ByPriority(items))
		case sortOptionDue:
			sort.Sort(vdir.ByDue(items))
		case sortOptionStatus:
			sort.Sort(vdir.ByStatus(items))
		case sortOptionCreated:
			sort.Sort(vdir.ByCreated(items))
		}
	}

	var m = make(map[string][]*vdir.Item)

	if c.byTag {
		emptyTag := vdir.Tag("[no tags]")

		items := []*vdir.Item{}
		for _, ii := range vd {
			items = append(items, ii...)
		}

		items, err := filterItems(items)
		if err != nil {
			return err
		}

		sortItems(items)

		for _, item := range items {
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
		for col, items := range vd {
			items, err := filterItems(items)
			if err != nil {
				return err
			}
			sortItems(items)

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

	col := color.New(color.Bold, color.FgGreen).SprintFunc()

	var sb = strings.Builder{}
	for _, key := range keys {
		sb.WriteString(col(fmt.Sprintf("-- %s --\n", key)))
		for _, i := range m[key] {
			if err := writeItem(&sb, *c, *i); err != nil {
				return err
			}
		}
	}

	fmt.Print(sb.String())

	return nil
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

func checkStatusFlag(flag string) error {
	switch vdir.ToDoStatus(strings.ToUpper(flag)) {
	case "":
		return nil
	case vdir.StatusNeedsAction, vdir.StatusCompleted, vdir.StatusCancelled, vdir.StatusAny:
		return nil
	default:
		return fmt.Errorf("Unknown status filter: %q, see %q", flag, "tdx list -h")
	}
}

func checkSortFlag(flag string) error {
	switch sortOption(strings.ToUpper(flag)) {
	case "":
		return nil
	case sortOptionStatus, sortOptionPrio, sortOptionDue, sortOptionCreated:
		return nil
	default:
		return fmt.Errorf("Unknown sort option: %q, see %q", flag, "tdx list -h")
	}
}
