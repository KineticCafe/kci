set positional-arguments

BIN_NAME := "kci"

# run linteres on the project
lint:
  @echo "Running linters..."
  golangci-lint run

# run the application
run *args='':
  go run main.go "$@"

# build the application
build: lint
  go build

release:

clean:
  rm -rf {{BIN_NAME}}
