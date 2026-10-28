package cmd

import (
	"fmt"

	"github.com/Guo-Chenxu/ctm/logic"
	"github.com/spf13/cobra"
)

// providerIndexCompletion completes provider indices as "<index>\t<name>"
func providerIndexCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	config, err := logic.LoadConfig()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0, len(config.Providers))
	for i, provider := range config.Providers {
		completions = append(completions, fmt.Sprintf("%d\t%s", i, provider.Name))
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
