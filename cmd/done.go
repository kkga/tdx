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
		usageLine: "<id>...",
	}}
	return c
}

type DoneCmd struct {
	Cmd
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

	for _, items := range c.allCollections {
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
		vtodo.Props.SetText(ical.PropStatus, string(vdir.StatusCompleted))

		if err := item.WriteFile(); err != nil {
			return err
		}

		t, err := item.Format()
		if err != nil {
			return err
		}
		sb.WriteString(fmt.Sprintf("%s\n", t))
	}

	fmt.Println(sb.String())

	return nil
}
