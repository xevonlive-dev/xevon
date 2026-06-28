package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xevonlive-dev/xevon/internal/config"
	"github.com/xevonlive-dev/xevon/pkg/database"
	"github.com/xevonlive-dev/xevon/pkg/terminal"
	"github.com/xevonlive-dev/xevon/public"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize xevon with default configuration and preset data",
	Long: `Create the ~/.xevon directory with a default config file, database schema,
scanning profiles, prompt templates, extensions, and SAST rules.

Safe to run on an existing installation — skips components that already exist
unless --force is passed, which regenerates the config (with a new API key)
and re-extracts all preset data.`,
	RunE: runInitCmd,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInitCmd(cmd *cobra.Command, args []string) error {
	defer syncLogger()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	xevonDir := filepath.Join(homeDir, ".xevon")
	settingsPath := filepath.Join(xevonDir, settingsFileName)

	if _, err := os.Stat(settingsPath); err == nil && !globalForce {
		fmt.Fprintf(os.Stderr, "%s Already initialized (%s exists). Use --force to reinitialize.\n",
			terminal.InfoSymbol(), terminal.Cyan(config.ContractPath(settingsPath)))
		return nil
	}

	if err := os.MkdirAll(xevonDir, 0755); err != nil {
		return fmt.Errorf("failed to create .xevon directory: %w", err)
	}

	// Write default config with a fresh API key
	configData := bytes.Replace(
		public.DefaultConfigYAML,
		[]byte(`auth_api_key: "auto-generated-on-first-run"`),
		[]byte(fmt.Sprintf(`auth_api_key: "%s"`, config.GenerateRandomHex(40))),
		1,
	)
	if err := os.WriteFile(settingsPath, configData, 0600); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	settings := config.DefaultSettings()

	// Initialize database
	if settings.Database.Enabled {
		db, err := database.NewDB(&settings.Database)
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer func() { _ = db.Close() }()

		ctx := context.Background()
		if err := db.CreateSchema(ctx); err != nil {
			return fmt.Errorf("failed to create database schema: %w", err)
		}
		if err := db.SeedDefaults(ctx); err != nil {
			return fmt.Errorf("failed to seed default data: %w", err)
		}
	}

	if globalForce {
		// Remove sentinel directories so bootstrap helpers re-extract
		for _, sub := range []string{"profiles", "extensions", "prompts"} {
			_ = os.RemoveAll(filepath.Join(xevonDir, sub))
		}
	}

	bootstrapDefaultProfiles(xevonDir)
	bootstrapExtensions(xevonDir)
	bootstrapPrompts(xevonDir)

	// Trigger the core dep install (chromium + nuclei templates) inline so
	// explicit `xevon init` is a complete first-run setup. Without --force
	// the marker may already exist and ensureCoreDeps short-circuits; under
	// --force we wipe the marker first so the dep check re-runs on the same
	// invocation (handy when --force is used to recover a broken setup).
	if globalForce {
		_ = os.Remove(filepath.Join(xevonDir, initMarkerFileName))
	}
	if err := ensureCoreDeps(); err != nil {
		return fmt.Errorf("core dependency install failed: %w", err)
	}

	fmt.Fprintf(os.Stderr, "%s %s\n", terminal.SuccessSymbol(), terminal.BoldGreen("xevon initialized successfully!"))
	fmt.Fprintf(os.Stderr, "  %s Config:   %s\n", terminal.InfoSymbol(), terminal.Cyan(config.ContractPath(settingsPath)))
	fmt.Fprintf(os.Stderr, "  %s Database: %s\n", terminal.InfoSymbol(), terminal.Cyan(config.ContractPath(config.ExpandPath(settings.Database.SQLite.Path))))
	fmt.Fprintf(os.Stderr, "  %s Docs:     %s\n", terminal.InfoSymbol(), terminal.Cyan("https://docs.xevon.live"))

	return nil
}
