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
  go build \
  -ldflags "-X 'main.BuildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ')'" 

# install application locally
install: lint test
  go install \
  -ldflags "-X 'main.BuildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ')'" 

# Create and upload a release
release:
  @echo "I am releasing to S3"


clean:
  rm -rf {{BIN_NAME}}
