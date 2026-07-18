#!/bin/bash
set -e
echo "Running Migrations..."
for f in $(ls migrations/*.sql | sort); do
    if [[ "$f" != *"022_seed_data.sql"* ]]; then
        echo "Applying $f..."
        docker exec -i registro-backend-db-1 psql -U user -d registro < "$f"
    fi
done
echo "Migrations finished."
