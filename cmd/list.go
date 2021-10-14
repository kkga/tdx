package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/kkga/tdx/vdir"
)

type listOptions struct {
	lists         []string
	listsExcluded []string
	allLists      bool
	sorting       string
	group         string
	due           int
	status        string
	tags          []string
	tagsExcluded  []string
	description   bool
	multiline     bool

	// TODO handle json flag
	// json = listCmd.Flags().Bool("json", false, "output in json")
}

// envListOptsVar is the environment variable for setting default list options
const envListOptsVar = "TDX_LIST_OPTS"

type sortOption string

const (
	sortOptionStatus  sortOption = "STATUS"
	sortOptionPrio    sortOption = "PRIO"
	sortOptionDue     sortOption = "DUE"
	sortOptionCreated sortOption = "CREATED"
)

type groupOption string

const (
	groupOptionList groupOption = "LIST"
	groupOptionTag  groupOption = "TAG"
	groupOptionNone groupOption = "NONE"
)

func NewListCmd() *cobra.Command {
	opts := &listOptions{
		description: false,
	}

	cmd := &cobra.Command{
		Use:     "list [flags] [query]",
		Aliases: []string{"ls", "l"},
		Short:   "List todos",
		Long:    "List todos, optionally filtered by query.",
		Example: heredoc.Doc(`
            $ tdx list --sort prio --due 2
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			vd := make(vdir.Vdir)
			if err := vd.Init(vdirPath); err != nil {
				return err
			}

			defaultOpts := os.Getenv(envListOptsVar)
			cmd.ParseFlags(strings.Split(defaultOpts, " "))

			if err := checkGroupFlag(opts.group); err != nil {
				return err
			}

			if err := checkStatusFlag(opts.status); err != nil {
				return err
			}

			if err := checkSortFlag(opts.sorting); err != nil {
				return err
			}

			// if lists flag set, delete other collections from vdir
			if len(opts.lists) > 0 && opts.allLists == false {
				for _, list := range opts.lists {
					if err := checkList(vd, list, false); err != nil {
						return err
					}
				}
				for col := range vd {
					if !containsString(opts.lists, col.Name) {
						delete(vd, col)
					}
				}
			}

			query := strings.Join(args, "")

			return runList(vd, query, opts)
		},
	}

	cmd.Flags().StringSliceVarP(&opts.lists, "lists", "l", []string{}, "filter by `LISTS`, comma-separated (e.g. 'tasks,other')")
	cmd.Flags().StringVarP(&opts.group, "group", "g", "list", "group listed todos, valid options: list, tag, none ")
	cmd.Flags().BoolVarP(&opts.allLists, "all", "a", false, "show todos from all lists (overrides -l)")
	cmd.Flags().IntVarP(&opts.due, "due", "d", 0, "filter by due date in next `N` days")
	cmd.Flags().StringVarP(&opts.status, "status", "S", "needs-action", "filter by `STATUS`: needs-action, completed, cancelled, any")
	cmd.Flags().StringSliceVarP(&opts.tags, "tag", "t", []string{}, "filter todos by given `TAGS`")
	cmd.Flags().StringSliceVarP(&opts.tagsExcluded, "no-tag", "T", []string{}, "exclude todos with given `TAGS`")
	cmd.Flags().StringVarP(&opts.sorting, "sort", "s", "prio", "sort by `FIELD`: prio, due, status, created")
	cmd.Flags().BoolVar(&opts.description, "description", false, "show description in output")
	cmd.Flags().BoolVar(&opts.multiline, "two-line", false, "use 2-line output for dates and description")
	cmd.Flags().SortFlags = false
	return cmd
}

func runList(vd vdir.Vdir, query string, opts *listOptions) error {

	filterItems := func(items []*vdir.Item) (filtered []*vdir.Item, err error) {
		filtered = items

		filtered, err = vdir.Filter(vdir.ByStatus(filtered), vdir.ToDoStatus(opts.status))
		if err != nil {
			return
		}

		tags := []vdir.Tag{}
		for _, tag := range opts.tags {
			tags = append(tags, vdir.Tag(tag))
		}
		filtered, err = vdir.Filter(vdir.ByTags(filtered), tags)
		if err != nil {
			return
		}

		tagsExcluded := []vdir.Tag{}
		for _, tag := range opts.tagsExcluded {
			tagsExcluded = append(tagsExcluded, vdir.Tag(tag))
		}
		filtered, err = vdir.Filter(vdir.ByTagsExcluded(filtered), tagsExcluded)
		if err != nil {
			return
		}

		filtered, err = vdir.Filter(vdir.ByDue(filtered), opts.due)
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
		switch sortOption(strings.ToUpper(opts.sorting)) {
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

	switch groupOption(strings.ToUpper(opts.group)) {
	case groupOptionList:
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
	case groupOptionTag:
		items := []*vdir.Item{}
		for _, ii := range vd {
			items = append(items, ii...)
		}

		items, err := filterItems(items)
		if err != nil {
			return err
		}

		sortItems(items)

		emptyTag := vdir.Tag("[no tags]")

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
	case groupOptionNone:
		items := []*vdir.Item{}
		for _, ii := range vd {
			items = append(items, ii...)
		}

		items, err := filterItems(items)
		if err != nil {
			return err
		}

		sortItems(items)

		noneKey := groupOptionNone
		for _, item := range items {
			// here comes an ugly hack
			m[string(noneKey)] = append(m[string(noneKey)], item)
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

	colGroup := color.New(color.Bold, color.FgGreen).SprintFunc()
	var sb = strings.Builder{}
	for _, key := range keys {
		if key != string(groupOptionNone) {
			sb.WriteString(colGroup(fmt.Sprintf("-- %s --\n", key)))
		}
		for _, i := range m[key] {
			if err := writeItem(&sb, *i, opts); err != nil {
				return err
			}
		}
	}

	fmt.Print(sb.String())

	return nil
}

func writeItem(sb *strings.Builder, item vdir.Item, opts *listOptions) error {
	formatOpts := []vdir.FormatOption{}
	if opts.multiline {
		formatOpts = append(formatOpts, vdir.FormatMultiline)
	}
	if opts.description {
		formatOpts = append(formatOpts, vdir.FormatDescription)
	}

	s, err := item.Format(formatOpts...)
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

func checkGroupFlag(flag string) error {
	switch groupOption(strings.ToUpper(flag)) {
	case "":
		return nil
	case groupOptionList, groupOptionTag, groupOptionNone:
		return nil
	default:
		return fmt.Errorf("Unknown group option: %q, see %q", flag, "tdx list -h")
	}
}
