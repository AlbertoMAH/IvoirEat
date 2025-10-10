#!/bin/bash

set -e

# Start the Next.js server in the background, using relative paths
echo "Starting Next.js server on 0.0.0.0:3000..."
(cd ./frontend/.next/standalone && HOSTNAME=0.0.0.0 PORT=3000 node server.js &)

# Wait a few seconds to ensure the Next.js server has time to start
sleep 5

# Start the Go server in the foreground, using a relative path
echo "Starting Go server..."
./backend/app
