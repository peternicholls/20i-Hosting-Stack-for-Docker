# Top-level Makefile for the 20i stack repo
# This proxies common TUI targets so you can run them from the repository root
# Example: `make build` will run `make build` inside the `tui` directory

TUI_DIR := tui

.PHONY: build install clean test test-coverage help

build:
	$(MAKE) -C $(TUI_DIR) build

install:
	$(MAKE) -C $(TUI_DIR) install

clean:
	$(MAKE) -C $(TUI_DIR) clean

test:
	$(MAKE) -C $(TUI_DIR) test

test-coverage:
	$(MAKE) -C $(TUI_DIR) test-coverage

help:
	@echo "Top-level Makefile; proxies targets to $(TUI_DIR)"
	@echo "Usage: make <target>"
	@echo "Available targets: build install clean test test-coverage"