# Start from the official Go image
FROM golang:1.23-alpine

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o todo-app ./cmd/server

# Expose port
EXPOSE 8080

# Run the application
CMD ["./todo-app"]
