#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Start the Next.js server in the background
echo "Starting Next.js server..."
(cd frontend/.next/standalone && node server.js > ../../../frontend_server.log 2>&1 &)

# Wait a few seconds to ensure the Next.js server has time to start
sleep 5

# Start the Go server in the foreground
# It will act as the main entry point and reverse proxy
echo "Starting Go server..."
./backend/app
