package cmd

import (
		flag "github.com/spf13/pflag"
	"fmt"
	"os"
	"strings"

	"github.com/kkga/tdx/vdir"
)

func NewDeleteCmd() *DeleteCmd {
	c := &DeleteCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("delete", flag.ExitOnError),
		alias:     []string{"del"},
		short:     "Delete todos",
		usageLine: "[options] <id>...",
	}}
	c.fs.BoolVar(&c.yes, "y", false, "do not ask for confimation")
	return c
}

type DeleteCmd struct {
	Cmd
	yes bool
}

func (c *DeleteCmd) Run() error {
	IDs, err := c.argsToIDs()
	if err != nil {
		return err
	}

	var toDelete []*vdir.Item

	for _, id := range IDs {
		item, err := c.vdir.ItemById(id)
		if err != nil {
			return err
		}
		toDelete = append(toDelete, item)
	}

	sb := strings.Builder{}

	for _, item := range toDelete {
		s, err := item.Format()
		if err != nil {
			return err
		}
		sb.WriteString(fmt.Sprintf("%s", s))
	}

	fmt.Print(sb.String())

	if !c.yes {
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
		for _, i := range toDelete {
			if err := os.Remove(i.Path); err != nil {
				return err
			}
		}
		fmt.Printf("Deleted: %d todos\n", len(toDelete))
	}

	return nil
}
