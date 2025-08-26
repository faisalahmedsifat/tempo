# Tempo - Webhook Scheduler

A lightweight, command-line webhook scheduler that allows you to schedule HTTP requests to webhooks at specified intervals using cron expressions.

## Features

- **Simple CLI Interface**: Easy-to-use commands for managing webhook jobs
- **Persistent Storage**: Jobs are saved to JSON files and survive restarts
- **Cron Scheduling**: Support for standard cron expressions
- **Multiple HTTP Methods**: GET, POST, PUT, DELETE support
- **Custom Headers & Body**: Full control over request configuration
- **Real-time Testing**: Test webhooks immediately with `tempo run`
- **Import/Export**: Backup and restore job configurations
- **Interactive Mode**: Guided setup for complex webhooks

## Installation

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd tempo

# Build the CLI
go build -o tempo cmd/cli/main.go

# Install globally (optional)
go install ./cmd/cli
```

Or, use the provided `Makefile` to build the CLI:

## Quick Start

### 1. Add Your First Job

```bash
# Simple GET request every 30 seconds
tempo add health-check --url "https://api.example.com/health" --schedule "*/30 * * * * *"

# POST request with JSON body
tempo add weekly-report \
  --url "https://api.example.com/reports" \
  --method POST \
  --schedule "0 0 9 * * 1" \
  --header "Content-Type=application/json" \
  --body '{"template": "weekly_metrics"}'
```

### 2. List Your Jobs

```bash
tempo list
```

### 3. Test a Job

```bash
# Test a saved job
tempo run health-check

# Test a one-off webhook
tempo run --url "https://httpbin.org/get" --method GET
```

### 4. Start the Scheduler

```bash
# Run in foreground (see logs in real-time)
tempo start --foreground

# Run in background
tempo start
```

## CLI Commands

### `tempo add [job-id]`

Add a new webhook job to the scheduler.

**Flags:**
- `--url, -u`: Webhook URL (required)
- `--method, -m`: HTTP method (GET, POST, PUT, DELETE) [default: GET]
- `--schedule, -s`: Cron schedule expression (required)
- `--body, -b`: Request body
- `--header, -H`: HTTP headers (format: 'Key=Value')
- `--interactive, -i`: Interactive mode for guided setup

**Examples:**
```bash
# Simple health check
tempo add health-check --url "https://api.example.com/health" --schedule "*/30 * * * * *"

# Complex webhook with headers and body
tempo add slack-notify \
  --url "https://hooks.slack.com/services/xxx/yyy/zzz" \
  --method POST \
  --schedule "0 9 * * 1-5" \
  --header "Content-Type=application/json" \
  --body '{"text": "Daily reminder!"}'

# Interactive mode
tempo add --interactive
```

### `tempo list`

List all configured jobs.

**Flags:**
- `--verbose, -v`: Show detailed job information

**Example:**
```bash
tempo list --verbose
```

### `tempo run [job-id]`

Execute a webhook job immediately for testing.

**Flags:**
- `--url, -u`: Webhook URL (for one-off execution)
- `--method, -m`: HTTP method [default: GET]
- `--body, -b`: Request body
- `--header, -H`: HTTP headers (format: 'Key=Value')

**Examples:**
```bash
# Run a saved job
tempo run health-check

# Test a one-off webhook
tempo run --url "https://httpbin.org/post" --method POST --body '{"test": "data"}'
```

### `tempo start`

Start the webhook scheduler.

**Flags:**
- `--foreground, -f`: Run in foreground mode (default)

**Example:**
```bash
tempo start --foreground
```

### `tempo logs [job-id]`

View execution logs for webhook jobs.

**Flags:**
- `--follow, -f`: Follow logs in real-time
- `--since, -s`: Show logs since time (e.g., '1h', '30m', '2024-01-01')
- `--limit, -n`: Number of log entries to show [default: 50]

**Examples:**
```bash
tempo logs
tempo logs health-check --follow
tempo logs --since "1h"
```

### `tempo remove [job-id]`

Remove a webhook job from the scheduler.

**Flags:**
- `--all, -a`: Remove all jobs

**Examples:**
```bash
tempo remove health-check
tempo remove --all
```

### `tempo export [filename]`

Export job configurations to a file.

**Flags:**
- `--format, -f`: Export format (json, yaml) [default: json]

**Examples:**
```bash
tempo export
tempo export backup.json
tempo export --format yaml jobs.yaml
```

### `tempo import [filename]`

Import job configurations from a file.

**Flags:**
- `--format, -f`: Import format (json, yaml) [default: json]

**Examples:**
```bash
tempo import backup.json
tempo import --format yaml jobs.yaml
```

## Cron Schedule Format

Tempo uses the standard cron format with seconds precision:

```
┌───────────── second (0 - 59)
│ ┌───────────── minute (0 - 59)
│ │ ┌───────────── hour (0 - 23)
│ │ │ ┌───────────── day of the month (1 - 31)
│ │ │ │ ┌───────────── month (1 - 12)
│ │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday)
│ │ │ │ │ │
* * * * * *
```

**Common Examples:**
- `*/30 * * * * *` - Every 30 seconds
- `0 */5 * * * *` - Every 5 minutes
- `0 0 * * * *` - Every hour
- `0 0 9 * * 1-5` - Weekdays at 9 AM
- `0 0 0 1 * *` - First day of every month

## Real-World Examples

### DevOps Health Monitoring

```bash
# Add health checks for multiple services
tempo add user-service-health \
  --url "https://user-api.company.com/health" \
  --schedule "*/30 * * * * *"

tempo add payment-service-health \
  --url "https://payment-api.company.com/health" \
  --schedule "*/30 * * * * *"

# Start scheduler
tempo start --foreground

# Monitor in real-time
tempo logs --follow
```

### Weekly Reports

```bash
# Interactive setup for complex webhook
tempo add --interactive
# Job ID: investor-report
# Webhook URL: https://api.company.com/reports/weekly
# HTTP Method [GET]: POST
# Cron Schedule: 0 0 9 * * 1
# Header: Authorization=Bearer abc123
# Header: Content-Type=application/json
# Request Body: {"template": "weekly_metrics", "recipients": ["investor@vc.com"]}

# Test the webhook immediately
tempo run investor-report

# View execution history
tempo logs investor-report
```

### E-commerce Inventory Sync

```bash
# Add inventory sync job
tempo add inventory-sync \
  --url "https://warehouse-api.company.com/sync" \
  --method POST \
  --schedule "0 */15 * * * *" \
  --header "API-Key=warehouse_123" \
  --header "Content-Type=application/json" \
  --body '{"source": "shopify", "target": "warehouse"}'

# Pause during maintenance
tempo remove inventory-sync

# Resume after maintenance  
tempo add inventory-sync \
  --url "https://warehouse-api.company.com/sync" \
  --method POST \
  --schedule "0 */15 * * * *" \
  --header "API-Key=warehouse_123" \
  --header "Content-Type=application/json" \
  --body '{"source": "shopify", "target": "warehouse"}'

# Export configuration for backup
tempo export > backup.json
```

## Configuration

Jobs are stored in `~/.tempo/jobs.json` by default. You can customize the data directory by modifying the storage configuration.

## Development

### Project Structure

```
tempo/
├── cmd/
│   ├── cli/           # CLI application
│   │   ├── main.go    # CLI entry point
│   │   └── commands/  # CLI commands
│   └── main.go        # Original test application
├── internal/
│   ├── service/       # Scheduler service
│   ├── storage/       # Job storage
│   └── types/         # Data types
└── README.md
```

### Building

```bash
# Build CLI
go build -o tempo cmd/cli/main.go

# Build original test app
go build -o tempo-test cmd/main.go
```

## License

MIT License - see LICENSE file for details.