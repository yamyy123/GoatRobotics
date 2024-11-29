#!/bin/bash

# Set the image name (you can change this to anything you like)
IMAGE_NAME="myapp"

# Build the Docker image
echo "Building Docker image..."
docker build -t $IMAGE_NAME .

# Check if the image was built successfully
if [ $? -eq 0 ]; then
  echo "Docker image built successfully!"
else
  echo "Docker image build failed!"
  exit 1
fi

# Run the Docker container (in detached mode)
echo "Running the Docker container..."
docker run -d -p 8080:8080 --name $IMAGE_NAME-container $IMAGE_NAME

# Check if the container is running
if [ $? -eq 0 ]; then
  echo "Docker container is running!"
else
  echo "Failed to start the Docker container!"
  exit 1
fi

# Print the status of the running container
docker ps

echo "You can access the app at http://localhost:8080"
