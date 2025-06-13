.PHONY: all test clean imports fmt lint

# Default target
all: test

# Run tests
test:
	go test -v ./...

# Run goimports to format and organize imports
imports:
	@echo "Running goimports..."
	@goimports -w .
	@echo "Done."

# Run gofmt
fmt:
	@echo "Running gofmt..."
	@gofmt -w -s .
	@echo "Done."

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run
	@echo "Done."

# Clean up
clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f coverage.out
	@echo "Done."
