## Name of the project
PROJECT_NAME = tempo

## Run the project
run:
	go run cmd/main.go

## Build the project
build:
	go build -o $(PROJECT_NAME) cmd/main.go


## Clean the project
clean:
	rm -f $(PROJECT_NAME)