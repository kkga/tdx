package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
	"github.com/spf13/cobra"
)

type doneOptions struct {
	toggle bool
}

func NewDoneCmd() *cobra.Command {
	opts := &doneOptions{}

	cmd := &cobra.Command{
		Use:     "done [options] <todo>...",
		Aliases: []string{"do"},
		Short:   "Complete todos",
		Long:    "Complete todos",
		Args:    cobra.MinimumNArgs(1),
		Example: heredoc.Doc(`
			$ tdx done 1
			$ tdx done 1 2 3`),
		RunE: func(cmd *cobra.Command, args []string) error {
			vd := make(vdir.Vdir)
			if err := vd.Init(vdirPath); err != nil {
				return err
			}

			IDs, err := stringsToInts(args)
			if err != nil {
				return err
			}

			var items []*vdir.Item
			for _, id := range IDs {
				item, err := vd.ItemById(id)
				if err != nil {
					return err
				}
				items = append(items, item)
			}

			return runDone(items, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.toggle, "toggle", "t", false, "toggle completed state")

	return cmd
}

func runDone(items []*vdir.Item, opts *doneOptions) error {
	sb := strings.Builder{}

	for _, item := range items {
		vtodo, err := item.Vtodo()
		if err != nil {
			return err
		}

		if opts.toggle {
			s, err := vtodo.Props.Text(ical.PropStatus)
			if err != nil {
				return err
			}
			switch vdir.ToDoStatus(s) {
			case vdir.StatusCompleted:
				vtodo.Props.SetText(ical.PropStatus, string(vdir.StatusNeedsAction))
			case vdir.StatusNeedsAction, vdir.StatusCancelled:
				vtodo.Props.SetText(ical.PropStatus, string(vdir.StatusCompleted))
			default:
				vtodo.Props.SetText(ical.PropStatus, string(vdir.StatusCompleted))
			}
		} else {
			vtodo.Props.SetText(ical.PropStatus, string(vdir.StatusCompleted))
		}

		if err := item.WriteFile(); err != nil {
			return err
		}

		s, err := item.Format()
		if err != nil {
			return err
		}
		sb.WriteString(fmt.Sprintf("%s", s))
	}

	fmt.Print(sb.String())
	return nil
}
