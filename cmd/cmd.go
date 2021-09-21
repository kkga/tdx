package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/emersion/go-ical"
	"github.com/kkga/ctodo/vdir"
)

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
	Alias() []string
}

type Cmd struct {
	fs        *flag.FlagSet
	alias     []string
	shortDesc string
	usageLine string
	cals      []ical.Calendar
	root      *vdir.VdirRoot
}

var dir = "/home/kkga/.local/share/calendars/"

func (c *Cmd) Run() error      { return nil }
func (c *Cmd) Name() string    { return c.fs.Name() }
func (c *Cmd) Alias() []string { return c.alias }

func (c *Cmd) Init(args []string) error {
	c.fs.Usage = c.usage

	if err := c.fs.Parse(args); err != nil {
		return err
	}

	c.root = vdir.NewVdirRoot(dir)

	return nil
}

func (c *Cmd) usage() {
	fmt.Println(c.shortDesc)
	fmt.Println()

	fmt.Println("USAGE")
	fmt.Printf("  ctodo %s %s\n\n", c.fs.Name(), c.usageLine)

	if strings.Contains(c.usageLine, "[options]") {
		fmt.Println("OPTIONS")
		c.fs.PrintDefaults()
	}
}

// TODO: refactor this for using the ical.Component directly
func FormatToDo(vtodo ical.Component) string {
	sb := strings.Builder{}
	// colorPrio := color.New(color.FgRed, color.Bold).SprintFunc()
	// colorDate := color.New(color.FgYellow).SprintFunc()
	// colorDesc := color.New(color.Faint).SprintFunc()

	for p, v := range vtodo.Props {
		switch p {
		case ical.PropStatus:
			if s := v[0].Value; s != "" {
				if s == vdir.StatusCompleted {
					sb.WriteString("[x] ")
				} else {
					sb.WriteString("[ ] ")
				}
			}
		case ical.PropSummary:
			if s := v[0].Value; s != "" {
				sb.WriteString(s)
				sb.WriteString("\n")
			}
		}
	}

	// summary := vtodo.Props.Get(ical.PropSummary)
	// description := vtodo.Props.Get(ical.PropDescription)
	// status := vtodo.Props.Get(ical.PropStatus)
	// prio := vtodo.Props.Get(ical.PropPriority)
	// due := vtodo.Props.Get(ical.PropDue)

	// fmt.Println(prio.Value)

	// if prio.Value != "" {
	// 	p, _ := strconv.Atoi(prio.Value)
	// 	var ps string
	// 	switch {
	// 	case p == vdir.PriorityHigh:
	// 		ps = "!!!"
	// 	case p > vdir.PriorityHigh && p <= vdir.PriorityMedium:
	// 		ps = "!!"
	// 	case p > vdir.PriorityMedium:
	// 		ps = "!"
	// 	}
	// 	sb.WriteString(fmt.Sprintf(" %s", colorPrio(ps)))
	// }

	// if d, err := due.DateTime(time.Local); err != nil {
	// 	date := d.Local().Format("Jan-06")
	// 	sb.WriteString(fmt.Sprintf(" %s", colorDate(date)))
	// }

	// if s, err := summary.Text(); err != nil {
	// 	sb.WriteString(" ")
	// 	sb.WriteString(s)
	// }

	// if d, err := description.Text(); err != nil {
	// 	sb.WriteString(colorDesc("\n    ↳ "))
	// 	sb.WriteString(colorDesc(d))
	// }

	return sb.String()
}
