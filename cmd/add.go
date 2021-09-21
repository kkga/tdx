package cmd

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/kkga/ctodo/vtodo"
)

func NewAddCmd() *AddCmd {
	c := &AddCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("add", flag.ExitOnError),
		alias:     []string{"a"},
		shortDesc: "Add todo",
		usageLine: "[options]",
	}}
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

	if len(args) == 0 {
		return errors.New("Provide a todo summary")
	}

	summary := strings.Join(args, " ")
	uid := vtodo.GenerateUID()

	t := ical.NewComponent("VTODO")
	t.Props.SetText(ical.PropSummary, summary)
	t.Props.SetText(ical.PropUID, uid)
	t.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())

	// TODO parse due date flag

	if c.description != "" {
		t.Props.SetText(ical.PropDescription, c.description)
	}

	if c.priority != "" {
		prioProp := ical.NewProp(ical.PropPriority)
		prioProp.SetValueType(ical.ValueInt)
		switch c.priority {
		case "high":
			prioProp.Value = fmt.Sprint(vtodo.PriorityHigh)
			t.Props.Add(prioProp)
		case "medium":
			prioProp.Value = fmt.Sprint(vtodo.PriorityMedium)
			t.Props.Add(prioProp)
		case "low":
			prioProp.Value = fmt.Sprint(vtodo.PriorityLow)
			t.Props.Add(prioProp)
		default:
			return fmt.Errorf("Unknown priority flag: %s, expected one of: high, medium, low", c.priority)
		}
	}

	calBuf, err := vtodo.Encode(t)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.ics", calDir, uid))
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.Write(calBuf.Bytes())
	if err != nil {
		return err
	}

	w.Flush()
	if err != nil {
		return err
	}

	return nil
}
