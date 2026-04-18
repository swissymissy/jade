#!/bin/sh 

set -e

# exit if command fails
echo "Running migration..."
goose -dir /app/sql/schema sqlite3 "$DB_PATH" up

# hands off to main app
echo "Starting server..."
exec /app/jade
