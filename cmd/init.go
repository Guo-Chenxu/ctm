package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [shell]",
	Short: "Generate shell initialization script",
	Long:  `Generate shell initialization script to enable automatic environment variable setting.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := "bash"
		if len(args) > 0 {
			shell = args[0]
		}

		switch shell {
		case "bash", "zsh":
			fmt.Println("# Add this to your ~/.bashrc or ~/.zshrc:")
			fmt.Println("ctm-use() {")
			fmt.Println("    eval $(ctm use \"$@\")")
			fmt.Println("}")
			fmt.Println()
			fmt.Println("# Then you can use: ctm-use <index>")
		case "fish":
			fmt.Println("# Add this to your ~/.config/fish/config.fish:")
			fmt.Println("function ctm-use")
			fmt.Println("    eval (ctm use $argv)")
			fmt.Println("end")
			fmt.Println()
			fmt.Println("# Then you can use: ctm-use <index>")
		default:
			fmt.Fprintf(cmd.ErrOrStderr(), "Unsupported shell: %s\n", shell)
			fmt.Fprintf(cmd.ErrOrStderr(), "Supported shells: bash, zsh, fish\n")
		}
	},
}
