package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Guo-Chenxu/ctm/logic"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:               "use <index>",
	Short:             "Use provider at specified index",
	Long:              `Switch to the provider at the specified index by writing its environment variables to Claude Code's settings file.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: providerIndexCompletion,
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

		removed, err := applyProvider(config, index)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			os.Exit(1)
		}

		config.CurrentIndex = index
		if err := logic.SaveConfig(config); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error saving config: %v\n", err)
			os.Exit(1)
		}

		printUseResult(cmd, config.Providers[index], index, removed)
	},
}

// applyProvider drops every CTM-managed key from settings.json, then writes the provider's env
func applyProvider(config *logic.Config, index int) ([]string, error) {
	provider := config.Providers[index]

	settings, err := logic.LoadClaudeSettings()
	if err != nil {
		return nil, err
	}
	env, err := settings.Env()
	if err != nil {
		return nil, err
	}

	removed := make([]string, 0)
	for _, key := range logic.GetAllEnvKeys(config) {
		if _, kept := provider.Env[key]; kept {
			continue
		}
		if _, exists := env[key]; exists {
			delete(env, key)
			removed = append(removed, key)
		}
	}
	for key, value := range provider.Env {
		env[key] = value
	}

	if err := settings.SetEnv(env); err != nil {
		return nil, err
	}
	if err := logic.SaveClaudeSettings(settings); err != nil {
		return nil, err
	}
	return removed, nil
}

func printUseResult(cmd *cobra.Command, provider logic.Provider, index int, removed []string) {
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Switched to provider %d: %s\n", index, provider.Name)
	fmt.Fprintf(out, "Environment written to %s\n", logic.GetClaudeSettingsPath())

	keys := make([]string, 0, len(provider.Env))
	for key := range provider.Env {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(out, "  set %s\n", key)
	}

	if len(removed) > 0 {
		fmt.Fprintf(out, "  removed %s\n", strings.Join(removed, ", "))
		fmt.Fprintln(out, "Restart Claude Code for removed variables to take effect.")
	}
}
