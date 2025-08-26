package service

// imports
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"tempo/internal/types"
	"time"

	"github.com/robfig/cron/v3"
)

// structs

type Scheduler struct {
	Cron   *cron.Cron
	ctx    context.Context
	cancel context.CancelFunc
}

/*
* WebhookError is an error type for webhook errors
* It contains the status code and message of the webhook error
 */
type WebhookError struct {
	StatusCode int
	Message    string
}

/*
* Error returns a string representation of the webhook error
* It returns the status code and message of the webhook error
 */
func (e *WebhookError) Error() string {
	return fmt.Sprintf("webhook returned status code: %d, message: %s", e.StatusCode, e.Message)
}

/*
* NewScheduler creates a new scheduler instance
* It creates a context with a cancel function and a cron instance
* It returns a new scheduler instance
 */
func NewScheduler() *Scheduler {
	// create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// create a cron instance
	cron := cron.New(cron.WithSeconds())

	// return a new scheduler instance
	return &Scheduler{
		Cron:   cron,
		ctx:    ctx,
		cancel: cancel,
	}
}

/*
* CallWebhook calls the webhook
* It calls the webhook and returns an error if the webhook returns a status code >= 400
 */
func CallWebhook(job types.Job) error {
	log.Printf("Calling webhook: %v, method: %v, body: %v, headers: %v", job.URL, job.Method, job.Body, job.Headers)

	// create a new http client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// create request with body, method and headers
	req, err := http.NewRequest(job.Method, job.URL, strings.NewReader(job.Body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return &WebhookError{
			StatusCode: 500,
			Message:    "Error creating request",
		}
	}

	// add headers to request
	for key, value := range job.Headers {
		req.Header.Add(key, value)
	}

	// send request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return &WebhookError{
			StatusCode: 500,
			Message:    "Error sending request",
		}
	}

	defer resp.Body.Close()

	// check if status code is >= 400
	if resp.StatusCode >= 400 {
		return &WebhookError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// log response
	log.Printf("Response: %v", resp)

	return nil
}

/*
* AddJob adds a new job to the scheduler
* It adds a new job to the scheduler using the cron expression
* It calls the webhook and logs the result
 */
func (s *Scheduler) AddJob(job types.Job) {
	_, err := s.Cron.AddFunc(job.CronExpr, func() {
		err := CallWebhook(job)

		if err != nil {
			log.Printf("Error calling webhook: %v", err)
		} else {
			log.Printf("Job %s executed successfully", job.URL)
		}

	})

	if err != nil {
		log.Printf("Error adding job: %v", err)
	}
}

/*
* Start starts the scheduler
* It starts the scheduler and logs a message
 */
func (s *Scheduler) Start() {
	s.Cron.Start()
	log.Println("[INFO] Scheduler started")
}

/*
* Stop stops the scheduler
* It stops the scheduler and logs a message
 */
func (s *Scheduler) Stop() {
	s.Cron.Stop()
	s.cancel()

	log.Println("[INFO] Scheduler stopped")
}
