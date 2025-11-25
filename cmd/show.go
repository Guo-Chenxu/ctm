package cmd

import (
	"fmt"
	"strconv"

	"github.com/Guo-Chenxu/ctm/logic"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [index]",
	Short: "Show provider details",
	Long:  `Display detailed information about a provider. If no index is specified, shows the current provider.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := logic.LoadConfig()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error loading config: %v\n", err)
			fmt.Fprintf(cmd.ErrOrStderr(), "Config file should be at: %s\n", logic.GetConfigPath())
			return
		}

		var index int
		if len(args) == 0 {
			// Show current provider
			index = config.CurrentIndex
		} else {
			// Show specified provider
			var err error
			index, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: Invalid index '%s'. Must be a number.\n", args[0])
				return
			}
		}

		if err := logic.ValidateIndex(config, index); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			return
		}

		provider := config.Providers[index]
		current := ""
		if index == config.CurrentIndex {
			current = " (current)"
		}

		fmt.Printf("Provider %d%s:\n", index, current)
		fmt.Printf("  Name: %s\n", provider.Name)
		fmt.Printf("  Description: %s\n", provider.Description)
		fmt.Println("  Environment variables:")
		for key, value := range provider.Env {
			fmt.Printf("    %s: %s\n", key, value)
		}
	},
}
