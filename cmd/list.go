package cmd

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/emersion/go-ical"

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
	c.fs.BoolVar(&c.byTag, "t", false, "organize by tags")
	c.fs.BoolVar(&c.description, "desc", false, "show description in output")
	c.fs.BoolVar(&c.multiline, "2l", false, "use 2-line output for dates and description")
	c.fs.StringVar(&c.listFilter, "l", "", "filter by `list`")
	c.fs.BoolVar(&c.allLists, "a", false, "show all lists (overrides -l)")
	c.fs.StringVar(&c.sortOption, "s", "prio", "sort by `field`: prio, due, status, created")
	c.fs.StringVar(&c.statusFilter, "S", "needs-action", "filter by `status`: needs-action, completed, cancelled, any")
	c.fs.IntVar(&c.dueFilter, "d", 0, "filter by due date in next N `days`")
	return c
}

type ListCmd struct {
	Cmd
	// json         bool
	byTag        bool
	multiline    bool
	description  bool
	allLists     bool
	dueFilter    int
	listFilter   string
	statusFilter string
	sortOption   string
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

	filterAndSortItems := func(ii []*vdir.Item) (items []*vdir.Item, err error) {
		items, err = filterByStatus(ii, vdir.ToDoStatus(strings.ToUpper(c.statusFilter)))
		if err != nil {
			return
		}
		items, err = filterByDue(items, c.dueFilter)
		if err != nil {
			return
		}
		items, err = filterByQuery(items, query)
		if err != nil {
			return
		}

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

		return
	}

	var m = make(map[string][]*vdir.Item)

	if c.byTag {
		emptyTag := vdir.Tag("[no tags]")

		allItems := []*vdir.Item{}
		for _, items := range vd {
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
		for col, items := range vd {
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
		sb.WriteString(col(fmt.Sprintf("-- %s\n", key)))
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
		vt, err := i.Vtodo()
		if err != nil {
			return nil, err
		}
		s, propErr := vt.Props.Text(ical.PropStatus)
		if propErr != nil {
			return nil, propErr
		}
		if s == string(status) {
			filtered = append(filtered, i)
		}
	}
	return
}

func filterByDue(items []*vdir.Item, dueDays int) (filtered []*vdir.Item, err error) {
	if dueDays == 0 {
		return items, nil
	}
	now := time.Now()
	inDueDays := now.AddDate(0, 0, dueDays)

	for _, i := range items {
		vt, err := i.Vtodo()
		if err != nil {
			return nil, err
		}
		due, err := vt.Props.DateTime(ical.PropDue, time.Local)
		if err != nil {
			return nil, err
		}
		if !due.IsZero() && due.Before(inDueDays) {
			filtered = append(filtered, i)
		}
	}
	return
}

func filterByQuery(items []*vdir.Item, query string) (filtered []*vdir.Item, err error) {
	if query == "" {
		return items, nil
	}

	for _, i := range items {
		vt, err := i.Vtodo()
		if err != nil {
			return nil, err
		}
		summary, propErr := vt.Props.Text(ical.PropSummary)
		if propErr != nil {
			return nil, propErr
		}
		if strings.Contains(strings.ToLower(summary), strings.ToLower(query)) {
			filtered = append(filtered, i)
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
