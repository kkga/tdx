package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
)

func NewPurgeCmd() *PurgeCmd {
	c := &PurgeCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("purge", flag.ExitOnError),
		shortDesc: "Remove completed and cancelled todos",
	}}
	return c
}

type PurgeCmd struct {
	Cmd
}

func (c *PurgeCmd) Run() error {
	var removed int

	for _, items := range c.vdir {
		for _, item := range items {
			vtodo, err := item.Vtodo()
			if err != nil {
				return err
			}
			status, err := vtodo.Props.Text(ical.PropStatus)
			if err != nil {
				return err
			}

			s := vdir.ToDoStatus(status)
			if s == vdir.StatusCompleted || s == vdir.StatusCancelled {
				if err := os.Remove(item.Path); err != nil {
					return err
				}
				removed++
			}
		}
	}

	fmt.Printf("Removed: %d\n", removed)

	return nil
}
