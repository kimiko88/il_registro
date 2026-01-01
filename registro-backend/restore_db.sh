#!/bin/bash
set -e

echo "Starting Database Restoration..."

# 1. Run Migrations
echo "Running Migrations..."
# Get all .sql files in migrations folder, sorted
for f in $(ls migrations/*.sql | sort); do
    # Skip the 022 seed file as we use the Go seeder
    if [[ "$f" != *"022_seed_data.sql"* ]]; then
        echo "Applying $f..."
        cat "$f" | sudo docker exec -i registro-backend-db-1 psql -U user -d registro
    fi
done

# 2. Run Go Seeder
echo "Running Go Seeder..."
# Ensure we can connect to localhost DB. The go seeder expects DB at localhost:5432
if go run cmd/seed/main.go; then
    echo "Seeder finished successfully."
else
    echo "Seeder failed. Attempting to run with internal container networking?"
    # Fallback or checks if port is exposed? 
    # Docker compose maps 5432:5432, so localhost should work.
    # If it failed before, maybe it was because of blocked sudo or network delay.
    exit 1
fi

# 3. Create Superadmin
echo "Creating Superadmin..."
cat create_superadmin.sql | sudo docker exec -i registro-backend-db-1 psql -U user -d registro

echo "Database successfully restored!"
