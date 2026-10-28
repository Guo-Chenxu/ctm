package cmd

import (
	"fmt"
	"io"
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
	Long:  `Write PATH setup and shell completion to the shell rc file.`,
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

	completionPath, err := writeCompletionScript(shell, homeDir)
	if err != nil {
		return err
	}

	currentContent, err := os.ReadFile(rcPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read %s: %w", rcPath, err)
	}

	block := buildInitBlock(ctmDir, completionPath)
	updatedContent := upsertInitBlock(string(currentContent), block)

	if err := os.WriteFile(rcPath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", rcPath, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "ctm init written to %s\n", rcPath)
	fmt.Fprintf(cmd.OutOrStdout(), "Completion script written to %s\n", completionPath)
	fmt.Fprintln(cmd.OutOrStdout(), "Restart your shell or source the rc file to activate ctm.")
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

// writeCompletionScript generates the completion script for the shell and returns its path
func writeCompletionScript(shell, homeDir string) (string, error) {
	path := filepath.Join(homeDir, ".ctm", "completion."+shell)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("failed to create completion directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create %s: %w", path, err)
	}
	defer file.Close()

	if err := generateCompletion(shell, file); err != nil {
		return "", err
	}
	return path, nil
}

func generateCompletion(shell string, w io.Writer) error {
	switch shell {
	case "bash":
		return rootCmd.GenBashCompletionV2(w, true)
	case "zsh":
		return rootCmd.GenZshCompletion(w)
	default:
		return fmt.Errorf("unsupported shell: %s. Supported shells: bash, zsh", shell)
	}
}

func buildInitBlock(ctmDir, completionPath string) string {
	return fmt.Sprintf(`%s
export PATH=%s:"$PATH"

if [ -f %s ]; then
    . %s
fi
%s
`, initBlockStart, shellQuote(ctmDir), shellQuote(completionPath), shellQuote(completionPath), initBlockEnd)
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
