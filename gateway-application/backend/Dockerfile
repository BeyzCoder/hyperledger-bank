# Stage 1: Build the Go binary
FROM golang:1.23-alpine as builder

# Install CompileDaemon for live reloading
RUN go install github.com/githubnemo/CompileDaemon@latest

# Set the working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Expose the Gin application port
EXPOSE 8080

# Command to run CompileDaemon for live reloading
CMD ["CompileDaemon", "--build=go build -o /app/main .", "--command=./main", "--directory=/app", "--polling"]