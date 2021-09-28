package cmd

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed embed/help
var helpTxt string

func Root(args []string, version string) error {
	if len(args) > 0 {
		switch args[0] {
		case "-h", "--help":
			printHelp()
			os.Exit(0)
		case "-v", "--version":
			fmt.Printf("tdx version %s\n", version)
			os.Exit(0)
		}
	}

	cmds := []Runner{
		NewListCmd(),
		NewAddCmd(),
		NewDoneCmd(),
		NewShowCmd(),
		NewDeleteCmd(),
		NewEditCmd(),
		NewPurgeCmd(),
	}

	if len(args) == 0 {
		defaultCmd := NewListCmd()
		if err := defaultCmd.Init([]string{}); err != nil {
			return err
		}
		return defaultCmd.Run()
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand || containsString(cmd.Alias(), subcommand) {
			if err := cmd.Init(os.Args[2:]); err != nil {
				return err
			}
			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %q", subcommand)
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func printHelp() {
	fmt.Print(helpTxt)
}
