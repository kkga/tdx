package cmd

import (
	"bytes"
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
	c.fs.StringVar(&c.due, "d", "", "due date")
	return c
}

type AddCmd struct {
	Cmd
	priority string
	due      string
}

func (c *AddCmd) Run() error {
	fmt.Println(c.fs.Args())

	// parse flags and create vtodo
	// encode into a new cal
	// write cal to file
	// print result

	// 	todo := &todo.ToDo{}
	// 	todoBuf, _ := todo.Encode()

	// 	f, _ := os.Create(fmt.Sprintf("%s/ctdo-testing-%d.ics", calDir, time.Now().UnixNano()))
	// 	defer f.Close()
	// 	w := bufio.NewWriter(f)
	// 	_, _ = w.Write(todoBuf.Bytes())
	// 	w.Flush()
	return nil
}

// generateUID returns a random string containing timestamp and hostname
func generateUID() string {
	sb := strings.Builder{}

	randStr := func(n int) string {
		var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyz")

		s := make([]rune, n)
		for i := range s {
			s[i] = letters[rand.Intn(len(letters))]
		}
		return string(s)
	}

	sb.WriteString(fmt.Sprint(time.Now().UnixNano()))
	sb.WriteString(fmt.Sprintf("-%s", randStr(8)))
	if hostname, _ := os.Hostname(); hostname != "" {
		sb.WriteString(fmt.Sprintf("@%s", hostname))
	}

	return sb.String()
}

// encode adds vtodo into a new Calendar and returns a buffer ready for writing
func encode(vtodo *ical.Component) (bytes.Buffer, error) {
	vtodo.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())

	if vtodo.Props.Get(ical.PropUID).Value == "" {
		vtodo.Props.SetText(ical.PropUID, generateUID())
	}

	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN")
	cal.Children = append(cal.Children, vtodo)

	var buf bytes.Buffer
	if err := ical.NewEncoder(&buf).Encode(cal); err != nil {
		return buf, err
	}
	return buf, nil
}
