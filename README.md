# Tempo

A simple, reliable webhook scheduler built with Go. Schedule HTTP requests to run on recurring intervals using cron expressions with built-in error handling and logging.

## Overview

Tempo is a lightweight webhook scheduler that allows you to schedule HTTP GET and POST requests to any endpoint using cron expressions. Perfect for periodic API calls, health checks, data synchronization, and automated notifications.

### Key Features

- **Cron-based Scheduling** - Uses standard cron expressions with second-level precision
- **HTTP Method Support** - GET and POST requests with custom headers and body
- **Error Handling** - Comprehensive error handling with custom error types
- **Real-time Logging** - Detailed execution logs with success/failure status
- **Graceful Shutdown** - Clean shutdown with signal handling
- **Simple Architecture** - Clean, maintainable code structure

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/faisalahmedsifat/tempo.git
cd tempo

# Install dependencies
go mod download

# Run the scheduler
go run cmd/main.go
```

### Project Structure

```
tempo/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── types/
│   │   └── job.go           # Job struct definition
│   └── service/
│       └── scheduler.go     # Scheduling logic and webhook execution
├── go.mod
└── go.sum
```

## Usage

### Creating Jobs

Jobs are defined using the `Job` struct:

```go
job := types.Job{
    ID:       "my-webhook",
    URL:      "https://api.example.com/webhook",
    Method:   "POST",
    CronExpr: "*/30 * * * * *", // Every 30 seconds
    Headers: map[string]string{
        "Content-Type": "application/json",
        "Authorization": "Bearer your-token",
    },
    Body: `{"message": "Hello from Tempo!"}`,
}
```

### Basic Example

```go
package main

import (
    "tempo/internal/service"
    "tempo/internal/types"
)

func main() {
    // Create scheduler
    scheduler := service.NewScheduler()

    // Define a job
    job := types.Job{
        ID:       "health-check",
        URL:      "https://httpbin.org/get",
        Method:   "GET",
        CronExpr: "*/10 * * * * *", // Every 10 seconds
    }

    // Add job to scheduler
    scheduler.AddJob(job)

    // Start scheduler
    scheduler.Start()

    // Keep running (add signal handling in production)
    select {}
}
```

### Supported Job Types

#### GET Request
```go
job := types.Job{
    ID:       "api-health-check",
    URL:      "https://api.example.com/health",
    Method:   "GET",
    CronExpr: "0 */5 * * * *", // Every 5 minutes
    Headers: map[string]string{
        "User-Agent": "Tempo-Scheduler/1.0",
    },
}
```

#### POST Request with JSON
```go
job := types.Job{
    ID:       "slack-notification",
    URL:      "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK",
    Method:   "POST",
    CronExpr: "0 0 9 * * 1-5", // Weekdays at 9 AM
    Headers: map[string]string{
        "Content-Type": "application/json",
    },
    Body: `{
        "text": "Good morning! Daily standup reminder.",
        "channel": "#general"
    }`,
}
```

#### POST Request with Form Data
```go
job := types.Job{
    ID:       "form-submission",
    URL:      "https://api.example.com/submit",
    Method:   "POST",
    CronExpr: "0 0 * * * *", // Every hour
    Headers: map[string]string{
        "Content-Type": "application/x-www-form-urlencoded",
    },
    Body: "key1=value1&key2=value2",
}
```

## Cron Expression Format

Tempo supports 6-field cron expressions with second precision:

```
┌───────────── second (0-59)
│ ┌───────────── minute (0-59)
│ │ ┌───────────── hour (0-23)
│ │ │ ┌───────────── day of month (1-31)
│ │ │ │ ┌───────────── month (1-12)
│ │ │ │ │ ┌───────────── day of week (0-6, 0=Sunday)
│ │ │ │ │ │
* * * * * *
```

### Common Patterns

| Expression | Description |
|------------|-------------|
| `*/10 * * * * *` | Every 10 seconds |
| `0 */5 * * * *` | Every 5 minutes |
| `0 0 * * * *` | Every hour |
| `0 0 9 * * *` | Every day at 9 AM |
| `0 0 9 * * 1-5` | Weekdays at 9 AM |
| `0 0 0 1 * *` | First day of every month |
| `0 30 14 * * 6` | Every Saturday at 2:30 PM |

## API Reference

### Types

#### Job
```go
type Job struct {
    ID       string            `json:"id"`       // Unique identifier
    URL      string            `json:"url"`      // Target webhook URL
    CronExpr string            `json:"cron_expr"` // Cron expression
    Method   string            `json:"method"`   // HTTP method (GET/POST)
    Headers  map[string]string `json:"headers"`  // HTTP headers
    Body     string            `json:"body"`     // Request body (for POST)
}
```

### Scheduler Methods

#### NewScheduler()
Creates a new scheduler instance.

```go
scheduler := service.NewScheduler()
```

#### AddJob(job Job)
Adds a job to the scheduler.

```go
scheduler.AddJob(job)
```

#### Start()
Starts the scheduler and begins executing jobs.

```go
scheduler.Start()
```

#### Stop()
Gracefully stops the scheduler.

```go
scheduler.Stop()
```

## Error Handling

Tempo includes comprehensive error handling:

### WebhookError
Custom error type for webhook-specific failures:

```go
type WebhookError struct {
    StatusCode int    // HTTP status code
    Message    string // Error message
}
```

### Error Scenarios

- **Network errors**: Connection timeouts, DNS failures
- **HTTP errors**: 4xx and 5xx status codes
- **Scheduling errors**: Invalid cron expressions
- **Request errors**: Invalid URLs, malformed requests

## Logging

Tempo provides detailed logging for monitoring and debugging:

```
Calling webhook: POST https://api.example.com/webhook
[INFO] Job slack-notification executed successfully

Calling webhook: GET https://api.invalid.com/endpoint
Error calling webhook: no such host
[ERROR] Job health-check failed: webhook returned status code: 500, message: Internal server error
```

## Testing

### Test with Online Services

```go
// Test with httpbin.org
jobs := []types.Job{
    {
        ID:       "test-get",
        URL:      "https://httpbin.org/get",
        Method:   "GET",
        CronExpr: "*/5 * * * * *",
    },
    {
        ID:       "test-post",
        URL:      "https://httpbin.org/post",
        Method:   "POST",
        CronExpr: "*/10 * * * * *",
        Body:     `{"test": "data"}`,
    },
}
```

### Create a Test Webhook

Use [webhook.site](https://webhook.site) to create a test endpoint and monitor incoming requests.

## Development

### Prerequisites

- Go 1.21 or higher
- Git

### Setup

1. Fork and clone the repository
2. Install dependencies: `go mod download`
3. Make your changes
4. Test your changes: `go run cmd/main.go`
5. Submit a pull request

### Code Structure

- **`internal/types/`**: Data structures and types
- **`internal/service/`**: Business logic and scheduling
- **`cmd/`**: Application entry points

## Configuration

### Environment Variables

Currently, all configuration is done through code. Future versions may support:

| Variable | Description | Default |
|----------|-------------|---------|
| `TEMPO_LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `TEMPO_TIMEOUT` | HTTP request timeout | `10s` |
| `TEMPO_MAX_RETRIES` | Maximum retry attempts | `3` |

## Limitations

- Jobs are defined in code and require restart to modify
- No persistence layer (jobs don't survive restarts)
- No web interface or API for job management
- Limited to GET and POST HTTP methods
- No job dependencies or complex workflows

## Roadmap

- [ ] REST API for job management
- [ ] Web dashboard
- [ ] Job persistence (database storage)
- [ ] More HTTP methods (PUT, DELETE, PATCH)
- [ ] Retry logic with exponential backoff
- [ ] Job execution history
- [ ] Metrics and monitoring
- [ ] Configuration file support
- [ ] Job templates and bulk operations
- [ ] Authentication and authorization

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Dependencies

- [robfig/cron](https://github.com/robfig/cron) - Cron expression parsing and scheduling

## Acknowledgments

- Inspired by the need for simple, reliable webhook scheduling
- Built with Go's excellent standard library
- Thanks to the robfig/cron library for robust cron expression support

---

*Simple. Reliable. Efficient. Never miss a scheduled webhook again.*