package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tempo/internal/service"
	"tempo/internal/storage"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the webhook scheduler",
	Long: `Start the webhook scheduler to execute configured jobs according to their schedules.

Examples:
  tempo start --foreground
  tempo start --daemon`,
	RunE: runStart,
}

var foreground bool

func init() {
	startCmd.Flags().BoolVarP(&foreground, "foreground", "f", false, "Run in foreground mode (default)")
}

func runStart(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Starting Tempo Scheduler...")

	// Load jobs from storage
	store, err := storage.NewStorage("")
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %v", err)
	}

	jobs := store.GetAllJobs()
	scheduler := service.NewScheduler()

	if len(jobs) == 0 {
		fmt.Println("No jobs configured. Use 'tempo add' to create jobs first.")
		fmt.Println("Starting scheduler in idle mode...")
	} else {
		fmt.Printf("Loading %d job(s)...\n", len(jobs))
		for _, job := range jobs {
			scheduler.AddJob(job)
			fmt.Printf("  ✓ %s: %s %s\n", job.ID, job.Method, job.URL)
		}
		fmt.Println()
	}

	scheduler.Start()

	if foreground {
		fmt.Println("Scheduler running in foreground mode.")
		fmt.Println("Press Ctrl+C to stop")

		// Wait for interrupt signal
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		fmt.Println("\nStopping scheduler...")
		scheduler.Stop()
		fmt.Println("Scheduler stopped.")
	} else {
		fmt.Println("Scheduler started in background mode.")
		fmt.Println("Use 'tempo logs' to view execution logs.")
	}

	return nil
}
