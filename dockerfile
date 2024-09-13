# Use the official Go image as the base image
FROM golang:1.23-bullseye AS build-stage

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download and cache the dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./main.go

# Use a smaller base image for the final stage
FROM debian:bullseye-slim AS build-release-stage

# Install curl for healthchecks
RUN apt-get update && apt-get install -y curl

# Set the Current Working Directory inside the container
WORKDIR /

# Copy the compiled binary from the build stage
COPY --from=build-stage /main /main

# Run the compiled file
CMD ["/main"]