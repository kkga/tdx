package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
	var toDelete []*vdir.Item

	for _, items := range c.vdir {
		for _, item := range items {
			vtodo, err := item.Vtodo()
			if err != nil {
				return err
			}
			s, err := vtodo.Props.Text(ical.PropStatus)
			if err != nil {
				return err
			}

			switch vdir.ToDoStatus(s) {
			case vdir.StatusCancelled, vdir.StatusCompleted:
				toDelete = append(toDelete, item)
			}
		}
	}

	if len(toDelete) > 0 {
		sb := strings.Builder{}

		for _, item := range toDelete {
			s, err := item.Format()
			if err != nil {
				return err
			}
			sb.WriteString(fmt.Sprintf("%s", s))
		}

		fmt.Print(sb.String())

		ok := promptConfirm("Delete listed todos?", false)
		if ok {
			for _, i := range toDelete {
				if err := os.Remove(i.Path); err != nil {
					return err
				}
			}
			fmt.Printf("Deleted: %d todos\n", len(toDelete))
		}
	} else {
		fmt.Println("No todos to purge")
	}

	return nil
}
