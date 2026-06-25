# tyk-sre-assignment

This repository contains the boilerplate projects for the SRE role interview assignments.

### Project

Location: https://github.com/TykTechnologies/tyk-sre-assignment/tree/main/golang

A Kubernetes monitoring and health check application built in Go that provides HTTP endpoints for cluster health monitoring, deployment status reporting, and application readiness checks.

---

## Project Structure

### Modules

The application is organized into three main packages:

#### 1. **main** (`main.go`)
- Entry point for the application
- Establishes Kubernetes connection with flexible configuration fallback:
  1. Explicit kubeconfig file (if provided)
  2. In-cluster config (for Kubernetes deployments)
  3. Local kubeconfig fallback (for development)
- Initializes HTTP server with registered handlers
- Validates cluster connectivity on startup

#### 2. **handlers** (`handlers/`)
- `deployments.go`: Handles `/deployments/health` endpoint - queries cluster for deployment status
- `readiness.go`: Handles `/readyz` endpoint - checks Kubernetes API server health
- Encapsulates HTTP request/response logic and error handling

#### 3. **kubernetes** (`kubernetes/`)
- `deployments.go`: `GetDeploymentHealth()` - lists all deployments across all namespaces and compares desired vs. ready replicas
- `health.go`: `CheckAPIServer()` - verifies Kubernetes API server connectivity and measures latency
- Core Kubernetes API operations with timeout handling

---

## Endpoints

The application exposes three HTTP endpoints:

### 1. `/healthz` - Application Health Check
**Method:** GET

**Purpose:** Simple application liveness probe

**Response:**
```
200 OK
ok
```

**Example:**
```bash
curl http://localhost:8080/healthz
```

---

### 2. `/deployments/health` - Kubernetes Deployments Status
**Method:** GET

**Purpose:** Returns health status of all Kubernetes deployments across all namespaces

**Query Parameters:**
- `pretty=true` - Optional. Pretty-prints JSON response with indentation

**Response Format:**
```json
[
  {
    "namespace": "kube-system",
    "name": "coredns",
    "desired": 2,
    "ready": 2,
    "healthy": true
  },
  {
    "namespace": "default",
    "name": "my-app",
    "desired": 3,
    "ready": 2,
    "healthy": false
  }
]
```

**Response Fields:**
- `namespace` - Kubernetes namespace
- `name` - Deployment name
- `desired` - Desired number of replicas (from spec)
- `ready` - Current number of ready replicas
- `healthy` - Boolean indicating if desired == ready && desired > 0

**Examples:**
```bash
# Standard JSON response
curl http://localhost:8080/deployments/health

# Pretty-printed response
curl http://localhost:8080/deployments/health?pretty=true
```

**Error Response:**
```
500 Internal Server Error
error fetching deployments: <error details>
```

---

### 3. `/readyz` - Kubernetes API Server Readiness
**Method:** GET

**Purpose:** Checks Kubernetes API server connectivity and latency

**Response Format:**
```json
{
  "connected": true,
  "serverVersion": "v1.28.0",
  "latencyMs": 15,
  "error": ""
}
```

**Response Fields:**
- `connected` - Boolean indicating API server connectivity
- `serverVersion` - Kubernetes server version (only when connected)
- `latencyMs` - Round-trip latency to API server in milliseconds (only when connected)
- `error` - Error message (only when not connected)

**Examples:**
```bash
curl http://localhost:8080/readyz
```

**Example Success Response:**
```json
{
  "connected": true,
  "serverVersion": "v1.28.0",
  "latencyMs": 12
}
```

**Example Failure Response:**
```json
{
  "connected": false,
  "error": "unable to connect to Kubernetes API server"
}
```

---

## Workflow

### Application Startup

1. **Parse Flags:** Application accepts two optional flags:
   - `--kubeconfig` - Path to kubeconfig file (default: empty)
   - `--address` - HTTP server listen address (default: `:8080`)

2. **Build Kubernetes Configuration:** The application uses a three-tier fallback strategy:
   - If `--kubeconfig` is provided, use that explicit config
   - Otherwise, try in-cluster configuration (useful when running in Kubernetes)
   - Fall back to `~/.kube/config` for local development

3. **Connect to Kubernetes:** Establishes client connection and validates by fetching server version

4. **Start HTTP Server:** Registers handlers and begins listening for requests

### Request Handling

- **Deployment Health Requests:** When `/deployments/health` is called:
  1. Handler receives request and extracts query parameters
  2. Calls `kubernetes.GetDeploymentHealth()` with 5-second timeout
  3. Function lists all deployments across all namespaces
  4. Compares desired vs. ready replicas for each deployment
  5. Returns JSON array of deployment statuses
  6. Handler applies optional pretty-printing

- **Readiness Requests:** When `/readyz` is called:
  1. Handler receives request
  2. Calls `kubernetes.CheckAPIServer()` to probe API server
  3. Measures round-trip latency
  4. Returns cluster connectivity status with version and latency information

---

## Building & Running

### Build the Project
```bash
cd golang
go mod tidy && go build
```

### Run Against Local Kubernetes
```bash
./tyk-sre-assignment --kubeconfig ~/.kube/config --address ":8080"
```

### Run Against Current Context
```bash
# Uses current kubeconfig context from ~/.kube/config
./tyk-sre-assignment
```

### Run Inside Kubernetes Cluster
```bash
# Uses in-cluster service account
./tyk-sre-assignment --address ":8080"
```

### Unit Tests
```bash
go test -v
```

---

## Usage Examples

### Monitor Cluster Health
```bash
# Check all deployments
curl http://localhost:8080/deployments/health?pretty=true

# Check API server connectivity
curl http://localhost:8080/readyz

# Health probe (for liveness checks)
curl http://localhost:8080/healthz
```
