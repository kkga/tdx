package cmd

import (
	"errors"
	"flag"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/emersion/go-ical"

	"github.com/kkga/tdx/vdir"
	"github.com/tj/go-naturaldate"
)

func NewAddCmd() *AddCmd {
	c := &AddCmd{Cmd: Cmd{
		// TODO add long description of due date parsing
		fs:           flag.NewFlagSet("add", flag.ExitOnError),
		alias:        []string{"a"},
		shortDesc:    "Add todo",
		usageLine:    "[options] <todo>",
		listRequired: true,
	}}
	c.fs.StringVar(&c.listFlag, "l", "", "`list` for new todo")
	c.fs.StringVar(&c.priority, "p", "", "`priority`: !!!, !!, !")
	c.fs.StringVar(&c.description, "d", "", "`description` text")
	return c
}

type AddCmd struct {
	Cmd
	priority    string
	description string
}

func (c *AddCmd) Run() error {
	args := c.fs.Args()

	if len(args) == 0 {
		return errors.New("Provide a todo text")
	}

	summary := strings.Join(args, " ")

	t := ical.NewComponent("VTODO")
	t.Props.SetText(ical.PropSummary, summary)
	t.Props.SetText(ical.PropStatus, string(vdir.StatusNeedsAction))

	uid := vdir.GenerateUID()
	t.Props.SetText(ical.PropUID, uid)

	now := time.Now()
	due, _ := naturaldate.Parse(summary, now, naturaldate.WithDirection(naturaldate.Future))
	if due != now {
		t.Props.SetDateTime(ical.PropDue, due)
	}

	if c.description != "" {
		t.Props.SetText(ical.PropDescription, c.description)
	}

	if c.priority != "" {
		prioProp := ical.NewProp(ical.PropPriority)
		prioProp.SetValueType(ical.ValueInt)
		switch c.priority {
		case "!!!":
			prioProp.Value = fmt.Sprint(vdir.PriorityHigh)
			t.Props.Add(prioProp)
		case "!!":
			prioProp.Value = fmt.Sprint(vdir.PriorityMedium)
			t.Props.Add(prioProp)
		case "!":
			prioProp.Value = fmt.Sprint(vdir.PriorityLow)
			t.Props.Add(prioProp)
		default:
			return fmt.Errorf("Unknown priority flag: %s, expected one of: !!!, !!, !", c.priority)
		}
	}

	cal := ical.NewCalendar()
	cal.Children = append(cal.Children, t)

	p := path.Join(c.collection.Path, fmt.Sprintf("%s.ics", uid))

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
