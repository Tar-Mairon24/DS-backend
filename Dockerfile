# Use the official Golang image to create a build environment
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code to the Working Directory inside the container
COPY . /app/

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o main main.go

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

RUN apk add --no-cache ca-certificates && chmod +x main && ls -l /app && sleep 10

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD [ "sh", "-c", "./main" ]