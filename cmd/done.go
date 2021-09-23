package cmd

import (
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
)

func NewDoneCmd() *DoneCmd {
	c := &DoneCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("done", flag.ExitOnError),
		alias:     []string{"do"},
		shortDesc: "Complete todo",
	}}
	return c
}

type DoneCmd struct {
	Cmd
}

func (c *DoneCmd) Run() error {
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

	// var collection vdir.Collection
	var item *vdir.Item

	for _, items := range collections {
		for _, i := range items {
			if i.Id == argID {
				// collection = *col
				item = i
			}
		}
	}

	if item == nil {
		return fmt.Errorf("Non-existing todo ID: %d", argID)
	}

	for _, comp := range item.Ical.Children {
		if comp.Name == ical.CompToDo {
			comp.Props.SetText(ical.PropStatus, string(vdir.StatusCompleted))
		}
	}

	if err := item.WriteFile(); err != nil {
		return err
	}

	for _, comp := range item.Ical.Children {
		if comp.Name == ical.CompToDo {
			t, err := item.Format()
			if err != nil {
				return err
			}
			fmt.Println(t)
		}
	}

	return nil
}
