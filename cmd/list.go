package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/kkga/tdx/vdir"
)

var (
	lists         []string
	listsExcluded []string
	allLists      bool

	sorting      string
	due          int
	status       string
	tags         []string
	tagsExcluded []string
	description  bool
	multiline    bool
	byTag        bool

	// TODO think of a better way to organize output by tags/lists/...
	// byTag = listCmd.Flags().Bool("by-tag", false, "organize by tags")
	// TODO handle json flag
	// json = listCmd.Flags().Bool("json", false, "output in json")
)

type sortOption string

const (
	sortOptionStatus  sortOption = "STATUS"
	sortOptionPrio    sortOption = "PRIO"
	sortOptionDue     sortOption = "DUE"
	sortOptionCreated sortOption = "CREATED"
)

func init() {
	listCmd.Flags().StringSliceVarP(&lists, "lists", "l", []string{}, "filter by `LISTS`, comma-separated (e.g. 'tasks,other')")
	listCmd.Flags().BoolVarP(&allLists, "all", "a", false, "show todos from all lists (overrides -l)")
	listCmd.Flags().IntVarP(&due, "due", "d", 0, "filter by due date in next `N` days")
	listCmd.Flags().StringVarP(&status, "status", "S", "needs-action", "filter by `STATUS`: needs-action, completed, cancelled, any")
	listCmd.Flags().StringSliceVarP(&tags, "tag", "t", []string{}, "filter todos by given `TAGS`")
	listCmd.Flags().StringSliceVarP(&tagsExcluded, "no-tag", "T", []string{}, "exclude todos with given `TAGS`")

	listCmd.Flags().StringVarP(&sorting, "sort", "s", "prio", "sort by `FIELD`: prio, due, status, created")
	listCmd.Flags().BoolVar(&description, "description", false, "show description in output")
	listCmd.Flags().BoolVar(&byTag, "by-tag", false, "organize by tags")
	listCmd.Flags().BoolVar(&multiline, "two-line", false, "use 2-line output for dates and description")

	listCmd.Flags().SortFlags = false
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list [flags] [query]",
	Aliases: []string{"ls", "l"},
	Short:   "List todos",
	Long:    "List todos, optionally filtered by query.",
	RunE: func(cmd *cobra.Command, args []string) error {
		vd := make(vdir.Vdir)
		if err := vd.Init(vdirPath); err != nil {
			return err
		}
		return runList(vd, cmd.Flags(), args)
	},
}

func runList(vd vdir.Vdir, flags *flag.FlagSet, args []string) error {
	var query string
	if len(args) > 0 {
		query = strings.Join(args, "")
	}

	// TODO: this is a good place to set default opts from env?
	// if len(c.conf.ListOpts) > 0 {
	// 	c.fs.Parse(strings.Split(c.conf.ListOpts, " "))
	// }

	// if err := c.fs.Parse(c.args); err != nil {
	// 	return err
	// }

	// if err := checkStatusFlag(*status); err != nil {
	// 	return err
	// }

	// if err := checkSortFlag(*sorting); err != nil {
	// 	return err
	// }

	// if list flag set, delete other collections from vdir
	collections := vd
	if len(lists) > 0 && allLists == false {
		for _, list := range lists {
			if err := checkList(collections, list, false); err != nil {
				return err
			}
		}
		for col := range collections {
			if !containsString(lists, col.Name) {
				delete(collections, col)
			}
		}
	}

	filterItems := func(items []*vdir.Item) (filtered []*vdir.Item, err error) {
		filtered = items

		filtered, err = vdir.Filter(vdir.ByStatus(filtered), vdir.ToDoStatus(status))
		if err != nil {
			return
		}
		// filtered, err = vdir.Filter(vdir.ByTag(filtered), tags)
		// if err != nil {
		// 	return
		// }
		// filtered, err = vdir.Filter(vdir.ByTagExclude(filtered), tagsExcluded)
		// if err != nil {
		// 	return
		// }
		filtered, err = vdir.Filter(vdir.ByDue(filtered), due)
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
		switch sortOption(strings.ToUpper(sorting)) {
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

	if byTag {
		emptyTag := vdir.Tag("[no tags]")

		items := []*vdir.Item{}
		for _, ii := range collections {
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
		for col, items := range collections {
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
			if err := writeItem(&sb, *i); err != nil {
				return err
			}
		}
	}

	fmt.Print(sb.String())

	return nil
}

func writeItem(sb *strings.Builder, item vdir.Item) error {
	opts := []vdir.FormatOption{}
	if multiline {
		opts = append(opts, vdir.FormatMultiline)
	}
	if description {
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
