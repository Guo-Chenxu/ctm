package cmd

import (
	"fmt"

	"github.com/Guo-Chenxu/ctm/logic"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show config file path",
	Long:  `Display the path to the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := logic.GetConfigPath()
		fmt.Printf("Config file path: %s\n", configPath)
	},
}
