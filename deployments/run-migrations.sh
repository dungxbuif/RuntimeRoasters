#!/bin/bash
set -e

echo "Starting Database Migrations..."

# Setup connection params
export PGHOST=${PG_DB_HOST:-postgres-service}
export PGPORT=${PG_DB_PORT:-5432}
export PGUSER=${PG_USER:-dev}
export PGPASSWORD=${PG_PASS:-password}

# Function to run migrations for a specific database
run_migration() {
  local db_name=$1
  local migration_dir=$2

  echo "Migrating $db_name..."
  for sql_file in $(ls /migrations/$migration_dir/*.up.sql | sort); do
    echo "  -> Running $sql_file"
    psql -d "$db_name" -f "$sql_file"
  done
}

# Run all
run_migration "auth_db" "auth"
run_migration "logistics_db" "logistics"
run_migration "retail_db" "retail"
run_migration "warehouse_db" "warehouse"
run_migration "farm_db" "farm"
run_migration "payment_db" "payment"
run_migration "audit_db" "audit"
run_migration "trace_db" "trace"

echo "All migrations completed successfully!"
