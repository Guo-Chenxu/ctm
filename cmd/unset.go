package cmd

import (
	"fmt"
	"os"

	"github.com/Guo-Chenxu/ctm/logic"
	"github.com/spf13/cobra"
)

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset all managed environment variables",
	Long:  `Output shell commands to unset all environment variables managed by CTM.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := logic.LoadConfig()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error loading config: %v\n", err)
			fmt.Fprintf(cmd.ErrOrStderr(), "Config file should be at: %s\n", logic.GetConfigPath())
			os.Exit(1)
		}

		commands, err := buildUnsetCommands(config)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			os.Exit(1)
		}

		for _, command := range commands {
			fmt.Fprintln(cmd.OutOrStdout(), command)
		}
	},
}

func buildUnsetCommands(config *logic.Config) ([]string, error) {
	keys := logic.GetAllEnvKeys(config)
	commands := make([]string, 0, len(keys))
	for _, key := range keys {
		if !isValidShellIdentifier(key) {
			return nil, fmt.Errorf("invalid environment variable name %q", key)
		}
		commands = append(commands, fmt.Sprintf("unset %s", key))
	}
	return commands, nil
}
