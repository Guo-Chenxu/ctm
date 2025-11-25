package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ctm",
	Short: "CTM - Configuration Token Manager",
	Long: `CTM is a command-line tool for managing configuration tokens.
It helps you manage multiple provider configurations and their environment variables.`,
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(initCmd)
}
