#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Start the Next.js server in the background, explicitly setting host and port
echo "Starting Next.js server on 0.0.0.0:3000..."
(cd frontend/.next/standalone && HOSTNAME=0.0.0.0 PORT=3000 node server.js &)

# Wait a few seconds to ensure the Next.js server has time to start
sleep 5

# Start the Go server in the foreground
# It will act as the main entry point and reverse proxy
echo "Starting Go server..."
./backend/app
