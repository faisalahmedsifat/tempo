package commands

import (
	"fmt"
	"strings"
	"tempo/internal/service"
	"tempo/internal/storage"
	"tempo/internal/types"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [job-id]",
	Short: "Run a job immediately",
	Long: `Execute a webhook job immediately for testing purposes.

Examples:
  tempo run health-check
  tempo run --url "https://api.example.com/test" --method POST`,
	Args: cobra.MaximumNArgs(1),
	RunE: runExecute,
}

var (
	runURL     string
	runMethod  string
	runBody    string
	runHeaders []string
)

func init() {
	runCmd.Flags().StringVarP(&runURL, "url", "u", "", "Webhook URL (for one-off execution)")
	runCmd.Flags().StringVarP(&runMethod, "method", "m", "GET", "HTTP method")
	runCmd.Flags().StringVarP(&runBody, "body", "b", "", "Request body")
	runCmd.Flags().StringSliceVarP(&runHeaders, "header", "H", []string{}, "HTTP headers (format: 'Key=Value')")
}

func runExecute(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && runURL == "" {
		return fmt.Errorf("either job ID or --url is required")
	}

	if len(args) > 0 {
		jobID := args[0]

		store, err := storage.NewStorage("")
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %v", err)
		}

		job, exists := store.GetJob(jobID)
		if !exists {
			return fmt.Errorf("job '%s' not found", jobID)
		}

		fmt.Printf("Running job '%s'...\n", jobID)
		fmt.Printf("Executing: %s %s\n", job.Method, job.URL)
		if job.Body != "" {
			fmt.Printf("Body: %.50s...\n", job.Body)
		}
		if len(job.Headers) > 0 {
			fmt.Printf("Headers: %v\n", job.Headers)
		}

		err = service.CallWebhook(job)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return err
		}

		fmt.Println("✅ Job executed successfully")
		return nil
	}

	// One-off execution
	if runURL == "" {
		return fmt.Errorf("URL is required for one-off execution")
	}

	// Parse headers
	headers := make(map[string]string)
	for _, header := range runHeaders {
		parts := strings.SplitN(header, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header format: %s. Use 'Key=Value'", header)
		}
		headers[parts[0]] = parts[1]
	}

	job := types.Job{
		ID:      "one-off",
		URL:     runURL,
		Method:  runMethod,
		Body:    runBody,
		Headers: headers,
	}

	fmt.Printf("Executing: %s %s\n", job.Method, job.URL)
	if job.Body != "" {
		fmt.Printf("Body: %.50s...\n", job.Body)
	}
	if len(headers) > 0 {
		fmt.Printf("Headers: %v\n", headers)
	}

	err := service.CallWebhook(job)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return err
	}

	fmt.Println("✅ Webhook executed successfully")
	return nil
}
