package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Guo-Chenxu/ctm/logic"
	"github.com/spf13/cobra"
)

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Remove all managed environment variables",
	Long:  `Remove every environment variable managed by CTM from Claude Code's settings file.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := logic.LoadConfig()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error loading config: %v\n", err)
			fmt.Fprintf(cmd.ErrOrStderr(), "Config file should be at: %s\n", logic.GetConfigPath())
			os.Exit(1)
		}

		settings, err := logic.LoadClaudeSettings()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			os.Exit(1)
		}
		env, err := settings.Env()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			os.Exit(1)
		}

		removed := make([]string, 0)
		for _, key := range logic.GetAllEnvKeys(config) {
			if _, exists := env[key]; exists {
				delete(env, key)
				removed = append(removed, key)
			}
		}

		out := cmd.OutOrStdout()
		if len(removed) == 0 {
			fmt.Fprintln(out, "No CTM-managed environment variables found.")
			return
		}

		if err := settings.SetEnv(env); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			os.Exit(1)
		}
		if err := logic.SaveClaudeSettings(settings); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(out, "Removed %d variable(s) from %s\n", len(removed), logic.GetClaudeSettingsPath())
		fmt.Fprintf(out, "  %s\n", strings.Join(removed, ", "))
		fmt.Fprintln(out, "Restart Claude Code for the removal to take effect.")
	},
}
