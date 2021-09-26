package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kkga/tdx/vdir"
)

func NewDeleteCmd() *DeleteCmd {
	c := &DeleteCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("delete", flag.ExitOnError),
		alias:     []string{"del"},
		shortDesc: "Delete todos",
		usageLine: "<id>...",
	}}
	return c
}

type DeleteCmd struct {
	Cmd
}

func (c *DeleteCmd) Run() error {
	if len(c.fs.Args()) == 0 {
		return errors.New("Specify one or multiple IDs")
	}

	IDs := make([]int, len(c.fs.Args()))

	for i, s := range c.fs.Args() {
		IDs[i], _ = strconv.Atoi(s)
	}

	var toDelete []*vdir.Item

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
				toDelete = append(toDelete, item)
			}
		}
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

	ok := promptConfirm("Delete listed todos?", false)
	if ok {
		for _, i := range toDelete {
			if err := os.Remove(i.Path); err != nil {
				return err
			}
		}
		fmt.Printf("Deleted: %d todos\n", len(toDelete))
	}

	return nil
}