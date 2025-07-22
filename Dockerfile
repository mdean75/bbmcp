# Build stage
FROM golang:1.23-bullseye AS builder

ENV CA_DIR=/usr/local/share/ca-certificates \
    SSL_DIR=/etc/ssl/certs/

COPY certificates/* "${CA_DIR}/mc/"

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && update-ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

## Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o bbmcp .

# Final stage
FROM debian:bullseye-slim

ENV CA_DIR=/usr/local/share/ca-certificates/

COPY --from=builder /usr/local/share/ca-certificates/mc/ /usr/local/share/ca-certificates/

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*


## Create non-root user
RUN addgroup --gid 1001 --system bbmcp && \
    adduser --uid 1001 --system --ingroup bbmcp bbmcp

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/bbmcp .

# Change ownership to non-root user
RUN chown bbmcp:bbmcp /app/bbmcp

# Switch to non-root user
USER bbmcp

# Set environment variables with defaults
ENV BITBUCKET_BASE_URL=""
ENV BITBUCKET_USERNAME=""
ENV BITBUCKET_PASSWORD=""
ENV BITBUCKET_TOKEN=""
ENV BITBUCKET_DEFAULT_PROJECT_KEY=""

# Run the MCP server
CMD ["./bbmcp"]
