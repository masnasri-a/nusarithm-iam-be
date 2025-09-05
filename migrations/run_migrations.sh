#!/bin/bash

set -e

echo "Running database migrations..."

# Database connection details (you can modify these or use environment variables)
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"5432"}
DB_NAME=${DB_NAME:-"nusarithm_iam"}
DB_USER=${DB_USER:-"postgres"}
DB_PASSWORD=${DB_PASSWORD:-""}

# Build connection string
if [ -n "$DB_PASSWORD" ]; then
    CONN_STR="postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
else
    CONN_STR="postgresql://$DB_USER@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
fi

echo "Connecting to database: $DB_NAME on $DB_HOST:$DB_PORT"

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo "Error: psql command not found. Please install PostgreSQL client."
    exit 1
fi

# Run migrations in order
for migration_file in $(ls migrations/*.sql | sort); do
    echo "Running migration: $migration_file"
    psql "$CONN_STR" -f "$migration_file"
    if [ $? -eq 0 ]; then
        echo "✓ Migration $migration_file completed successfully"
    else
        echo "✗ Migration $migration_file failed"
        exit 1
    fi
done

echo "All migrations completed successfully!"
