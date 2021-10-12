//go:build skip

package cmd

import (
	"fmt"
	"strings"

	"github.com/kkga/tdx/vdir"
	flag "github.com/spf13/pflag"
)

func NewShowCmd() *ShowCmd {
	c := &ShowCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("show", flag.ExitOnError),
		name:      "show",
		short:     "Show detailed info for todo",
		usageLine: "[options] <id>...",
	}}
	c.fs.BoolVarP(&c.raw, "raw", "r", false, "raw output")
	c.fs.StringVarP(&c.test, "test", "t", "", "test flag")
	return c
}

type ShowCmd struct {
	Cmd
	raw  bool
	test string
}

func (c *ShowCmd) Run() error {
	IDs, err := c.argsToIDs()
	if err != nil {
		return err
	}

	fmt.Println(c.fs.Args())

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
