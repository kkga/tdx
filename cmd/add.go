package cmd

import (
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/kkga/tdx/vdir"
	"github.com/spf13/cobra"
)

type addOptions struct {
	list        string
	description string
}

func NewAddCmd() *cobra.Command {
	opts := &addOptions{}

	cmd := &cobra.Command{
		Use:     "add [options] <todo>",
		Aliases: []string{"a"},
		Short:   "Add todo",
		Long:    "Add new todo.",
		Args:    cobra.MinimumNArgs(1),
		Example: heredoc.Doc(`
            $ tdx add --sort prio --due 2
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			vd := make(vdir.Vdir)
			if err := vd.Init(vdirPath); err != nil {
				return err
			}

			todo := strings.Join(args, "")

			return runAdd(vd, opts, todo)
		},
	}

	cmd.Flags().StringVar(&opts.list, "l", "", "`list` for new todo")
	cmd.Flags().StringVar(&opts.description, "d", "", "`description` text")

	return cmd
}

func runAdd(vd vdir.Vdir, opts *addOptions, todo string) error {
	return nil
}

// func NewAddCmd() *AddCmd {
// 	c := &AddCmd{Cmd: Cmd{
// 		fs:        flag.NewFlagSet("add", flag.ExitOnError),
// 		alias:     []string{"a"},
// 		short:     "Add new todo",
// 		usageLine: "[options] <todo>",
// 		long: `AUTOMATIC PROPERTY PARSING
//   due date
//         If todo text contains a date in any of the following
//         forms, it will be converted to due date:
//         - "today", "tomorrow", "next tuesday"
//         - "in 3 days", "in a few days", "in 2 weeks", "in a month"
//         - "december 1st", "15 nov", "jul"
//   priority
//         If todo text contains one or more "!" chars,
//         they will be converted to priority:
//         - "!!!" (high)
//         - "!!"  (medium)
//         - "!"   (low)

// ENVIRONMENT VARIABLES
//   TDX_ADD_OPTS
//         default options for <add> command;
//         example: use a default list for new todos...
//             TDX_ADD_OPTS='-l=myList'`,
// 	}}
// 	return c
// }

// type AddCmd struct {
// 	Cmd
// 	description string
// 	list        string
// 	due         int
// }

// func (c *AddCmd) Run() error {
// 	if len(c.conf.AddOpts) > 0 {
// 		c.fs.Parse(strings.Split(c.conf.AddOpts, " "))
// 	}

// 	if err := c.fs.Parse(c.args); err != nil {
// 		return err
// 	}

// 	var collection *vdir.Collection
// 	if len(c.vdir) > 1 {
// 		// if err := checkList(c.list, vd, true, list); err != nil {
// 		// 	return err
// 		// }
// 		for col := range c.vdir {
// 			if col.Name == c.list {
// 				collection = col
// 			}
// 		}
// 	} else {
// 		// if only one collection, use it without requiring a list flag
// 		for col := range c.vdir {
// 			collection = col
// 		}
// 	}

// 	args := c.fs.Args()
// 	if len(args) == 0 {
// 		return errors.New("Provide a todo text")
// 	}

// 	cal := ical.NewCalendar()
// 	t := ical.NewComponent(ical.CompToDo)
// 	uid := vdir.GenerateUID()
// 	t.Props.SetText(ical.PropStatus, string(vdir.StatusNeedsAction))
// 	t.Props.SetText(ical.PropUID, uid)
// 	cal.Children = append(cal.Children, t)

// 	summary := strings.Join(args, " ")

// 	if strings.Contains(summary, "!!!") {
// 		summary = strings.Trim(strings.Replace(summary, "!!!", "", 1), " ")
// 		prioProp := ical.NewProp(ical.PropPriority)
// 		prioProp.Value = fmt.Sprint(vdir.PriorityHigh)
// 		t.Props.Add(prioProp)
// 	} else if strings.Contains(summary, "!!") {
// 		summary = strings.Trim(strings.Replace(summary, "!!", "", 1), " ")
// 		prioProp := ical.NewProp(ical.PropPriority)
// 		prioProp.Value = fmt.Sprint(vdir.PriorityMedium)
// 		t.Props.Add(prioProp)
// 	} else if strings.Contains(summary, "!") {
// 		summary = strings.Trim(strings.Replace(summary, "!", "", 1), " ")
// 		prioProp := ical.NewProp(ical.PropPriority)
// 		prioProp.Value = fmt.Sprint(vdir.PriorityLow)
// 		t.Props.Add(prioProp)
// 	}

// 	if c.description != "" {
// 		t.Props.SetText(ical.PropDescription, c.description)
// 	}

// 	if due, text, err := parseDate(summary); err == nil {
// 		t.Props.SetDateTime(ical.PropDue, due)
// 		summary = strings.Trim(strings.Replace(summary, text, "", 1), " ")
// 	}

// 	t.Props.SetText(ical.PropSummary, summary)

// 	p := path.Join(collection.Path, fmt.Sprintf("%s.ics", uid))

// 	item := &vdir.Item{
// 		Path: p,
// 		Ical: cal,
// 	}
// 	item.WriteFile()

// 	if err := c.vdir.Init(c.conf.Path); err != nil {
// 		return err
// 	}

// 	addedItem, err := c.vdir.ItemByPath(p)
// 	if err != nil {
// 		return err
// 	}

// 	s, err := addedItem.Format()
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Print(s)

// 	return nil
// }
