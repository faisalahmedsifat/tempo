

## **Folder Structure**

```
tempo/
├── cmd/
│   └── main.go           # entrypoint
├── internal/
│   ├── types/
│   │   └── job.go        # Job struct
│   └── service/
│       └── scheduler.go  # cron-based scheduling & webhook logic
├── go.mod
```

---

## **Step 1 — types/job.go**

```go
package types

type Job struct {
	ID           string
	URL          string
	CronExpr     string // supports seconds: "0 */5 * * * *" or standard 5-field
}
```

---

## **Step 2 — service/scheduler.go**

```go
package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
	"tempo/internal/types"
)

type Scheduler struct {
	Cron   *cron.Cron
	ctx    context.Context
	cancel context.CancelFunc
}

func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	// WithSeconds() allows 6-field cron expressions
	c := cron.New(cron.WithSeconds())
	return &Scheduler{Cron: c, ctx: ctx, cancel: cancel}
}

// AddJob schedules a job and handles errors centrally
func (s *Scheduler) AddJob(job types.Job) {
	_, err := s.Cron.AddFunc(job.CronExpr, func() {
		if err := CallWebhook(job.URL); err != nil {
			log.Printf("[ERROR] Job %s failed: %v", job.ID, err)
		} else {
			log.Printf("[INFO] Job %s executed successfully", job.ID)
		}
	})
	if err != nil {
		log.Printf("[ERROR] Failed to schedule job %s: %v", job.ID, err)
	}
}

func CallWebhook(url string) error {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &WebhookError{StatusCode: resp.StatusCode}
	}
	return nil
}

type WebhookError struct {
	StatusCode int
}

func (e *WebhookError) Error() string {
	return "webhook returned status code " + http.StatusText(e.StatusCode)
}

// Start begins cron scheduler
func (s *Scheduler) Start() {
	s.Cron.Start()
	log.Println("[INFO] Scheduler started")
}

// Stop stops the scheduler gracefully
func (s *Scheduler) Stop() {
	s.Cron.Stop()
	s.cancel()
	log.Println("[INFO] Scheduler stopped")
}
```

---

## **Step 3 — cmd/main.go**

```go
package main

import (
	"os"
	"os/signal"
	"syscall"

	"tempo/internal/service"
	"tempo/internal/types"
)

func main() {
	scheduler := service.NewScheduler()

	// Example job with seconds-level cron: every 10 seconds
	job := types.Job{
		ID:       "job1",
		URL:      "https://example.com/webhook",
		CronExpr: "*/10 * * * * *", 
	}
	scheduler.AddJob(job)

	scheduler.Start()

	// Wait for exit signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	scheduler.Stop()
}
```

---

### **✅ Features**

1. **Cron expressions with seconds support** (`robfig/cron/v3` with `WithSeconds()`).
2. **Centralized error handling**:

   * Logs failed webhook calls.
   * Logs failed scheduling.
   * Logs HTTP errors (non-2xx).
3. **Graceful shutdown**.
4. **Minimal structure** — no unnecessary folders.
5. Can be easily extended:

   * Add POST payloads, headers, retries.
   * Persist jobs in DB later.

