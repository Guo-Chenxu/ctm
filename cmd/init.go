package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	initBlockStart = "# >>> ctm init >>>"
	initBlockEnd   = "# <<< ctm init <<<"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Write shell initialization script",
	Long:  `Write shell initialization script to enable automatic environment variable setting.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInit(cmd); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runInit(cmd *cobra.Command) error {
	shell, err := detectShell(os.Getenv("SHELL"))
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	rcPath, err := shellRCPath(shell, homeDir)
	if err != nil {
		return err
	}

	ctmDir, err := ctmExecutableDir()
	if err != nil {
		return err
	}

	currentContent, err := os.ReadFile(rcPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read %s: %w", rcPath, err)
	}

	block := buildInitBlock(ctmDir)
	updatedContent := upsertInitBlock(string(currentContent), block)

	if err := os.WriteFile(rcPath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", rcPath, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "ctm init written to %s\n", rcPath)
	fmt.Fprintln(cmd.OutOrStdout(), "Restart your shell or source the file to use ctm-use and ctm-unset.")
	return nil
}

func detectShell(shellPath string) (string, error) {
	if shellPath == "" {
		return "", fmt.Errorf("could not detect shell from SHELL. Supported shells: bash, zsh")
	}

	shell := filepath.Base(shellPath)
	switch shell {
	case "bash", "zsh":
		return shell, nil
	default:
		return "", fmt.Errorf("unsupported shell: %s. Supported shells: bash, zsh", shell)
	}
}

func shellRCPath(shell, homeDir string) (string, error) {
	switch shell {
	case "bash":
		return filepath.Join(homeDir, ".bashrc"), nil
	case "zsh":
		return filepath.Join(homeDir, ".zshrc"), nil
	default:
		return "", fmt.Errorf("unsupported shell: %s. Supported shells: bash, zsh", shell)
	}
}

func ctmExecutableDir() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	return filepath.Dir(executablePath), nil
}

func buildInitBlock(ctmDir string) string {
	return fmt.Sprintf(`%s
export PATH=%s:"$PATH"

ctm-use() {
    local ctm_exports
    local ctm_unsets
    ctm_exports=$(ctm use "$@") || return
    ctm_unsets=$(ctm unset) || return
    eval "$ctm_unsets"
    eval "$ctm_exports"
    ctm show
}

ctm-unset() {
    local ctm_unsets
    ctm_unsets=$(ctm unset "$@") || return
    eval "$ctm_unsets"
}
%s
`, initBlockStart, shellQuote(ctmDir), initBlockEnd)
}

func upsertInitBlock(content, block string) string {
	block = strings.TrimRight(block, "\n") + "\n"

	start := strings.Index(content, initBlockStart)
	end := strings.Index(content, initBlockEnd)
	if start >= 0 && end >= start {
		end += len(initBlockEnd)
		if end < len(content) && content[end] == '\n' {
			end++
		}
		return content[:start] + block + content[end:]
	}

	if strings.TrimSpace(content) == "" {
		return block
	}

	if strings.HasSuffix(content, "\n") {
		return content + "\n" + block
	}
	return content + "\n\n" + block
}
