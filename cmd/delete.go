package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/kkga/tdx/vdir"
	"github.com/spf13/cobra"
)

type deleteOptions struct {
	yes bool
}

func NewDeleteCmd() *cobra.Command {
	opts := &deleteOptions{}

	cmd := &cobra.Command{
		Use:     "delete <id>...",
		Aliases: []string{"del"},
		Short:   "Delete todos",
		Long:    "Permanently delete todos.",
		Args:    cobra.MinimumNArgs(1),
		Example: heredoc.Doc(`
			$ tdx delete 1
			$ tdx delete 1 2 3`),
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

			return runDelete(items, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.yes, "yes", "y", false, "do not ask for confirmation")

	return cmd
}

func runDelete(items []*vdir.Item, opts *deleteOptions) error {
	sb := strings.Builder{}

	for _, item := range items {
		s, err := item.Format()
		if err != nil {
			return err
		}
		sb.WriteString(fmt.Sprintf("%s", s))
	}

	fmt.Print(sb.String())

	if !opts.yes {
		ok := promptConfirm("Delete listed todos?", false)
		if ok {
			for _, i := range items {
				if err := os.Remove(i.Path); err != nil {
					return err
				}
			}
			fmt.Printf("Deleted: %d todos\n", len(items))
		}
	} else {
		for _, i := range items {
			if err := os.Remove(i.Path); err != nil {
				return err
			}
		}
		fmt.Printf("Deleted: %d todos\n", len(items))
	}

	return nil
}
