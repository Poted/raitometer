#!/bin/sh

echo "Waiting for postgres..."
while ! nc -z raitometer_db 5432; do
  sleep 0.1
done
echo "PostgreSQL started"

echo "Running database migrations..."
/app/migrate -database "$DATABASE_URL" -path /app/internal/migrations up

echo "Starting the application..."
/app/main