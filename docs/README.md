# Mock Netbox Service

A mock Netbox service for E2E testing of network documentation and device management applications without requiring a real Netbox instance. Provides HTTP/REST endpoints for managing devices, sites, and interfaces.

## Overview

This service emulates the Netbox API for testing purposes. It supports:

- **Devices**: Physical network devices (switches, routers, servers) lifecycle management (create, list, get, update, delete)
- **Sites**: Physical locations/data centers (create, list, get, delete)
- **Interfaces**: Network interfaces on devices (create, list, get, delete)

All data is stored in-memory and reset on service restart.

## Endpoints

### Health Checks

- `GET /health` - Liveness probe (always returns 200)
- `GET /ready` - Readiness probe (returns 200 when ready)

### Devices

- `GET /api/dcim/devices` - List all devices
- `POST /api/dcim/devices` - Create a new device
  - Body: `{"name": string, "site": int, "type": string, "status": string, "serial": string, "asset_tag": string}`
  - Required: `name`, `site`, `type`
  - Default status: "active"
- `GET /api/dcim/devices/{id}` - Get device details
- `PUT /api/dcim/devices/{id}` - Update device
  - Body: Partial fields allowed
- `DELETE /api/dcim/devices/{id}` - Delete device

### Interfaces

- `GET /api/dcim/interfaces` - List all interfaces
- `POST /api/dcim/interfaces` - Create a new interface
  - Body: `{"device": int, "name": string, "type": string, "status": string, "enabled": bool, "mtu": int, "mode": string, "mac_address": string}`
  - Required: `device`, `name`, `type`
  - Default status: "active"
  - Default enabled: true
- `GET /api/dcim/interfaces/{id}` - Get interface details
- `DELETE /api/dcim/interfaces/{id}` - Delete interface

### Sites

- `GET /api/sites` - List all sites
- `POST /api/sites` - Create a new site
  - Body: `{"name": string, "slug": string, "region": string, "facility": string, "asn": int, "time_zone": string, "description": string}`
  - Required: `name`
  - Auto-generated slug if not provided
- `GET /api/sites/{id}` - Get site details
- `DELETE /api/sites/{id}` - Delete site

## Configuration

All configuration is via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | 8080 | HTTP server port |
| `LOG_LEVEL` | info | Logging level (debug, info, warn, error) |
| `GIN_MODE` | release | Gin framework mode (debug, release) |

## Example Usage

### Local Development

```bash
# Install dependencies
go mod tidy

# Run lint checks
make lint

# Run tests
make test

# Start server
make run

# In another terminal, test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/dcim/devices

# Create a site
curl -X POST http://localhost:8080/api/sites \
  -H "Content-Type: application/json" \
  -d '{"name": "us-west-1"}'

# Create a device
curl -X POST http://localhost:8080/api/dcim/devices \
  -H "Content-Type: application/json" \
  -d '{"name": "switch-01", "site": 1, "type": "switch"}'
```

### Docker

```bash
# Build image
make docker-build

# Run container
make docker-run

# Test in another terminal
curl http://localhost:8080/health
```

### Kubernetes

```bash
# Deploy to cluster
kubectl apply -f k8s/manifest.yaml

# Check status
kubectl get pods -n mock-services
kubectl logs -n mock-services -l app=mock-netbox

# Port forward for local testing
kubectl port-forward -n mock-services svc/mock-netbox 8080:8080
curl http://localhost:8080/health
```

## Response Format

All responses are JSON. Successful responses follow this format:

```json
{
  "count": 5,
  "results": [...]
}
```

For single resource responses:

```json
{
  "id": 1,
  "name": "device-01",
  ...
}
```

Error responses:

```json
{
  "error": "description"
}
```

## Request Tracing

All requests include a unique `X-Request-ID` header for tracing:

```bash
curl -H "X-Request-ID: my-trace-id" http://localhost:8080/api/dcim/devices
```

If not provided, a UUID is generated automatically.

## Logging

All logs are structured JSON with the following fields:

- `timestamp` - RFC3339 timestamp
- `level` - Log level (INFO, WARN, ERROR)
- `msg` - Log message
- `request_id` - Request trace ID
- Additional context fields depending on operation

## Testing

```bash
# Run all tests
go test -v -race ./...

# With coverage report
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Run specific test
go test -v -run TestDeviceHandler ./internal/handler
```

## Building

### Local Build

```bash
go build -o mock-netbox ./cmd/main
./mock-netbox
```

### Docker Build

The Dockerfile uses multi-stage builds:

1. **Builder stage**: Go 1.26.4, downloads dependencies, runs golangci-lint and tests
2. **Final stage**: scratch image with only the binary and CA certificates

Target image size: ≤50MB

```bash
docker build -t mock-netbox:latest .
```

## Architecture

```
cmd/main/
  └── main.go              # Entry point, Gin setup, route registration

internal/
  ├── handler/
  │   ├── health.go        # Liveness/readiness probes
  │   ├── device.go        # Device management
  │   ├── interface.go     # Interface management
  │   ├── site.go          # Site management
  │   ├── init.go          # Shared store initialization
  │   └── handlers_test.go # Handler tests
  ├── middleware/
  │   └── middleware.go    # Gin middleware stack
  └── store/
      └── store.go         # In-memory data storage
```

## Standards

This service follows the mock-servers standards:

- **Language**: Go 1.26.4
- **Framework**: Gin for HTTP
- **Logging**: stdlib `log/slog` with JSON output
- **Linting**: golangci-lint 2.12.2 (no exceptions)
- **Testing**: ≥70% coverage on handler logic
- **Kubernetes**: Deployment with health probes and resource limits
- **CI/CD**: GitHub Actions with lint → test → build pipeline

See [STANDARDS.md](../docs/STANDARDS.md) in the mock-servers repo for full guidelines.

## Contributing

Follow the [Contributing Guidelines](../docs/CONTRIBUTING.md) in the mock-servers repo.

## License

Mocks for testing purposes only.
