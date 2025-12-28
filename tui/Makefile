# Project: 20i Stack Manager TUI
# File: Makefile
# Purpose: Build, install, clean, and test automation for the TUI.
# Version: 0.1.0
# Updated: 2025-12-28
.PHONY: build install clean test test-coverage

# Build the TUI binary
build:
	go build -o bin/20i-stack-manager main.go

# Install the binary to GOPATH/bin
install:
	go install

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run all tests
test:
	go test -v ./...

# Run tests with coverage report
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
