package cmd

import (
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
)

func NewUndoCmd() *UndoCmd {
	c := &UndoCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("undo", flag.ExitOnError),
		alias:     []string{""},
		shortDesc: "Mark todo as not done",
	}}
	return c
}

type UndoCmd struct {
	Cmd
}

func (c *UndoCmd) Run() error {
	if len(c.fs.Args()) == 0 {
		return errors.New("Specify a todo ID.")
	}

	argID, err := strconv.Atoi(c.fs.Arg(0))
	if err != nil {
		return err
	}

	collections, err := c.root.Collections()
	if err != nil {
		return err
	}

	var item *vdir.Item
	var vtodo *ical.Component

	for _, items := range collections {
		for _, i := range items {
			if i.Id == argID {
				item = i
			}
		}
	}
	if item == nil {
		return fmt.Errorf("Non-existing todo ID: %d", argID)
	}

	for _, comp := range item.Ical.Children {
		if comp.Name == ical.CompToDo {
			vtodo = comp
		}
	}

	vtodo.Props.SetText(ical.PropStatus, string(vdir.StatusNeedsAction))

	if err := item.WriteFile(); err != nil {
		return err
	}

	f, err := item.Format()
	fmt.Println(f)

	return nil
}
