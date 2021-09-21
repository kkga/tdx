package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-ical"
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

const (
	ToDoPriorityHigh   = 1
	ToDoPriorityMedium = 5
	ToDoPriorityLow    = 6
)

func (c *AddCmd) Run() error {
	args := c.fs.Args()

	if len(args) == 0 {
		return errors.New("Provide a todo summary")
	}

	summary := strings.Join(args, " ")
	uid := generateUID()

	vtodo := ical.NewComponent("VTODO")
	vtodo.Props.SetText(ical.PropSummary, summary)
	vtodo.Props.SetText(ical.PropUID, uid)
	vtodo.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())

	// TODO parse due date flag

	if c.description != "" {
		vtodo.Props.SetText(ical.PropDescription, c.description)
	}

	if c.priority != "" {
		prioProp := ical.NewProp(ical.PropPriority)
		prioProp.SetValueType(ical.ValueInt)
		switch c.priority {
		case "high":
			prioProp.Value = fmt.Sprint(ToDoPriorityHigh)
			vtodo.Props.Add(prioProp)
		case "medium":
			prioProp.Value = fmt.Sprint(ToDoPriorityMedium)
			vtodo.Props.Add(prioProp)
		case "low":
			prioProp.Value = fmt.Sprint(ToDoPriorityLow)
			vtodo.Props.Add(prioProp)
		default:
			return fmt.Errorf("Unknown priority flag: %s, expected one of: high, medium, low", c.priority)
		}
	}

	calBuf, err := encode(vtodo)
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

// generateUID returns a random string containing timestamp and hostname
func generateUID() string {
	sb := strings.Builder{}

	time := time.Now().UnixNano()

	randStr := func(n int) string {
		rs := rand.NewSource(time)
		r := rand.New(rs)
		var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyz")

		s := make([]rune, n)
		for i := range s {
			s[i] = letters[r.Intn(len(letters))]
		}
		return string(s)
	}

	sb.WriteString(fmt.Sprint(time))
	sb.WriteString(fmt.Sprintf("-%s", randStr(8)))
	if hostname, _ := os.Hostname(); hostname != "" {
		sb.WriteString(fmt.Sprintf("@%s", hostname))
	}

	return sb.String()
}

// encode adds vtodo into a new Calendar and returns a buffer ready for writing
func encode(vtodo *ical.Component) (*bytes.Buffer, error) {
	cal := ical.NewCalendar()
	// TODO move this data somewhere
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN")
	cal.Children = append(cal.Children, vtodo)

	var buf bytes.Buffer
	if err := ical.NewEncoder(&buf).Encode(cal); err != nil {
		return &buf, err
	}
	return &buf, nil
}
