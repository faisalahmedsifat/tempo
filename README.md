# Tempo

A modern, reliable webhook scheduler that eliminates the need for managing cron job infrastructure. Schedule any HTTP endpoint to run on a recurring basis with built-in monitoring, retries, and alerts.

## Overview

Tempo provides a simple API and dashboard for scheduling webhooks (HTTP requests) to run at specified intervals. Instead of managing servers, crontabs, or complex orchestration tools, developers can schedule jobs in seconds and monitor their execution through a clean interface.

### Key Features

- **Simple Scheduling** - Natural language scheduling or cron expressions
- **Reliable Execution** - 99.99% scheduling accuracy with automatic retries
- **Complete Visibility** - Dashboard showing execution history, logs, and metrics
- **Flexible Webhooks** - Support for any HTTP method, headers, and request bodies
- **Smart Alerts** - Get notified only when things actually need attention
- **Developer Friendly** - REST API, CLI, and SDKs for popular languages

## Quick Start

### Installation

```bash
# Using Docker
docker-compose up

# Or run locally
go mod download
go run github.com/steebchen/prisma-client-go generate
npx prisma migrate dev
go run cmd/server/main.go
```

### Basic Usage

#### Via CLI

```bash
# Install CLI
npm install -g @tempo/cli

# Schedule a webhook
tempo schedule "daily-backup" \
  --url "https://api.example.com/backup" \
  --schedule "every day at 2am"

# List all jobs
tempo list
# View execution history
tempo logs daily-backup
```

#### Via API

```bash
# Create a job
curl -X POST https://api.tempo.dev/v1/jobs \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "sync-data",
    "url": "https://api.example.com/sync",
    "schedule": "every hour",
    "method": "POST",
    "headers": {
      "Authorization": "Bearer token"
    }
  }'
```

#### Via SDK

```javascript
// JavaScript/TypeScript
import { Tempo } from '@tempo/sdk';

const tempo = new Tempo({ apiKey: 'tmp_xxx' });

await tempo.jobs.create({
  name: 'daily-report',
  url: 'https://api.example.com/reports/daily',
  schedule: 'every day at 9am',
  retries: 3,
  timeout: '5m'
});
```

## Architecture

Tempo is built with a clean, domain-driven architecture:

```
tempo/
├── cmd/server/          # Application entry point
├── internal/
│   ├── domain/          # Core business entities
│   ├── service/         # Business logic
│   ├── repository/      # Data access layer
│   ├── handler/         # HTTP handlers
│   └── infrastructure/  # External services
├── prisma/              # Database schema
└── pkg/                 # Shared utilities
```

### Technology Stack

- **Backend**: Go with Fiber framework
- **Database**: PostgreSQL with Prisma ORM
- **Queue**: Redis with Asynq
- **Scheduling**: Cron expression parser with distributed locking
- **API**: RESTful with OpenAPI specification

## API Reference

### Authentication

All API requests require an API key in the Authorization header:

```bash
Authorization: Bearer YOUR_API_KEY
```

### Endpoints

#### Create Job
`POST /api/v1/jobs`

```json
{
  "name": "job-name",
  "url": "https://webhook.url",
  "schedule": "every hour",
  "method": "POST",
  "headers": {},
  "body": "{}",
  "max_retries": 3,
  "timeout_seconds": 30
}
```

#### List Jobs
`GET /api/v1/jobs?limit=20&offset=0`

#### Get Job
`GET /api/v1/jobs/:id`

#### Update Job
`PUT /api/v1/jobs/:id`

#### Delete Job
`DELETE /api/v1/jobs/:id`

#### Run Job Immediately
`POST /api/v1/jobs/:id/run`

#### Get Execution History
`GET /api/v1/jobs/:id/executions`

### Schedule Formats

Tempo accepts both natural language and cron expressions:

#### Natural Language
- `every minute`
- `every 5 minutes`
- `every hour`
- `every day at 9am`
- `every monday at 2pm`
- `every first friday of month`

#### Cron Expressions
- `* * * * *` - Every minute
- `0 * * * *` - Every hour
- `0 9 * * 1-5` - Weekdays at 9am
- `0 0 1 * *` - First day of month

## Development

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Node.js 18+ (for Prisma CLI)

### Setup

1. Clone the repository
```bash
git clone https://github.com/faisalahmedsiifat/tempo.git
cd tempo
```

2. Install dependencies
```bash
go mod download
npm install
```

3. Set up environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run database migrations
```bash
npx prisma migrate dev
```

5. Generate Prisma client
```bash
go run github.com/steebchen/prisma-client-go generate
```

6. Start the server
```bash
go run cmd/server/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/service/...
```

### Building

```bash
# Build binary
go build -o bin/tempo cmd/server/main.go

# Build Docker image
docker build -t tempo:latest .

# Build with version info
go build -ldflags "-X main.Version=1.0.0" -o bin/tempo cmd/server/main.go
```

## Deployment

### Using Docker Compose

```bash
docker-compose up -d
```

### Using Kubernetes

```bash
kubectl apply -f k8s/
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://localhost/tempo` |
| `REDIS_URL` | Redis connection string | `localhost:6379` |
| `SECRET_KEY` | JWT signing key | - |
| `ENVIRONMENT` | Environment (development/production) | `development` |

## Configuration

### Rate Limits

Default rate limits per API key:
- 1000 requests per hour
- 100 concurrent jobs
- 10,000 executions per day

### Retention

- Execution logs: 30 days (free), 90 days (paid)
- Metrics: 90 days
- Failed job details: 7 days

## Monitoring

Tempo exposes Prometheus metrics at `/metrics`:

- `tempo_jobs_total` - Total number of jobs
- `tempo_executions_total` - Total executions by status
- `tempo_execution_duration_seconds` - Execution duration histogram
- `tempo_webhook_response_time_seconds` - Webhook response times
- `tempo_scheduler_lag_seconds` - Scheduling delay

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Run `go fmt` before committing
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Add tests for new features
- Update documentation as needed

## Roadmap

- [ ] Workflow orchestration (job dependencies)
- [ ] Event-driven triggers
- [ ] Webhook request/response transformations
- [ ] Team collaboration features
- [ ] Advanced scheduling (timezone-aware, holidays)
- [ ] Webhook authentication methods (OAuth, mTLS)
- [ ] Bulk operations API
- [ ] Terraform provider
- [ ] Grafana dashboard templates

## Support

- **Documentation**: [docs.tempo.dev](https://docs.tempo.dev)
- **Issues**: [GitHub Issues](https://github.com/faisalahmedsiifat/tempo/issues)
- **Discord**: [Join our community](https://discord.gg/tempo)
- **Email**: support@tempo.dev

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Asynq](https://github.com/hibiken/asynq) for reliable task processing
- [Fiber](https://github.com/gofiber/fiber) for fast HTTP handling
- [Prisma](https://www.prisma.io/) for type-safe database access
- [Robfig/cron](https://github.com/robfig/cron) for cron expression parsing

---

Built with precision scheduling in mind. Never miss a webhook again.