package cmd

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/emersion/go-ical"
	"github.com/fatih/color"
	"github.com/kkga/tdx/vdir"
	"github.com/kkga/tdx/vtodo"
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
func FormatToDo(comp ical.Component) (string, error) {
	sb := strings.Builder{}
	colorStatusDone := color.New(color.FgGreen).SprintFunc()
	colorStatusUndone := color.New(color.FgRed, color.Bold).SprintFunc()
	// colorPrio := color.New(color.FgRed, color.Bold).SprintFunc()
	// colorDate := color.New(color.FgYellow).SprintFunc()
	// colorDesc := color.New(color.Faint).SprintFunc()

	if comp.Name != ical.CompToDo {
		return "", fmt.Errorf("Not VTODO component: %v", comp)
	}

	var (
		status      string
		summary     string
		description string
		prio        string
		// due         string
	)

	for name, prop := range comp.Props {
		p := prop[0]

		switch name {

		case ical.PropStatus:
			if p.Value == vtodo.StatusCompleted {
				status = colorStatusDone("[x]")
			} else {
				status = colorStatusUndone("[ ]")
			}

		case ical.PropSummary:
			summary = p.Value

		case ical.PropDescription:
			description = p.Value

		case ical.PropPriority:
			v, err := strconv.Atoi(p.Value)
			if err != nil {
				return "", err
			}
			switch {
			case v == vtodo.PriorityHigh:
				prio = "!!!"
			case v > vtodo.PriorityHigh && v <= vtodo.PriorityMedium:
				prio = "!!"
			case v > vtodo.PriorityMedium:
				prio = "!"
			}

		}
	}

	sb.WriteString(status)
	sb.WriteString(prio)
	sb.WriteString(summary)
	sb.WriteString(description)

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

	return sb.String(), nil
}
