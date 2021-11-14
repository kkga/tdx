package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/kkga/tdx/vdir"
	"github.com/spf13/cobra"
)

type showOptions struct {
	raw bool
}

func NewShowCmd() *cobra.Command {
	opts := &showOptions{}

	cmd := &cobra.Command{
		Use:   "show [options] <todo>...",
		Short: "Show todos",
		Long:  "Show detailed info about todos",
		Args:  cobra.MinimumNArgs(1),
		Example: heredoc.Doc(`
			$ tdx show 1
			$ tdx show 1 2 3`),
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

			return runShow(items, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.raw, "raw", "r", false, "raw output")

	return cmd
}

func runShow(items []*vdir.Item, opts *showOptions) error {
	sb := strings.Builder{}

	for i, item := range items {
		var s string
		if opts.raw {
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

		if i < len(items)-1 {
			sb.WriteString("\n")
		}
	}

	fmt.Print(sb.String())

	return nil
}
