# ==============================================================================
# STAGE 1: Build the Go application static binary
# ==============================================================================
FROM golang:1.25-alpine AS builder

# Install system dependencies needed for compiling
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory inside the container
WORKDIR /app

# Copy dependency definition files to leverage Docker layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application static binary:
# - CGO_ENABLED=0 disables dynamic linking, ensuring a completely static binary
# - GOOS=linux compiles the binary for Linux
# - -ldflags="-w -s" strips debugging symbols and DWARF tables to minimize binary size
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/api \
    ./cmd/api/main.go

# Create a non-root system user for security best practices
# Running as non-root prevents privilege escalation attacks inside the container
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 10001 \
    appuser

# ==============================================================================
# STAGE 2: Create the final minimal and secure runtime image
# ==============================================================================
FROM alpine:latest

# Install CA certificates and Timezone database
# - ca-certificates is critical for outbound secure HTTPS calls (e.g. Cloudinary uploads)
# - tzdata ensures that dates, times, and logs align correctly with local timezone settings
RUN apk --no-cache add ca-certificates tzdata

# Copy user/group details from builder to run as non-root
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Set the working directory for the application
WORKDIR /app

# Copy the compiled static binary from the builder stage
# Set ownership to our non-root appuser
COPY --from=builder --chown=appuser:appuser /app/api /app/api

# Copy .env.example as a template reference
COPY --from=builder --chown=appuser:appuser /app/.env.example /app/.env.example

# Set standard environment defaults (can be overridden at runtime)
ENV PORT=3010
ENV TZ=Asia/Jakarta

# Switch to the secure non-root user
USER appuser

# Expose the application port
EXPOSE 3010

# Define the entrypoint to run the Go application binary
ENTRYPOINT ["/app/api"]
