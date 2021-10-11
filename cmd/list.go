package cmd

import (
	"fmt"
	"strings"

	"github.com/kkga/tdx/vdir"
	"github.com/spf13/cobra"
)

var (
	list         = listCmd.Flags().StringP("list", "l", "", "filter by `LIST`")
	allLists     = listCmd.Flags().BoolP("all", "a", false, "show todos from all lists (overrides -l)")
	sort         = listCmd.Flags().StringP("sort", "s", "prio", "sort by `FIELD`: prio, due, status, created")
	due          = listCmd.Flags().IntP("due", "d", 0, "filter by due date in next N `DAYS`")
	status       = listCmd.Flags().StringP("status", "S", "needs-action", "filter by `STATUS`: needs-action, completed, cancelled, any")
	tags         = listCmd.Flags().StringP("tag", "t", "", "filter todos by given `TAGS`")
	tagsExcluded = listCmd.Flags().StringP("no-tag", "T", "", "exclude todos with given `TAGS`")
	description  = listCmd.Flags().Bool("description", false, "show description in output")
	multiline    = listCmd.Flags().Bool("two-line", false, "use 2-line output for dates and description")
	// TODO think of a better way to organize output by tags/lists/...
	byTag = listCmd.Flags().Bool("by-tag", false, "organize by tags")
	// TODO handle json flag
	json = listCmd.Flags().Bool("json", false, "output in json")
)

func init() {
	listCmd.Flags().SortFlags = false
	rootCmd.AddCommand(listCmd)

	// c.fs.BoolVar(&c.json, "json", false, "json output")
	// c.fs.BoolVar(&c.byTag, "t", false, "organize by tags")

	// listCmd.Flags().StringVarP(listFilter, "list", "l", "", "filter by `LIST`")
	// listCmd.Flags().BoolVarP(allLists, "all", "a", false, "show todos from all lists (overrides -l)")
	// listCmd.Flags().StringVarP(sortOption, "sort", "s", "prio", "sort by `FIELD`: prio, due, status, created")
	// listCmd.Flags().IntVarP(dueFilter, "due", "d", 0, "filter by due date in next N `DAYS`")
	// listCmd.Flags().StringVarP(statusFilter, "status", "S", "needs-action", "filter by `STATUS`: needs-action, completed, cancelled, any")
	// listCmd.Flags().StringVarP(tagFilter, "tag", "t", "", "filter todos by given `TAGS`")
	// listCmd.Flags().StringVarP(tagExcludeFilter, "no-tag", "T", "", "exclude todos with given `TAGS`")
	// listCmd.Flags().BoolVar(description, "description", false, "show description in output")
	// listCmd.Flags().BoolVar(multiline, "two-line", false, "use 2-line output for dates and description")
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List todos",
	Long:    "List todos, optionally filtered by query.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// func NewListCmd() *ListCmd {
// 	c := &ListCmd{Cmd: Cmd{
// 		fs:        flag.NewFlagSet("list", flag.ExitOnError),
// 		name:      "list",
// 		alias:     []string{"ls", "l"},
// 		usageLine: "[options] [query]",
// 		short:     "List todos, optionally filtered by query",
// 		long: `ENVIRONMENT VARIABLES
//   TDX_LIST_OPTS         default options for <list> command;
//                         example: filter by due date in next 2 days, from 'myList', organized by tags...
//                             TDX_LIST_OPTS='-d=2 -l=myList -t'`,
// 	}}
// 	c.fs.SortFlags = false

// 	// TODO handle json flag
// 	// c.fs.BoolVar(&c.json, "json", false, "json output")

// 	// TODO think of a better way to organize output by tags/lists/...
// 	// c.fs.BoolVar(&c.byTag, "t", false, "organize by tags")

// 	c.fs.StringVarP(&c.listFilter, "list", "l", "", "filter by `LIST`")
// 	c.fs.BoolVarP(&c.allLists, "all", "a", false, "show todos from all lists (overrides -l)")

// 	c.fs.StringVarP(&c.sortOption, "sort", "s", "prio", "sort by `FIELD`: prio, due, status, created")
// 	c.fs.IntVarP(&c.dueFilter, "due", "d", 0, "filter by due date in next N `DAYS`")
// 	c.fs.StringVarP(&c.statusFilter, "status", "S", "needs-action", "filter by `STATUS`: needs-action, completed, cancelled, any")

// 	c.fs.StringVarP(&c.tagFilter, "tag", "t", "", "filter todos by given `TAGS`")
// 	c.fs.StringVarP(&c.tagExcludeFilter, "no-tag", "T", "", "exclude todos with given `TAGS`")

// 	c.fs.BoolVar(&c.description, "description", false, "show description in output")
// 	c.fs.BoolVar(&c.multiline, "two-line", false, "use 2-line output for dates and description")

// 	// c.fs.BoolP("help", "h", false, "show help")

// 	return c
// }

// type ListCmd struct {
// 	Cmd
// 	// json         bool
// 	byTag            bool
// 	multiline        bool
// 	description      bool
// 	allLists         bool
// 	tagFilter        string
// 	tagExcludeFilter string
// 	dueFilter        int
// 	listFilter       string
// 	statusFilter     string
// 	sortOption       string
// }

type sortOption string

const (
	sortOptionStatus  sortOption = "STATUS"
	sortOptionPrio    sortOption = "PRIO"
	sortOptionDue     sortOption = "DUE"
	sortOptionCreated sortOption = "CREATED"
)

// func (c *ListCmd) Run() error {
// 	var query string

// 	if len(c.fs.Args()) > 0 {
// 		query = strings.Join(c.fs.Args(), "")
// 	}

// 	// if len(c.conf.ListOpts) > 0 {
// 	// 	c.fs.Parse(strings.Split(c.conf.ListOpts, " "))
// 	// }

// 	// if err := c.fs.Parse(c.args); err != nil {
// 	// 	return err
// 	// }

// 	if err := checkStatusFlag(c.statusFilter); err != nil {
// 		return err
// 	}

// 	if err := checkSortFlag(c.sortOption); err != nil {
// 		return err
// 	}

// 	// if list flag set, delete other collections from vdir
// 	vd := c.vdir
// 	if c.listFilter != "" && c.allLists == false {
// 		if err := c.checkListFlag(c.listFilter, false, c); err != nil {
// 			return err
// 		}
// 		for col := range vd {
// 			if col.Name != c.listFilter {
// 				delete(vd, col)
// 			}
// 		}
// 	}

// 	filterItems := func(items []*vdir.Item) (filtered []*vdir.Item, err error) {
// 		filtered = items

// 		filtered, err = vdir.Filter(vdir.ByStatus(filtered), vdir.ToDoStatus(c.statusFilter))
// 		if err != nil {
// 			return
// 		}
// 		filtered, err = vdir.Filter(vdir.ByTag(filtered), vdir.Tag(c.tagFilter))
// 		if err != nil {
// 			return
// 		}
// 		filtered, err = vdir.Filter(vdir.ByTagExclude(filtered), vdir.Tag(c.tagExcludeFilter))
// 		if err != nil {
// 			return
// 		}
// 		filtered, err = vdir.Filter(vdir.ByDue(filtered), c.dueFilter)
// 		if err != nil {
// 			return
// 		}
// 		filtered, err = vdir.Filter(vdir.ByText(filtered), query)
// 		if err != nil {
// 			return
// 		}

// 		return

// 	}

// 	sortItems := func(items []*vdir.Item) {
// 		switch sortOption(strings.ToUpper(c.sortOption)) {
// 		case sortOptionPrio:
// 			sort.Sort(vdir.ByPriority(items))
// 		case sortOptionDue:
// 			sort.Sort(vdir.ByDue(items))
// 		case sortOptionStatus:
// 			sort.Sort(vdir.ByStatus(items))
// 		case sortOptionCreated:
// 			sort.Sort(vdir.ByCreated(items))
// 		}
// 	}

// 	var m = make(map[string][]*vdir.Item)

// 	if c.byTag {
// 		emptyTag := vdir.Tag("[no tags]")

// 		items := []*vdir.Item{}
// 		for _, ii := range vd {
// 			items = append(items, ii...)
// 		}

// 		items, err := filterItems(items)
// 		if err != nil {
// 			return err
// 		}

// 		sortItems(items)

// 		for _, item := range items {
// 			tags, err := item.Tags()
// 			if err != nil {
// 				return err
// 			}
// 			if len(tags) > 0 {
// 				for _, tag := range tags {
// 					m[tag.String()] = append(m[tag.String()], item)
// 				}
// 			} else {
// 				m[emptyTag.String()] = append(m[emptyTag.String()], item)
// 			}
// 		}
// 	} else {
// 		for col, items := range vd {
// 			items, err := filterItems(items)
// 			if err != nil {
// 				return err
// 			}
// 			sortItems(items)

// 			for _, item := range items {
// 				m[col.String()] = append(m[col.String()], item)
// 			}
// 		}

// 	}

// 	if len(m) == 0 {
// 		return fmt.Errorf("No todos found")
// 	}

// 	// sort map keys and prepare output
// 	keys := []string{}
// 	for key := range m {
// 		keys = append(keys, key)
// 	}
// 	sort.Strings(keys)

// 	col := color.New(color.Bold, color.FgGreen).SprintFunc()

// 	var sb = strings.Builder{}
// 	for _, key := range keys {
// 		sb.WriteString(col(fmt.Sprintf("-- %s --\n", key)))
// 		for _, i := range m[key] {
// 			if err := writeItem(&sb, *c, *i); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	fmt.Print(sb.String())

// 	return nil
// }

// func writeItem(sb *strings.Builder, c ListCmd, item vdir.Item) error {
// 	opts := []vdir.FormatOption{}
// 	if c.multiline {
// 		opts = append(opts, vdir.FormatMultiline)
// 	}
// 	if c.description {
// 		opts = append(opts, vdir.FormatDescription)
// 	}

// 	s, err := item.Format(opts...)
// 	if err != nil {
// 		return err
// 	}

// 	sb.WriteString(s)

// 	return nil
// }

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
