# Build stage
FROM golang:1.25-bookworm AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory to the project root in the container.
# This makes it easier to reference both kvStore-service and schemas modules.
WORKDIR /app

# Copy the schemas module files
COPY schemas/go.mod ./schemas/

# Copy kvStore-service module files
COPY kvStore-service/go.mod kvStore-service/go.sum ./kvStore-service/

# Download dependencies. This will respect the replace directive because
# the schemas module exists at the expected path.
RUN cd kvStore-service && go mod download

# Copy source code for the kvStore-service
COPY kvStore-service/ ./kvStore-service/

# Copy schemas source code
COPY schemas/ ./schemas/

# Build the binary inside the kvStore-service directory
RUN cd kvStore-service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a \
    -o /app/kvstore-server ./main/main.go

# Runtime stage - use distroless for minimal image
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary
COPY --from=builder /app/kvstore-server /kvstore-server

# Use non-root user
USER nonroot:nonroot

# Expose gRPC port
EXPOSE 50051

# Run the server
ENTRYPOINT ["/kvstore-server"]
