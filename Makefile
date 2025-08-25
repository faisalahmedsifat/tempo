## Name of the project
PROJECT_NAME = tempo

## Go build flags
GO_BUILD_FLAGS = -v -o $(PROJECT_NAME)

## Run the project
run:
	go run main.go

## Build the project
build:
	go build $(GO_BUILD_FLAGS)


## Clean the project
clean:
	rm -f $(PROJECT_NAME)