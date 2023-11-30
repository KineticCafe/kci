BIN_NAME := "kci"

lint:
  @echo "Running linters..."
  gofmt -l .
  golangci-lint ./...

