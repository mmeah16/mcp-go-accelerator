
# Use a Go v1.24.3 base image
FROM golang:1.24.2

# Create a directory for the application
WORKDIR /app

# Copy the Go module files to app
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application with entry point ./cmd/server to point to main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /mcp-go-accelerator ./cmd/server

# Expose the port the app runs on
EXPOSE 8080

CMD ["/mcp-go-accelerator"]