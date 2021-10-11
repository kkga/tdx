package cmd

import (
	_ "embed"
	"os"

	"github.com/spf13/cobra"
)

//go:embed embed/help
var helpTxt string

var vdirPath string

var rootCmd = &cobra.Command{
	Use:          "tdx",
	Short:        "tdx -- todo manager for vdir (iCalendar) files",
	Long:         `long output`,
	Version:      "TODO",
	SilenceUsage: true,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// fmt.Println("some stuff")
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	BinaryName = os.Args[0]
	rootCmd.Flags().StringVarP(&vdirPath, "path", "p", "", "path to vdir folder")
	rootCmd.MarkFlagRequired("path")
	// rootCmd.Flags().BoolP("version", "v", false, "print version")
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
