## Setup

```bash
# Create project
mkdir tempo && cd tempo
go mod init tempo
```

```bash
# Run Redis (For AsynQ) and Postgres
# Run Redis (for Asynq)
docker run -d -p 6379:6379 redis:alpine

# Run Postgres
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:15
```

## Create the main.go file