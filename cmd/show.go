package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/kkga/tdx/vdir"
)

func NewShowCmd() *ShowCmd {
	c := &ShowCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("show", flag.ExitOnError),
		shortDesc: "Show detailed info for todo",
		usageLine: "[options] <id>...",
	}}
	c.fs.BoolVar(&c.raw, "r", false, "raw output")
	return c
}

type ShowCmd struct {
	Cmd
	raw bool
}

func (c *ShowCmd) Run() error {
	IDs, err := c.argsToIDs()
	if err != nil {
		return err
	}

	var toShow []*vdir.Item

	for _, id := range IDs {
		item, err := c.vdir.ItemById(id)
		if err != nil {
			return err
		}
		toShow = append(toShow, item)
	}

	sb := strings.Builder{}

	for i, item := range toShow {
		var s string
		if c.raw {
			ff, err := item.FormatFull(vdir.FormatFullRaw)
			if err != nil {
				return err
			}
			s = ff
		} else {
			ff, err := item.FormatFull()
			if err != nil {
				return err
			}
			s = ff
		}
		sb.WriteString(s)
		if i < len(toShow)-1 {
			sb.WriteString("\n")
		}
	}

	fmt.Print(sb.String())

	return nil
}
