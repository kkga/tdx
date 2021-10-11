package cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//go:embed embed/help
var helpTxt string

var vdirPath string

// var vd vdir.Vdir

var rootCmd = &cobra.Command{
	Use:              "tdx",
	Short:            "tdx -- todo manager for vdir (iCalendar) files",
	Long:             `long output`,
	Version:          "devdddd",
	SilenceUsage:     true,
	TraverseChildren: true,
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	vd = vdir.Vdir{}
	// 	if err := vd.Init(c.conf.Path); err != nil {
	// 		return err
	// 	}
	// 	return vd.Init(vdirPath)
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	defaultPath := fmt.Sprintf("%s/.local/share/calendars/migadu", home)

	rootCmd.Flags().StringVarP(&vdirPath, "path", "p", defaultPath, "path to vdir folder")
	rootCmd.MarkFlagRequired("path")

}

// func Root(args []string, version string) error {
// 	if len(args) > 0 {
// 		switch args[0] {
// 		case "-h", "--help":
// 			printHelp()
// 			os.Exit(0)
// 		case "-v", "--version":
// 			fmt.Printf("tdx version %s\n", version)
// 			os.Exit(0)
// 		}
// 	}

// 	cmds := []Runner{
// 		NewListCmd(),
// 		NewAddCmd(),
// 		NewDoneCmd(),
// 		NewShowCmd(),
// 		NewDeleteCmd(),
// 		NewEditCmd(),
// 		NewPurgeCmd(),
// 	}

// 	if len(args) == 0 {
// 		defaultCmd := NewListCmd()
// 		if err := defaultCmd.Init([]string{}); err != nil {
// 			return err
// 		}
// 		return defaultCmd.Run()
// 	}

// 	subcommand := os.Args[1]

// 	for _, cmd := range cmds {
// 		if cmd.Name() == subcommand || containsString(cmd.Alias(), subcommand) {
// 			if err := cmd.Init(os.Args[2:]); err != nil {
// 				return err
// 			}
// 			return cmd.Run()
// 		}
// 	}

// 	return fmt.Errorf("Unknown subcommand: %q. See 'tdx -h'", subcommand)
// }

// func containsString(s []string, e string) bool {
// 	for _, a := range s {
// 		if a == e {
// 			return true
// 		}
// 	}
// 	return false
// }

// func printHelp() {
// 	fmt.Print(helpTxt)
// }
