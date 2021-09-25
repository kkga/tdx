package cmd

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
	"github.com/tj/go-naturaldate"
)

func NewEditCmd() *EditCmd {
	c := &EditCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("edit", flag.ExitOnError),
		alias:     []string{"e"},
		shortDesc: "Edit todo",
		usageLine: "[options] <id>",
	}}
	return c
}

type EditCmd struct {
	Cmd
}

const (
	layoutDateTime = "2 Jan 2006 15:04"
	layoutDate     = "2 Jan 2006"
)

func (c *EditCmd) Run() error {
	if len(c.fs.Args()) == 0 {
		return errors.New("Specify todo ID")
	}

	id, err := strconv.Atoi(c.fs.Arg(0))
	if err != nil {
		return fmt.Errorf("Invalid todo ID: %s", c.fs.Arg(0))
	}

	item, err := c.vdir.ItemById(id)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp("", "tdx")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	editor := os.Getenv("VISUAL")
	editorBin, err := exec.LookPath(editor)
	if err != nil {
		return errors.New("Set the VISUAL env variable to edit todos")
	}

	vtodo, err := item.Vtodo()
	if err != nil {
		return err
	}

	template := generateTemplate(*vtodo)
	_, err = tmp.Write([]byte(template))
	if err != nil {
		return err
	}

	cmd := exec.Command(editorBin, tmp.Name())
	cmd.SysProcAttr = &syscall.SysProcAttr{Foreground: true}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Run()

	newProps, err := parseTemplate(tmp)
	if err != nil {
		return err
	}

	for p, newVal := range newProps {
		curP := vtodo.Props.Get(p)
		newP := ical.NewProp(p)

		if curP == nil && newVal == "" {
			continue
		}

		if curP != nil && curP.Value == "" && newVal == "" {
			continue
		}

		if curP != nil && curP.Value != "" && newVal == "" {
			vtodo.Props.Del(p)
			continue
		}

		switch p {
		case ical.PropDue:
			if t, _ := time.Parse(layoutDateTime, newVal); !t.IsZero() {
				newP.SetDateTime(t)
			} else if t, _ := time.Parse(layoutDate, newVal); !t.IsZero() {
				newP.SetDateTime(t)
			} else {
				now := time.Now()
				due, err := naturaldate.Parse(newVal, now, naturaldate.WithDirection(naturaldate.Future))
				if err != nil {
					return err
				}
				if due != now {
					newP.SetDateTime(due)
				}
			}
		case ical.PropPriority:
			prioMap := map[string]vdir.ToDoPriority{
				"high":   vdir.PriorityHigh,
				"medium": vdir.PriorityMedium,
				"low":    vdir.PriorityLow,
			}
			prio := prioMap[newVal]
			newP.Value = fmt.Sprint(prio)
		default:
			newP.Value = newVal
		}

		vtodo.Props.Set(newP)
	}

	if err := item.WriteFile(); err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	f, err := item.Format(vdir.FormatDescription, vdir.FormatMultiline)
	if err != nil {
		return err
	}
	fmt.Print(f)

	return nil
}

func parseTemplate(f *os.File) (map[string]string, error) {
	props := make(map[string]string)

	f.Seek(0, 0)

	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() == "" || strings.HasPrefix(s.Text(), "-") {
			break
		}
		s := strings.SplitN(s.Text(), ":", 2)
		k := strings.Trim(s[0], " ")
		v := strings.Trim(s[1], " ")
		props[k] = v
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return props, nil
}

func generateTemplate(vtodo ical.Component) string {
	sb := strings.Builder{}

	props := map[string]string{
		ical.PropSummary:     "",
		ical.PropDescription: "",
		ical.PropStatus:      "",
		ical.PropPriority:    "",
		ical.PropDue:         "",
		ical.PropLocation:    "",
	}

	prioMap := map[vdir.ToDoPriority]string{
		vdir.PriorityHigh:   "high",
		vdir.PriorityMedium: "medium",
		vdir.PriorityLow:    "low",
	}

	for name, prop := range vtodo.Props {
		p := prop[0]
		if _, exists := props[name]; exists {
			var val string
			switch name {
			case ical.PropDue:
				date, _ := p.DateTime(time.Local)
				if !date.IsZero() {
					val = date.Format(layoutDateTime)
				}
			case ical.PropPriority:
				curPrio, _ := strconv.Atoi(p.Value)
				val = prioMap[vdir.ToDoPriority(curPrio)]
			default:
				val = p.Value
			}
			props[p.Name] = val
		}
	}

	order := []string{
		ical.PropSummary,
		ical.PropDescription,
		ical.PropStatus,
		ical.PropPriority,
		ical.PropDue,
		ical.PropLocation,
	}

	for _, p := range order {
		sb.WriteString(fmt.Sprintf("%s: %s\n", p, props[p]))
	}

	sb.WriteString(templateHelp)

	return sb.String()
}

const templateHelp = `
--------------------

DUE accepts following formats:
- "2 Jan 2006 15:04"
- "2 Jan 2006"
- natural date, same as <add> command, see "tdx add -h"

PRIORITY can be:
- low
- medium
- high
- [empty]

STATUS can be:
- NEEDS-ACTION
- COMPLETED
- CANCELLED
`
