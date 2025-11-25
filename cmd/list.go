package cmd

import (
	"fmt"

	"github.com/Guo-Chenxu/ctm/logic"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all provider names",
	Long:  `Display a list of all configured providers with their indices.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := logic.LoadConfig()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error loading config: %v\n", err)
			fmt.Fprintf(cmd.ErrOrStderr(), "Config file should be at: %s\n", logic.GetConfigPath())
			return
		}

		if len(config.Providers) == 0 {
			fmt.Println("No providers configured.")
			return
		}

		fmt.Println("Available providers:")
		for i, provider := range config.Providers {
			if i == config.CurrentIndex {
				// Current provider: add * prefix and color green
				fmt.Printf("\033[32m* %d: %s\033[0m\n", i, provider.Name)
			} else {
				fmt.Printf("  %d: %s\n", i, provider.Name)
			}
		}
	},
}
