#!/bin/bash

# Replace with your actual project name
PROJECT_NAME="gohtwind"

# Name of the Docker image and container
IMAGE_NAME="${PROJECT_NAME}-dev-image"
CONTAINER_NAME="${PROJECT_NAME}-dev-container"

# Build the Docker image using the development Dockerfile
docker build -f Dockerfile.dev -t "$IMAGE_NAME" .

# Check if the container is already running, and if so, stop it
if docker ps -q -f name="$CONTAINER_NAME" | grep -q .; then
    docker stop "$CONTAINER_NAME"
fi

# Remove the container if it exists
if docker ps -aq -f name="$CONTAINER_NAME" | grep -q .; then
    docker rm "$CONTAINER_NAME"
fi

# Run the new container from the built image
docker run -d \
    --name "$CONTAINER_NAME" \
    -p 8080:8080 \
    -p 40000:40000 \
    -v "$(pwd)":/app \
    "$IMAGE_NAME"

# Print the logs for the running container
docker logs -f "$CONTAINER_NAME"
