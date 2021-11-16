package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func NewDocsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "docs",
		Short:                 "Generates cli docs",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Root().DisableAutoGenTag = true
			return doc.GenMarkdownTree(cmd.Root(), "doc")
		},
	}

	return cmd
}
