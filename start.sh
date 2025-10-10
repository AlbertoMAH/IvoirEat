#!/bin/bash

# --- Diagnostic Step ---
echo "[DEBUG] Writing file list to /app/file_list.txt..."
ls -R /app/frontend > /app/file_list.txt
echo "[DEBUG] Finished writing file list."
# --- End of Diagnostic Step ---

set -e

# Start the Next.js server in the background, explicitly setting host and port
echo "Starting Next.js server on 0.0.0.0:3000..."
(cd /app/frontend/.next/standalone && HOSTNAME=0.0.0.0 PORT=3000 node server.js &)

# Wait a few seconds to ensure the Next.js server has time to start
sleep 5

# Start the Go server in the foreground
# It will act as the main entry point and reverse proxy
echo "Starting Go server..."
/app/backend/app
