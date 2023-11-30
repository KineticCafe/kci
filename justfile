set positional-arguments

BIN_NAME := "kci"

default: 
  just --list

# run all tests
test:
  go test ./...

# run linters on the project
lint:
  @echo "Running linters..."
  golangci-lint run

# run the application
run *args='':
  go run main.go "$@"

# build the application
build: lint test
  go build

# install application locally
install: lint test
  go install

# Create and upload a release
release:
  @echo "I am releasing to S3"


clean:
  rm -rf {{BIN_NAME}}
