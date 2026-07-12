package cmd

import (
	"fmt"
	"os"
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
			os.Exit(1)
		}

		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: Invalid index '%s'. Must be a number.\n", args[0])
			os.Exit(1)
		}

		if err := logic.ValidateIndex(config, index); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			os.Exit(1)
		}

		provider := config.Providers[index]
		commands, err := buildUseCommands(config, provider)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			os.Exit(1)
		}

		// Update currentIndex and save config only after all shell commands validate.
		config.CurrentIndex = index
		if err := logic.SaveConfig(config); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error saving config: %v\n", err)
			os.Exit(1)
		}

		for _, command := range commands {
			fmt.Fprintln(cmd.OutOrStdout(), command)
		}
	},
}

func buildUseCommands(config *logic.Config, provider logic.Provider) ([]string, error) {
	allKeys := logic.GetAllEnvKeys(config)
	commands := make([]string, 0, len(allKeys))

	for _, key := range allKeys {
		if !isValidShellIdentifier(key) {
			return nil, fmt.Errorf("invalid environment variable name %q", key)
		}

		if value, exists := provider.Env[key]; exists {
			commands = append(commands, fmt.Sprintf("export %s=%s", key, shellQuote(value)))
		} else {
			commands = append(commands, fmt.Sprintf("unset %s", key))
		}
	}

	return commands, nil
}
