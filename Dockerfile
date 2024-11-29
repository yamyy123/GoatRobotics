# Use the official Golang image as the base image (Go 1.21)
FROM golang:1.21-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download all the dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o myapp .

# Start a new stage from a smaller base image (to minimize the size of the final image)
FROM alpine:latest  

# Install the required CA certificates (needed to handle SSL/TLS)
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the pre-built binary from the builder image
COPY --from=builder /app/myapp .

# Copy the config.json file into the container
COPY config.json ./config.json

# Expose the port the app runs on (default Go port is 8080)
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
