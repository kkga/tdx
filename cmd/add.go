package cmd

import (
	"errors"
	"flag"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	// "github.com/kkga/tdx/vdir"
	"github.com/kkga/tdx/vdir"
)

func NewAddCmd() *AddCmd {
	c := &AddCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("add", flag.ExitOnError),
		alias:     []string{"a"},
		shortDesc: "Add todo",
		usageLine: "[options]",
		listReq:   true,
	}}
	c.fs.StringVar(&c.list, "l", "", "list")
	c.fs.StringVar(&c.priority, "p", "", "priority")
	c.fs.StringVar(&c.due, "D", "", "due date")
	c.fs.StringVar(&c.description, "d", "", "description")
	return c
}

type AddCmd struct {
	Cmd
	priority    string
	due         string
	description string
}

func (c *AddCmd) Run() error {
	args := c.fs.Args()

	if c.list == "" {
		return errors.New("Specify a list with '-l' or set default list with 'TDX_DEFAULT_LIST'")
	}

	if len(args) == 0 {
		return errors.New("Provide a todo summary")
	}

	summary := strings.Join(args, " ")
	uid := vdir.GenerateUID()

	t := ical.NewComponent("VTODO")
	t.Props.SetText(ical.PropSummary, summary)
	t.Props.SetText(ical.PropUID, uid)
	t.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())
	t.Props.SetText(ical.PropStatus, string(vdir.StatusNeedsAction))

	// TODO parse due date flag

	if c.description != "" {
		t.Props.SetText(ical.PropDescription, c.description)
	}

	if c.priority != "" {
		prioProp := ical.NewProp(ical.PropPriority)
		prioProp.SetValueType(ical.ValueInt)
		switch c.priority {
		case "high":
			prioProp.Value = fmt.Sprint(vdir.PriorityHigh)
			t.Props.Add(prioProp)
		case "medium":
			prioProp.Value = fmt.Sprint(vdir.PriorityMedium)
			t.Props.Add(prioProp)
		case "low":
			prioProp.Value = fmt.Sprint(vdir.PriorityLow)
			t.Props.Add(prioProp)
		default:
			return fmt.Errorf("Unknown priority flag: %s, expected one of: high, medium, low", c.priority)
		}
	}

	cal := ical.NewCalendar()
	// TODO move this data somewhere
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN")
	cal.Children = append(cal.Children, t)

	p := path.Join(c.collection.Path, fmt.Sprintf("%s.ics", uid))

	item := &vdir.Item{
		Path: p,
		Ical: cal,
	}
	item.WriteFile()

	fmt.Printf("Added: %s\n", summary)

	return nil
}
