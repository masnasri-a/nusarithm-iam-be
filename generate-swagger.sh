#!/bin/bash

set -e  # Exit on any error

echo "Starting Swagger documentation generation..."

# Check if swag command is available
if ! command -v swag &> /dev/null; then
    echo "Error: swag command not found. Please install it first."
    echo "Run: go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
fi

# Generate Swagger documentation
echo "Generating Swagger documentation..."
swag init -g main.go

# Check if docs.go was created
if [ ! -f "docs/docs.go" ]; then
    echo "Error: docs/docs.go was not generated."
    exit 1
fi

# Clean up the generated documentation using Go program
echo "Cleaning up generated documentation..."
go run cmd/clean-docs.go

echo "Swagger documentation generated and cleaned successfully!"
echo "You can access the documentation at: http://localhost:8080/swagger/index.html"
