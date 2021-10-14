package cmd

import (
	"errors"
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

			defaultOpts := os.Getenv(envAddOptsVar)
			cmd.ParseFlags(strings.Split(defaultOpts, " "))

			var collection *vdir.Collection
			if len(vd) > 1 {
				if err := checkList(vd, opts.list, true); err != nil {
					return err
				}
				for col := range vd {
					if col.Name == opts.list {
						collection = col
					}
				}
			} else {
				// if only one collection, use it without requiring a list flag
				for col := range vd {
					collection = col
				}
			}

			todo := strings.Join(args, "")

			return runAdd(collection, opts, todo)
		},
	}

	cmd.Flags().StringVarP(&opts.list, "list", "l", "", "`list` for new todo")
	cmd.MarkFlagRequired("list")
	cmd.Flags().StringVar(&opts.description, "d", "", "`description` text")

	return cmd
}

func runAdd(collection *vdir.Collection, opts *addOptions, todo string) error {
	return nil
}

func (c *AddCmd) Run() error {

	args := c.fs.Args()
	if len(args) == 0 {
		return errors.New("Provide a todo text")
	}

	cal := ical.NewCalendar()
	t := ical.NewComponent(ical.CompToDo)
	uid := vdir.GenerateUID()
	t.Props.SetText(ical.PropStatus, string(vdir.StatusNeedsAction))
	t.Props.SetText(ical.PropUID, uid)
	cal.Children = append(cal.Children, t)

	summary := strings.Join(args, " ")

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

	if c.description != "" {
		t.Props.SetText(ical.PropDescription, c.description)
	}

	if due, text, err := parseDate(summary); err == nil {
		t.Props.SetDateTime(ical.PropDue, due)
		summary = strings.Trim(strings.Replace(summary, text, "", 1), " ")
	}

	t.Props.SetText(ical.PropSummary, summary)

	p := path.Join(collection.Path, fmt.Sprintf("%s.ics", uid))

	item := &vdir.Item{
		Path: p,
		Ical: cal,
	}
	item.WriteFile()

	if err := c.vdir.Init(c.conf.Path); err != nil {
		return err
	}

	addedItem, err := c.vdir.ItemByPath(p)
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
