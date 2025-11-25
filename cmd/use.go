package cmd

import (
	"fmt"
	"strconv"

	"github.com/Guo-Chenxu/ctm/logic"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <index>",
	Short: "Use provider at specified index",
	Long:  `Switch to the provider at the specified index and export environment variables.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := logic.LoadConfig()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error loading config: %v\n", err)
			fmt.Fprintf(cmd.ErrOrStderr(), "Config file should be at: %s\n", logic.GetConfigPath())
			return
		}

		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: Invalid index '%s'. Must be a number.\n", args[0])
			return
		}

		if err := logic.ValidateIndex(config, index); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			return
		}

		// Update currentIndex and save config
		config.CurrentIndex = index
		if err := logic.SaveConfig(config); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error saving config: %v\n", err)
			return
		}

		provider := config.Providers[index]

		// Get all environment variable keys from all providers
		allKeys := logic.GetAllEnvKeys(config)

		// Output export commands for all keys
		// Current provider's keys: set to their values
		// Other keys: set to empty string
		for _, key := range allKeys {
			if value, exists := provider.Env[key]; exists {
				fmt.Printf("export %s=\"%s\"\n", key, value)
			} else {
				fmt.Printf("export %s=\"\"\n", key)
			}
		}
	},
}
