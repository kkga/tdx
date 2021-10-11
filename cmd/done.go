package cmd

import (
	"fmt"
	"strings"

	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
	flag "github.com/spf13/pflag"
)

func NewDoneCmd() *DoneCmd {
	c := &DoneCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("done", flag.ExitOnError),
		alias:     []string{"do"},
		short:     "Complete todos",
		usageLine: "[options] <id>...",
	}}
	c.fs.BoolVar(&c.toggle, "t", false, "toggle completed state")
	return c
}

type DoneCmd struct {
	Cmd
	toggle bool
}

func (c *DoneCmd) Run() error {
	IDs, err := c.argsToIDs()
	if err != nil {
		return err
	}

	var toComplete []*vdir.Item

	for _, id := range IDs {
		item, err := c.vdir.ItemById(id)
		if err != nil {
			return err
		}
		toComplete = append(toComplete, item)
	}

	sb := strings.Builder{}

	for _, item := range toComplete {
		vtodo, err := item.Vtodo()
		if err != nil {
			return err
		}

		if c.toggle {
			s, err := vtodo.Props.Text(ical.PropStatus)
			if err != nil {
				return err
			}
			switch vdir.ToDoStatus(s) {
			case vdir.StatusCompleted:
				vtodo.Props.SetText(ical.PropStatus, string(vdir.StatusNeedsAction))
			case vdir.StatusNeedsAction, vdir.StatusCancelled:
				vtodo.Props.SetText(ical.PropStatus, string(vdir.StatusCompleted))
			default:
				vtodo.Props.SetText(ical.PropStatus, string(vdir.StatusCompleted))
			}
		} else {
			vtodo.Props.SetText(ical.PropStatus, string(vdir.StatusCompleted))
		}

		if err := item.WriteFile(); err != nil {
			return err
		}

		s, err := item.Format()
		if err != nil {
			return err
		}
		sb.WriteString(fmt.Sprintf("%s", s))
	}

	fmt.Print(sb.String())

	return nil
}
