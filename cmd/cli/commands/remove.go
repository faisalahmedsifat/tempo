package commands

import (
	"fmt"
	"tempo/internal/storage"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [job-id]",
	Short: "Remove a webhook job",
	Long: `Remove a webhook job from the scheduler.

Examples:
  tempo remove health-check
  tempo remove --all`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRemove,
}

var removeAll bool

func init() {
	removeCmd.Flags().BoolVarP(&removeAll, "all", "a", false, "Remove all jobs")
}

func runRemove(cmd *cobra.Command, args []string) error {
	store, err := storage.NewStorage("")
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %v", err)
	}

	if removeAll {
		if err := store.RemoveAllJobs(); err != nil {
			return fmt.Errorf("failed to remove all jobs: %v", err)
		}
		fmt.Println("✓ Removed all jobs")
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("job ID is required")
	}

	jobID := args[0]

	if err := store.RemoveJob(jobID); err != nil {
		return fmt.Errorf("failed to remove job: %v", err)
	}

	fmt.Printf("✓ Removed job '%s'\n", jobID)
	return nil
}
