package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	vdirPath   string
	version    = "dev"
	defaultCmd = NewListCmd()

	rootCmd = &cobra.Command{
		Use:          "tdx",
		Short:        "tdx -- todo manager for vdir (iCalendar) files.",
		Version:      version,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return defaultCmd.Execute()
			}
			return nil
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	const envPathVar = "TDX_PATH"
	var defaultPath string

	if defaultPath = os.Getenv(envPathVar); defaultPath == "" {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		defaultPath = fmt.Sprintf("%s/.local/share/calendars", home)
	}

	rootCmd.PersistentFlags().StringVarP(&vdirPath, "path", "p", defaultPath, "path to vdir folder")
	rootCmd.MarkFlagRequired("path")

	cobra.EnableCommandSorting = false
	rootCmd.AddCommand(
		NewAddCmd(),
		NewListCmd(),
		NewDoneCmd(),
		NewEditCmd(),
		NewShowCmd(),
		NewDeleteCmd(),
		NewPurgeCmd(),
		NewDocsCmd(),
	)
}
