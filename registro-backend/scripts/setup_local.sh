#!/bin/bash
echo "Starting Local Database..."
docker-compose -f ../docker/supabase/docker-compose.yml up -d

echo "Waiting for DB to be healthy..."
until [ "`docker inspect -f {{.State.Health.Status}} $(docker-compose -f ../docker/supabase/docker-compose.yml ps -q db)`" == "healthy" ]; do
    sleep 2
    echo -n "."
done
echo "DB is ready!"

echo "Seeding Data..."
docker exec -i $(docker-compose -f ../docker/supabase/docker-compose.yml ps -q db) psql -U postgres -d postgres < ../seed.sql
