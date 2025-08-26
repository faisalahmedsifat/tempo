package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tempo/internal/storage"
	"tempo/internal/types"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [job-id]",
	Short: "Add a new webhook job",
	Long: `Add a new webhook job to the scheduler. You can specify all parameters via flags
or use --interactive for a guided setup.

Examples:
  tempo add health-check --url "https://api.example.com/health" --schedule "*/30 * * * * *"
  tempo add --interactive`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAdd,
}

var (
	jobURL      string
	jobMethod   string
	jobSchedule string
	jobBody     string
	jobHeaders  []string
	interactive bool
)

func init() {
	addCmd.Flags().StringVarP(&jobURL, "url", "u", "", "Webhook URL")
	addCmd.Flags().StringVarP(&jobMethod, "method", "m", "GET", "HTTP method (GET, POST, PUT, DELETE)")
	addCmd.Flags().StringVarP(&jobSchedule, "schedule", "s", "", "Cron schedule expression (e.g., '*/30 * * * * *')")
	addCmd.Flags().StringVarP(&jobBody, "body", "b", "", "Request body")
	addCmd.Flags().StringSliceVarP(&jobHeaders, "header", "H", []string{}, "HTTP headers (format: 'Key=Value')")
	addCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode for guided setup")
}

func runAdd(cmd *cobra.Command, args []string) error {
	var jobID string
	if len(args) > 0 {
		jobID = args[0]
	}

	if interactive {
		return runInteractiveAdd(jobID)
	}

	// Validate required fields
	if jobURL == "" {
		return fmt.Errorf("URL is required. Use --url or --interactive")
	}
	if jobSchedule == "" {
		return fmt.Errorf("schedule is required. Use --schedule or --interactive")
	}
	if jobID == "" {
		return fmt.Errorf("job ID is required")
	}

	// Parse headers
	headers := make(map[string]string)
	for _, header := range jobHeaders {
		parts := strings.SplitN(header, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header format: %s. Use 'Key=Value'", header)
		}
		headers[parts[0]] = parts[1]
	}

	job := types.Job{
		ID:       jobID,
		URL:      jobURL,
		Method:   jobMethod,
		CronExpr: jobSchedule,
		Body:     jobBody,
		Headers:  headers,
	}

	// Add job to storage
	store, err := storage.NewStorage("")
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %v", err)
	}

	if err := store.AddJob(job); err != nil {
		return fmt.Errorf("failed to add job: %v", err)
	}

	fmt.Printf("✓ Added job '%s': %s %s\n", jobID, jobMethod, jobURL)
	if jobBody != "" {
		fmt.Printf("  Body: %.50s...\n", jobBody)
	}
	if len(headers) > 0 {
		fmt.Printf("  Headers: %v\n", headers)
	}

	return nil
}

func runInteractiveAdd(jobID string) error {
	reader := bufio.NewReader(os.Stdin)

	// Get Job ID
	if jobID == "" {
		fmt.Print("Job ID: ")
		id, _ := reader.ReadString('\n')
		jobID = strings.TrimSpace(id)
		if jobID == "" {
			return fmt.Errorf("job ID cannot be empty")
		}
	}

	// Get URL
	fmt.Print("Webhook URL: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Get Method
	fmt.Print("HTTP Method [GET]: ")
	method, _ := reader.ReadString('\n')
	method = strings.TrimSpace(method)
	if method == "" {
		method = "GET"
	}
	method = strings.ToUpper(method)

	// Get Schedule
	fmt.Print("Cron Schedule (e.g., '*/30 * * * * *'): ")
	schedule, _ := reader.ReadString('\n')
	schedule = strings.TrimSpace(schedule)
	if schedule == "" {
		return fmt.Errorf("schedule cannot be empty")
	}

	// Get Body
	fmt.Print("Request Body (optional): ")
	body, _ := reader.ReadString('\n')
	body = strings.TrimSpace(body)

	// Get Headers
	headers := make(map[string]string)
	fmt.Println("Add headers (press Enter when done):")
	for {
		fmt.Print("Header (Key=Value): ")
		header, _ := reader.ReadString('\n')
		header = strings.TrimSpace(header)
		if header == "" {
			break
		}
		parts := strings.SplitN(header, "=", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid format. Use 'Key=Value'")
			continue
		}
		headers[parts[0]] = parts[1]
	}

	job := types.Job{
		ID:       jobID,
		URL:      url,
		Method:   method,
		CronExpr: schedule,
		Body:     body,
		Headers:  headers,
	}

	// Add job to storage
	store, err := storage.NewStorage("")
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %v", err)
	}

	if err := store.AddJob(job); err != nil {
		return fmt.Errorf("failed to add job: %v", err)
	}

	fmt.Printf("\n✓ Added job '%s': %s %s\n", jobID, method, url)
	if body != "" {
		fmt.Printf("  Body: %.50s...\n", body)
	}
	if len(headers) > 0 {
		fmt.Printf("  Headers: %v\n", headers)
	}

	return nil
}
