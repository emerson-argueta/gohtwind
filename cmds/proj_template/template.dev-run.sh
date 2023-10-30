#!/bin/bash

# Function to start the Go watcher
start_go_watcher() {
    # Install debugging tools
    go install github.com/go-delve/delve/cmd/dlv@latest
    # Install air for live-reloading
    go install github.com/cosmtrek/air@latest
    # start air
    air -c .air.toml
}

# Function to start the CSS watcher
start_css_watcher() {
    # Install frontend dependencies
    chmod +x ./frontend/bin/tailwindcss
    ./frontend/bin/tailwindcss -i ./frontend/static/css/main.css -o ./frontend/static/css/output.css --watch
}

# Prepare Go modules
echo "Running go mod tidy..."
go mod tidy

# Start the watchers in the background
start_go_watcher &
GO_WATCHER_PID=$!

start_css_watcher
CSS_WATCHER_PID=$!


# Function to stop background processes when this script is stopped
stop_watchers() {
    kill $GO_WATCHER_PID
    kill $CSS_WATCHER_PID
}

# Handle script termination
trap stop_watchers EXIT

# Wait indefinitely until the script is explicitly stopped
while true; do
    sleep 1
done
