.PHONY: build release test clean

# Build the CLI tool
build:
	@echo "Building Sato Framework CLI..."
	@go build -o bin/sato cli/sato.go

# Create a new release
release:
	@echo "Creating new release..."
	@git tag v0.1.0
	@git push origin v0.1.0

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@go clean

# Install the CLI tool
install:
	@echo "Installing Sato Framework CLI..."
	@go install ./cli/sato

# Generate documentation
docs:
	@echo "Generating documentation..."
	@go doc -all ./... > docs.md

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@go vet ./... 