.PHONY: test bench bench-baseline bench-compare bench-ci clean
.DEFAULT_GOAL := test

# Variables
BENCHMARK_BASELINE := benchmark.baseline
BENCHMARK_CURRENT := benchmark.current
THRESHOLD := 20

# Run all tests
test:
	go test -v -race ./...

# Run all tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run benchmarks and save to current file
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem -run=^$$ ./... > $(BENCHMARK_CURRENT)
	@echo "Benchmarks saved to $(BENCHMARK_CURRENT)"

# Generate baseline benchmark file
bench-baseline:
	@echo "Generating benchmark baseline..."
	go test -bench=. -benchmem -run=^$$ ./... > $(BENCHMARK_BASELINE)
	@echo "Baseline saved to $(BENCHMARK_BASELINE)"

# Compare current benchmarks with baseline
bench-compare: bench
	@if [ ! -f $(BENCHMARK_BASELINE) ]; then \
		echo "Baseline file not found. Run 'make bench-baseline' first."; \
		exit 1; \
	fi
	@chmod +x ./scripts/compare-benchmarks.sh
	@echo "Comparing benchmarks with threshold of $(THRESHOLD)%..."
	@./scripts/compare-benchmarks.sh $(BENCHMARK_BASELINE) $(BENCHMARK_CURRENT) $(THRESHOLD)

# CI benchmark check (used in GitHub Actions)
bench-ci:
	@echo "Running CI benchmark check..."
	@if [ ! -f $(BENCHMARK_BASELINE) ]; then \
		echo "No baseline found, skipping benchmark regression check"; \
		exit 0; \
	fi
	@$(MAKE) bench-compare

# Clean benchmark files
clean:
	@echo "Cleaning benchmark files..."
	@rm -f $(BENCHMARK_BASELINE) $(BENCHMARK_CURRENT) coverage.txt

# Build the project
build:
	go build -v ./...

# Download dependencies
deps:
	go mod download
	go mod tidy