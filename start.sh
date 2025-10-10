#!/bin/bash

# --- Diagnostic Steps ---
echo "--- Running Diagnostic start.sh ---"

echo "[DEBUG] Current working directory:"
pwd

echo "[DEBUG] Listing files in root directory:"
ls -lA

echo "[DEBUG] Listing files in frontend directory:"
ls -lA frontend

echo "[DEBUG] Listing files in frontend/.next/standalone directory:"
ls -lA frontend/.next/standalone

# --- Original Logic with enhanced logging ---
# We are not using 'set -e' here to allow the script to continue and report diagnostics
# even if the server fails to start immediately.

# Define log file path at the root
LOG_FILE="/app/frontend_server.log"

# Start the Next.js server in the background
echo "Attempting to start Next.js server... Logging to $LOG_FILE"
# The cd is necessary to run node server.js from the correct directory
(cd /app/frontend/.next/standalone && node server.js > "$LOG_FILE" 2>&1 &)


# Wait a few seconds to ensure the Next.js server has time to start or fail
echo "Waiting for 5 seconds..."
sleep 5

# --- Diagnostic Check ---
echo "[DEBUG] Checking for log file at $LOG_FILE..."
if [ -f "$LOG_FILE" ]; then
    echo "[DEBUG] Log file '$LOG_FILE' found. Contents:"
    cat "$LOG_FILE"
else
    echo "[DEBUG] Log file '$LOG_FILE' was NOT created."
fi
echo "--- End of Diagnostic Checks ---"


# Start the Go server in the foreground
echo "Starting Go server..."
./backend/app
