package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
	"github.com/spf13/cobra"
)

func NewPurgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge",
		Short: "Delete done todos",
		Long:  "Permanently delete completed and cancelled todos.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			vd := make(vdir.Vdir)
			if err := vd.Init(vdirPath); err != nil {
				return err
			}

			var toDelete []*vdir.Item
			for _, items := range vd {
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

			if len(toDelete) == 0 {
				return errors.New("No items to purge")
			}

			return runPurge(toDelete)
		},
	}

	return cmd
}

func runPurge(items []*vdir.Item) error {
	sb := strings.Builder{}

	for _, item := range items {
		s, err := item.Format()
		if err != nil {
			return err
		}
		sb.WriteString(fmt.Sprintf("%s", s))
	}

	fmt.Print(sb.String())

	ok := promptConfirm("Delete listed todos?", false)
	if ok {
		for _, i := range items {
			if err := os.Remove(i.Path); err != nil {
				return err
			}
		}
		fmt.Printf("Deleted: %d todos\n", len(items))
	}

	return nil
}
