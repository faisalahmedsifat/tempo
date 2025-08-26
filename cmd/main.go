package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tempo/internal/service"
	"tempo/internal/types"
)

func main() {
	fmt.Println("🚀 Starting Tempo Scheduler Test...")

	scheduler := service.NewScheduler()
	jobs := []types.Job{
		// GET request (backwards compatible)
		{
			ID:       "get-test",
			URL:      "https://httpbin.org/get",
			CronExpr: "*/10 * * * * *",
			Method:   "GET",
		},

		// POST request with JSON body
		{
			ID:       "post-json",
			URL:      "https://httpbin.org/post",
			CronExpr: "*/15 * * * * *",
			Method:   "POST",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"User-Agent":   "Tempo-Scheduler/1.0",
			},
			Body: `{"message": "Hello from Tempo!", "timestamp": "2025-01-01"}`,
		},

		// POST request with form data
		{
			ID:       "post-form",
			URL:      "https://httpbin.org/post",
			CronExpr: "*/20 * * * * *",
			Method:   "POST",
			Headers: map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			},
			Body: "key1=value1&key2=value2",
		},

		// POST request to webhook.site (create your own URL)
		{
			ID:       "webhook-site",
			URL:      "https://webhook.site/unique-id", // Replace with your webhook.site URL
			CronExpr: "*/30 * * * * *",
			Method:   "POST",
			Headers: map[string]string{
				"Content-Type":    "application/json",
				"X-Custom-Header": "tempo-test",
			},
			Body: `{
                "event": "scheduled_job",
                "job_id": "webhook-site",
                "data": {
                    "message": "This is a test from Tempo scheduler",
                    "timestamp": "2025-01-01T00:00:00Z"
                }
            }`,
		},

		// POST with empty body
		{
			ID:       "post-empty",
			URL:      "https://httpbin.org/post",
			CronExpr: "*/25 * * * * *",
			Method:   "POST",
			Headers: map[string]string{
				"Authorization": "Bearer your-token-here",
			},
			Body: "",
		},
	}

	for _, job := range jobs {
		scheduler.AddJob(job)
		fmt.Printf("  ✓ %s: %s %s\n", job.ID, job.Method, job.URL)
		if job.Body != "" {
			fmt.Printf("    Body: %.50s...\n", job.Body)
		}
		if len(job.Headers) > 0 {
			fmt.Printf("    Headers: %v\n", job.Headers)
		}
		fmt.Println()
	}

	scheduler.Start()

	fmt.Println("Press Ctrl+C to stop")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	scheduler.Stop()
}
