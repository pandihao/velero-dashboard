# Velero API Server

REST API wrapper for Velero backup and restore operations.

## API Documentation

After starting the server, access the Swagger documentation at:

```
http://localhost:8080/docs
```

## Quick Start

### Start the server

```bash
# Local development
./velero-api-server --kubeconfig ~/.kube/config --port 8080

# In Kubernetes cluster (auto-detects kubeconfig)
./velero-api-server --port 8080

# Skip TLS verification (for self-signed certs)
./velero-api-server --kubeconfig ~/.kube/config --insecure-skip-tls
```

### Command-line options

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | 8080 | HTTP server listen port |
| `--kubeconfig` | "" | Path to kubeconfig file (empty for in-cluster) |
| `--namespace` | velero | Velero installation namespace |
| `--insecure-skip-tls` | false | Skip TLS certificate verification |

### Environment variables

- `VELERO_NAMESPACE`: Override the Velero namespace
- `KUBECONFIG`: Path to kubeconfig file

## API Examples

### Create a backup

```bash
curl -X POST http://localhost:8080/api/v1/backups \
  -H "Content-Type: application/json" \
  -d '{
    "name": "backup-default-all",
    "includedNamespaces": ["default"],
    "ttl": "720h"
  }'
```

### List backups

```bash
curl http://localhost:8080/api/v1/backups
```

### Get backup details

```bash
curl http://localhost:8080/api/v1/backups/backup-default-all
```

### Create a restore

```bash
curl -X POST http://localhost:8080/api/v1/restores \
  -H "Content-Type: application/json" \
  -d '{
    "name": "restore-default-20260623",
    "backupName": "backup-default-all"
  }'
```

### Create a scheduled backup

```bash
curl -X POST http://localhost:8080/api/v1/schedules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "daily-backup",
    "schedule": "0 2 * * *",
    "includedNamespaces": ["default"],
    "ttl": "720h"
  }'
```

## Build

### For current platform

```bash
go build -o velero-api-server cmd/server/main.go
```

### Cross-compile for Linux

```bash
GOOS=linux GOARCH=amd64 go build -o velero-api-server-linux cmd/server/main.go
```

### Optimized build (smaller binary)

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o velero-api-server-linux cmd/server/main.go
```

## Project Structure

```
velero-api-server/
├── cmd/server/          # Application entry point
│   └── main.go
├── internal/
│   ├── handler/         # HTTP handlers
│   │   ├── backup.go
│   │   ├── restore.go
│   │   ├── schedule.go
│   │   ├── bsl.go
│   │   └── routes.go
│   ├── model/           # Request/response models
│   │   └── types.go
│   └── service/         # Business logic
│       └── velero.go
├── pkg/k8s/            # Kubernetes client
│   └── client.go
└── docs/               # API documentation
    ├── swagger.yaml
    └── swagger.html
```

## Requirements

- Go 1.22+
- Access to a Kubernetes cluster with Velero installed
- kubeconfig file (for out-of-cluster mode)
