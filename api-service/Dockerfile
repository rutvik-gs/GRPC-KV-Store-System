# Build stage
FROM golang:1.25-bookworm AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory to the project root in the container.
# This makes it easier to reference both api-service and schemas modules.
WORKDIR /app

# Copy the schemas module files
COPY schemas/go.mod ./schemas/

# Copy api-service module files
COPY api-service/go.mod api-service/go.sum ./api-service/

# Download dependencies. This will respect the replace directive because
# the schemas module exists at the expected path.
RUN cd api-service && go mod download

# Copy source code for the api-service
COPY api-service/ ./api-service/

# Copy schemas source code
COPY schemas/ ./schemas/

# Build the binary inside the api-service directory
RUN cd api-service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a \
    -o /app/api-server ./main/main.go

# Runtime stage - use distroless for minimal image
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary and OpenAPI spec
COPY --from=builder /app/api-server /api-server
COPY --from=builder /app/schemas/rest/openapi.yaml /schemas/rest/openapi.yaml

# Use non-root user
USER nonroot:nonroot

# Expose REST API port
EXPOSE 8080

# Run the server
ENTRYPOINT ["/api-server", "--spec", "/schemas/rest/openapi.yaml"]
