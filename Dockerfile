# Build stage
FROM golang:1.23-alpine AS builder

# Install git for go modules
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bbcli .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S bbcli && \
    adduser -u 1001 -S bbcli -G bbcli

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/bbcli .

# Change ownership to non-root user
RUN chown bbcli:bbcli /app/bbcli

# Switch to non-root user
USER bbcli

# Expose port (though MCP uses STDIO, this is for health checks if needed)
EXPOSE 8080

# Set environment variables with defaults
ENV BITBUCKET_BASE_URL=""
ENV BITBUCKET_USERNAME=""
ENV BITBUCKET_PASSWORD=""
ENV BITBUCKET_TOKEN=""
ENV BITBUCKET_DEFAULT_PROJECT_KEY=""

# Run the MCP server
CMD ["./bbcli"]