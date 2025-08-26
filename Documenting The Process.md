# Documenting Tempo

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

## Run prisma generate

```bash
go run github.com/steebchen/prisma-client-go generate
```

## Database synchronisation

```bash
go run github.com/steebchen/prisma-client-go generate
```

## Something I learned

Functions and Methods are different in Go!
Functions - Those that are not associated with a type
Methods - Those that are associated with a type
