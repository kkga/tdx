package cmd

import (
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

func (c *EditCmd) Run() error {
	if len(c.fs.Args()) == 0 {
		return errors.New("Specify todo ID")
	}

	id, err := strconv.Atoi(c.fs.Arg(0))
	if err != nil {
		return err
	}

	var item *vdir.Item

	for _, items := range c.allCollections {
		for _, i := range items {
			if i.Id == id {
				item = i
			}
		}
	}

	editor := os.Getenv("VISUAL")

	tmp, err := os.CreateTemp("", "tdx")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	binary, err := exec.LookPath(editor)
	if err != nil {
		return err
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

	args := []string{binary, tmp.Name()}
	env := os.Environ()

	err = syscall.Exec(binary, args, env)
	if err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	return nil
}

func parseTemplate(content string) map[string]string {
	props := make(map[string]string)

	return props
}

func generateTemplate(vtodo ical.Component) string {
	sb := strings.Builder{}

	for name, prop := range vtodo.Props {
		p := prop[0]
		switch name {
		case ical.PropSummary:
			sb.WriteString(fmt.Sprintf("SUMMARY: %s\n", p.Value))
		case ical.PropDescription:
			sb.WriteString(fmt.Sprintf("DESCRIPTION: %s\n", p.Value))
		case ical.PropLocation:
			sb.WriteString(fmt.Sprintf("LOCATION: %s\n", p.Value))
		case ical.PropStatus:
			sb.WriteString(fmt.Sprintf("STATUS: %s\n", p.Value))
		case ical.PropDue:
			date, _ := p.DateTime(time.Local)
			if !date.IsZero() {
				sb.WriteString(fmt.Sprintf("DUE: %s\n", date.Format("2 Jan 2006 15:04")))
			}
		case ical.PropDateTimeStart:
			date, _ := p.DateTime(time.Local)
			if !date.IsZero() {
				sb.WriteString(fmt.Sprintf("START: %s\n", date.Format("2 Jan 2006 15:04")))
			}
		}
	}

	return sb.String()
}
