package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
	"github.com/spf13/cobra"
)

type addOptions struct {
	list        string
	description string
}

// envAddOptsVar is the environment variable for setting default add options
const envAddOptsVar = "TDX_ADD_OPTS"

func NewAddCmd() *cobra.Command {
	opts := &addOptions{}

	cmd := &cobra.Command{
		Use:        "add <todo summary>",
		Aliases:    []string{"a"},
		Short:      "Add todo",
		Long:       "Add new todo.",
		SuggestFor: []string{"new"},
		Args:       cobra.MinimumNArgs(1),
		Example: heredoc.Doc(`
			$ tdx add buy milk -l shopping`),
		RunE: func(cmd *cobra.Command, args []string) error {
			vd := make(vdir.Vdir)
			if err := vd.Init(vdirPath); err != nil {
				return err
			}

			if err := checkList(vd, opts.list, true); err != nil {
				return err
			}

			var collection *vdir.Collection
			for col := range vd {
				if col.Name == opts.list {
					collection = col
				}
			}

			rawTodo := strings.Join(args, " ")

			return runAdd(collection, opts, rawTodo)
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVarP(&opts.list, "list", "l", "", "`LIST` for new todo")
	cmd.MarkFlagRequired("list")
	cmd.Flags().StringVarP(&opts.description, "description", "d", "", "description text")

	defaultOpts := os.Getenv(envAddOptsVar)
	cmd.ParseFlags(strings.Split(defaultOpts, " "))

	return cmd
}

func runAdd(collection *vdir.Collection, opts *addOptions, rawTodo string) error {
	cal := ical.NewCalendar()
	t := ical.NewComponent(ical.CompToDo)
	uid := vdir.GenerateUID()
	t.Props.SetText(ical.PropStatus, string(vdir.StatusNeedsAction))
	t.Props.SetText(ical.PropUID, uid)
	cal.Children = append(cal.Children, t)

	summary := rawTodo

	if strings.Contains(summary, "!!!") {
		summary = strings.Trim(strings.Replace(summary, "!!!", "", 1), " ")
		prioProp := ical.NewProp(ical.PropPriority)
		prioProp.Value = fmt.Sprint(vdir.PriorityHigh)
		t.Props.Add(prioProp)
	} else if strings.Contains(summary, "!!") {
		summary = strings.Trim(strings.Replace(summary, "!!", "", 1), " ")
		prioProp := ical.NewProp(ical.PropPriority)
		prioProp.Value = fmt.Sprint(vdir.PriorityMedium)
		t.Props.Add(prioProp)
	} else if strings.Contains(summary, "!") {
		summary = strings.Trim(strings.Replace(summary, "!", "", 1), " ")
		prioProp := ical.NewProp(ical.PropPriority)
		prioProp.Value = fmt.Sprint(vdir.PriorityLow)
		t.Props.Add(prioProp)
	}

	if due, text, err := parseDate(summary); err == nil {
		t.Props.SetDateTime(ical.PropDue, due)
		summary = strings.Trim(strings.Replace(summary, text, "", 1), " ")
	}

	if opts.description != "" {
		t.Props.SetText(ical.PropDescription, opts.description)
	}

	t.Props.SetText(ical.PropSummary, summary)

	p := path.Join(collection.Path, fmt.Sprintf("%s.ics", uid))

	item := &vdir.Item{
		Path: p,
		Ical: cal,
	}
	item.WriteFile()

	vd := make(vdir.Vdir)
	if err := vd.Init(vdirPath); err != nil {
		return err
	}

	addedItem, err := vd.ItemByPath(p)
	if err != nil {
		return err
	}

	s, err := addedItem.Format()
	if err != nil {
		return err
	}
	fmt.Print(s)

	return nil
}
