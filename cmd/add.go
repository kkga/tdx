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
		fs:    flag.NewFlagSet("add", flag.ExitOnError),
		alias: []string{"a"},
		short: "Add new todo",
		long: `DUE DATE
  If todo text contains a date in one of the following
  formats, it will be applied as due date:
  * "today", "tomorrow", "in 3 days", "in 2 weeks"
  * "next week", "next month", "next monday"
  * ordinal date: "december 1st", "15th november"

PRIORITY
  If todo text contains one or more "!" chars,
  they will be converted to priority:
  * "!!!" (high)
  * "!!"  (medium)
  * "!"   (low)`,
		usageLine:    "[options] <todo>",
		listRequired: true,
	}}
	c.fs.StringVar(&c.listFlag, "l", "", "`list` for new todo")
	c.fs.StringVar(&c.description, "d", "", "`description` text")
	return c
}

type AddCmd struct {
	Cmd
	description string
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

	t.Props.SetText(ical.PropSummary, summary)

	if c.description != "" {
		t.Props.SetText(ical.PropDescription, c.description)
	}

	now := time.Now()
	due, _ := naturaldate.Parse(summary, now, naturaldate.WithDirection(naturaldate.Future))
	if due != now {
		t.Props.SetDateTime(ical.PropDue, due)
	}

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
