package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// envPathVar is the environment variable for setting default vdir path
const envPathVar = "TDX_PATH"

var (
	// vdirPath is the path to user's vdir folder
	vdirPath string

	rootCmd = &cobra.Command{
		Use:          "tdx",
		Short:        "tdx -- todo manager for vdir (iCalendar) files.",
		Long:         "tdx -- todo manager for vdir (iCalendar) files.",
		Version:      "dev",
		SilenceUsage: true,
		// TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				listCmd := NewListCmd()
				return listCmd.Execute()
			}
			return nil
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
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
	)
}

// TODO: run default command if no subcommands
// if len(args) == 0 {
// 	defaultCmd := NewListCmd()
// 	if err := defaultCmd.Init([]string{}); err != nil {
// 		return err
// 	}
// 	return defaultCmd.Run()
// }
