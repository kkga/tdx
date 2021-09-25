package cmd

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
)

func NewDoneCmd() *DoneCmd {
	c := &DoneCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("done", flag.ExitOnError),
		alias:     []string{"do"},
		shortDesc: "Complete todos",
		usageLine: "[options] <id>...",
	}}
	c.fs.BoolVar(&c.toggle, "t", false, "toggle complete state")
	return c
}

type DoneCmd struct {
	Cmd
	toggle bool
}

func (c *DoneCmd) Run() error {
	if len(c.fs.Args()) == 0 {
		return errors.New("Specify one or multiple IDs")
	}

	IDs := make([]int, len(c.fs.Args()))

	for i, s := range c.fs.Args() {
		IDs[i], _ = strconv.Atoi(s)
	}

	var toComplete []*vdir.Item

	containsInt := func(ii []int, i int) bool {
		for _, v := range ii {
			if v == i {
				return true
			}
		}
		return false
	}

	// TODO: rewrite with ItemById
	for _, items := range c.vdir {
		for _, item := range items {
			if containsInt(IDs, item.Id) {
				toComplete = append(toComplete, item)
			}
		}
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
