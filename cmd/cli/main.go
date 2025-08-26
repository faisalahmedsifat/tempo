package main

import (
	"fmt"
	"os"
	"tempo/cmd/cli/commands"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tempo",
	Short: "Tempo - A simple webhook scheduler",
	Long: `Tempo is a lightweight webhook scheduler that allows you to schedule HTTP requests 
to webhooks at specified intervals using cron expressions.

Examples:
  tempo add health-check --url "https://api.example.com/health" --schedule "*/30 * * * * *"
  tempo list
  tempo start --foreground
  tempo logs --follow`,
}

func main() {
	// Register all commands
	commands.RegisterCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
