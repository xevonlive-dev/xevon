package configcmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xevonlive-dev/xevon/internal/config"
	"github.com/xevonlive-dev/xevon/pkg/terminal"
)

func newCleanCmd(deps Deps, example string) *cobra.Command {
	return &cobra.Command{
		Use:     "clean",
		Short:   "Reset xevon to a clean state",
		Long:    "Remove the ~/.xevon/ directory (config, database, extensions) and regenerate fresh defaults.",
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigClean(deps)
		},
	}
}

func runConfigClean(deps Deps) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	xevonDir := filepath.Join(homeDir, ".xevon")

	displayDir := config.ContractPath(xevonDir)

	// Check if directory exists
	if _, err := os.Stat(xevonDir); os.IsNotExist(err) {
		fmt.Printf("%s Nothing to clean — %s does not exist.\n", terminal.InfoSymbol(), displayDir)
		return nil
	}

	fmt.Printf("%s This will remove %s (config, database, and all local data)\n", terminal.BoldRed(terminal.SymbolFailed+" Warn:"), terminal.Cyan(displayDir))

	if !deps.Force() {
		fmt.Print("\nProceed? (type 'yes' to confirm): ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	if err := os.RemoveAll(xevonDir); err != nil {
		return fmt.Errorf("failed to remove %s: %w", xevonDir, err)
	}

	fmt.Printf("%s Removed %s\n", terminal.SuccessSymbol(), displayDir)

	// Regenerate fresh defaults
	if err := deps.Reinitialize(); err != nil {
		return fmt.Errorf("failed to reinitialize: %w", err)
	}

	return nil
}
