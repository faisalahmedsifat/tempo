package commands

import (
	"fmt"
	"tempo/internal/storage"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured jobs",
	Long: `List all configured webhook jobs with their details including ID, method, URL, and schedule.

Examples:
  tempo list
  tempo list --verbose`,
	RunE: runList,
}

var verbose bool

func init() {
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed job information")
}

func runList(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage("")
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %v", err)
	}

	jobs := store.GetAllJobs()

	if len(jobs) == 0 {
		fmt.Println("No jobs configured yet.")
		fmt.Println("Use 'tempo add' to create your first job.")

		if verbose {
			fmt.Println("\nExample jobs:")
			fmt.Println("  tempo add health-check --url 'https://api.example.com/health' --schedule '*/30 * * * * *'")
			fmt.Println("  tempo add weekly-report --url 'https://api.example.com/reports' --method POST --schedule '0 0 9 * * 1'")
		}
		return nil
	}

	fmt.Printf("Found %d job(s):\n\n", len(jobs))
	for _, job := range jobs {
		fmt.Printf("ID: %s\n", job.ID)
		fmt.Printf("  Method: %s\n", job.Method)
		fmt.Printf("  URL: %s\n", job.URL)
		fmt.Printf("  Schedule: %s\n", job.CronExpr)
		if job.Body != "" {
			fmt.Printf("  Body: %.50s...\n", job.Body)
		}
		if len(job.Headers) > 0 {
			fmt.Printf("  Headers: %v\n", job.Headers)
		}
		fmt.Println()
	}

	return nil
}
