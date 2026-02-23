# Week 4 – Provisioning and Interaction via RESTful API

APIs are the front door to your platform. Make them robust, secure, and developer-friendly.

---

## Overview

This week focuses on:
- **RESTful API Development**: Building a production-ready API for PaaS product management
- **OpenAPI Specification**: Documenting API endpoints with industry-standard specifications
- **Unit Testing**: Implementing comprehensive tests for all API endpoints
- **Containerization**: Creating Docker images and deploying to STACKIT Container Registry
- **Kubernetes Deployment**: Provisioning the API on SKE with GitOps practices
- **Auto-Scaling**: Configuring Horizontal Pod Autoscaler (HPA) for production resilience
- **Performance Testing**: Load testing with k6 to validate scaling behavior

The implementation demonstrates **API-first development**, **container orchestration**, and **cloud-native scalability patterns**.

---

## Architecture

```
┌─────────────────────────────────────────────────┐
│           External Clients / Users              │
└──────────────────┬──────────────────────────────┘
                   │ HTTP/REST
┌──────────────────▼──────────────────────────────┐
│              PaaS API Service                   │
│    (RESTful API for Product Management)         │
│  - Create Instance    - Delete Instance         │
│  - List Instances     - Get Connection Info     │
└──────────────────┬──────────────────────────────┘
                   │ Kubernetes Client
┌──────────────────▼──────────────────────────────┐
│         PostgreSQL Cluster (CRD)                │
│          (CloudNativePG Operator)               │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────┐
│            STACKIT Kubernetes Engine            │
│              (Managed K8s Cluster)              │
└─────────────────────────────────────────────────┘
```

**Flow for Instance Creation:**
```
User Request → API Endpoint → Validate Input → Create K8s CR
    → Operator Reconciles → PostgreSQL Deployed → Return Instance Info
```

---

## Prerequisites

- STACKIT Kubernetes Engine (SKE) cluster running (from Week 3)
- CloudNativePG Operator installed
- `kubectl` configured for cluster access
- `docker` CLI for building container images
- Go 1.21+ for local development
- (Optional) k6 for performance testing

---

## API Implementation

### 1. API Structure

The API is built with Go and organized as follows:

```
paas-api/
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── config/config.go         # Configuration management
│   ├── http/
│   │   ├── routes.go            # Route definitions
│   │   └── handlers/
│   │       ├── handlers.go      # HTTP handlers
│   │       └── handlers_test.go # Unit tests
│   ├── kube/client.go           # Kubernetes client
│   ├── model/instance.go        # Data models
│   └── service/
│       ├── logic.go             # Business logic
│       └── ports.go             # Service ports
├── openapi.yaml                 # OpenAPI 3.0 specification
├── go.mod                       # Go dependencies
└── Dockerfile                   # Container image definition
```

### 2. API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check endpoint |
| `POST` | `/api/v1/instances` | Create a new PostgreSQL instance |
| `GET` | `/api/v1/instances` | List all instances |
| `GET` | `/api/v1/instances/{id}` | Get instance details and connection info |
| `DELETE` | `/api/v1/instances/{id}` | Delete an instance |

### 3. OpenAPI Specification

The API is fully documented using OpenAPI 3.0 specification.

**Viewing the Specification:**

 **Using Swagger Editor** (online):
   - Visit [editor.swagger.io](https://editor.swagger.io/)
   - Copy and paste the contents of `openapi.yaml`
   - Or import the file directly


**Key Features:**
- Complete request/response schemas
- Authentication requirements (if implemented)
- Error response definitions
- Example payloads
- Interactive API testing (with Swagger UI)

### 4. Running Unit Tests

```bash
cd Week_4/paas-api

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./internal/http/handlers
```

**Test Coverage:**
- Handler tests for all endpoints
- Mock Kubernetes client interactions
- Input validation testing
- Error handling scenarios

---

## Local Development

### 1. Install Dependencies

```bash
cd Week_4/paas-api
go mod download
```

### 2. Configure Kubernetes Access

Ensure your `kubeconfig.yml` points to the SKE cluster:

```bash
export KUBECONFIG=/path/to/Week_4/paas-api/kubeconfig.yml
```

### 3. Run the API Locally

```bash
go run cmd/api/main.go
```

The API will start on `http://localhost:8080`

### 4. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Create an instance
curl -X POST http://localhost:8080/api/v1/instances \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-db",
    "instances": 3,
    "storage": "10Gi"
  }'

# List instances
curl http://localhost:8080/api/v1/instances

# Get instance details
curl http://localhost:8080/api/v1/instances/test-db

# Delete instance
curl -X DELETE http://localhost:8080/api/v1/instances/test-db
```

---

## Container Image & Deployment

### 1. Build Docker Image

```bash
cd Week_4

# Build the image
docker build -t paas-api:latest -f Dockerfile paas-api/

# Tag for STACKIT Container Registry
docker tag paas-api:latest registry.stackit.cloud/your-project/paas-api:v1.0.0
```

### 2. Push to STACKIT Container Registry

```bash
# Login to STACKIT registry
docker login registry.stackit.cloud

# Push the image
docker push registry.stackit.cloud/your-project/paas-api:v1.0.0
```

### 3. Deploy to SKE

#### Manual Deployment

```bash
cd gitops/apps/paas-api

# Create namespace
kubectl apply -f namespace.yaml

# Deploy the application
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

#### GitOps Deployment with ArgoCD

```bash
# Apply ArgoCD application
kubectl apply -f gitops/argo/app-paas-api.yaml

# Check sync status
argocd app get paas-api

# Sync the application
argocd app sync paas-api
```

### 4. Verify Deployment

```bash
# Check pods
kubectl get pods -n paas-api

# Check service
kubectl get svc -n paas-api

# Check logs
kubectl logs -n paas-api -l app=paas-api
```

---

## Auto-Scaling with HPA

### 1. Configure Horizontal Pod Autoscaler

The HPA automatically scales the API based on CPU and memory utilization:

```bash
kubectl apply -f gitops/apps/paas-api/hpa.yaml
```

**HPA Configuration:**
- Min replicas: 2
- Max replicas: 10
- Target CPU: 70%
- Target Memory: 80%

### 2. Monitor HPA Status

```bash
# Check HPA status
kubectl get hpa -n paas-api

# Watch HPA in real-time
kubectl get hpa -n paas-api -w

# Describe HPA for detailed metrics
kubectl describe hpa paas-api -n paas-api
```

---

## Performance Testing

### 1. k6 Load Test Setup

The load tests are defined using k6 and deployed as Kubernetes Jobs.

```javascript
// gitops/apps/paas-api/loadtest/k6.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 },  // Ramp up to 100 users
    { duration: '5m', target: 100 },  // Stay at 100 users
    { duration: '2m', target: 200 },  // Ramp up to 200 users
    { duration: '5m', target: 200 },  // Stay at 200 users
    { duration: '2m', target: 0 },    // Ramp down to 0 users
  ],
};

export default function() {
  let response = http.get('http://paas-api.paas-api.svc.cluster.local:8080/health');
  check(response, { 'status is 200': (r) => r.status === 200 });
  sleep(1);
}
```

### 2. Run Performance Tests

```bash
# Apply ConfigMap with k6 script
kubectl apply -f gitops/apps/paas-api/loadtest/k6-configmap.yaml

# Run the load test job
kubectl apply -f gitops/apps/paas-api/loadtest/k6-job.yaml

# Monitor the test
kubectl logs -n paas-api -l job-name=k6-loadtest -f

# Watch HPA scale the pods
kubectl get hpa -n paas-api -w
```

### 3. Analyze Results

```bash
# Check final HPA state
kubectl describe hpa paas-api -n paas-api

# View pod scaling events
kubectl get events -n paas-api --sort-by='.lastTimestamp'

# Check number of running pods
kubectl get pods -n paas-api -l app=paas-api
```

**Expected Behavior:**
- API should scale from 2 to 10 replicas during load
- Response times should remain stable under load
- No errors or failed requests
- Automatic scale-down after load decreases

---

## Understanding the Create Flow

### Instance Creation Process

```mermaid
flowchart TD

    A[Client POST /instances] --> B{Validate Request}

    B -->|Invalid| E1[400 Bad Request]
    B -->|Valid| C{Auth / Quota Check}

    C -->|Failed| E2[401 / 403]
    C -->|OK| D[Create Instance Record<br/>Status: PROVISIONING]

    D --> F[Create Kubernetes<br/>Custom Resource]

    F --> G[Return 202 Accepted<br/>Instance ID]

    G --> H[Operator Reconciles<br/>Desired State]

    H --> I{Resource Ready?}

    I -->|No| J{Resource Failed?}
    J -->|Yes| K[Update Status: FAILED]
    J -->|No| H

    I -->|Yes| L[Read Connection Info<br/>from Service / Secret]
    L --> M[Update Status: READY]

    M --> N[Provisioning Complete]

    %% ---------- STYLES ----------

    style A fill:#90ee90,stroke:#2d5016,stroke-width:2px,color:#000
    style N fill:#90ee90,stroke:#2d5016,stroke-width:2px,color:#000

    style B fill:#4da6ff,stroke:#0066cc,stroke-width:2px,color:#fff
    style C fill:#4da6ff,stroke:#0066cc,stroke-width:2px,color:#fff
    style D fill:#4da6ff,stroke:#0066cc,stroke-width:2px,color:#fff
    style F fill:#4da6ff,stroke:#0066cc,stroke-width:2px,color:#fff
    style G fill:#ffa500,stroke:#cc6600,stroke-width:2px,color:#000

<!-- [MermaidChart: 4435770d-635c-4046-8254-f9ed94189783] -->
<!-- [MermaidChart: b3ece2c1-4585-436b-94bd-5960dfa41a2e] -->
    style H fill:#9966cc,stroke:#663399,stroke-width:2px,color:#fff
    style I fill:#9966cc,stroke:#663399,stroke-width:2px,color:#fff
    style J fill:#9966cc,stroke:#663399,stroke-width:2px,color:#fff
    style L fill:#9966cc,stroke:#663399,stroke-width:2px,color:#fff

    style M fill:#90ee90,stroke:#2d5016,stroke-width:2px,color:#000
    style K fill:#ff6b6b,stroke:#8b0000,stroke-width:2px,color:#fff

    style E1 fill:#ff6b6b,stroke:#8b0000,stroke-width:2px,color:#fff
    style E2 fill:#ff6b6b,stroke:#8b0000,stroke-width:2px,color:#fff
```

**Flow Stages:**

1. **Request Validation** (Blue) - API validates input and checks authorization
2. **Provisioning** (Orange) - Returns 202 with instance ID, provisioning starts asynchronously
3. **Reconciliation** (Purple) - Operator watches and reconciles Kubernetes resources
4. **Success** (Green) - Instance ready with connection information
5. **Failure** (Red) - Error handling at each stage

---

## Bonus Features

### 1. Automated API Deployment ✓

GitOps integration for automated deployment:
- **ArgoCD Application**: Declarative configuration in [`gitops/argo/app-paas-api.yaml`](../gitops/argo/app-paas-api.yaml)
- **Kubernetes Manifests**: Deployment, Service, HPA in [`gitops/apps/paas-api/`](../gitops/apps/paas-api/)
- **Automated Sync**: Changes to Git repository trigger automatic deployment
- **Self-healing**: ArgoCD automatically reconciles drift from desired state

### 2. Auto-Scaling and Performance Tests ✓

Horizontal scaling with load testing:
- **HPA Configuration**: [`gitops/apps/paas-api/hpa.yaml`](../gitops/apps/paas-api/hpa.yaml) - scales 2-10 replicas based on CPU/memory
- **k6 Load Tests**: [`gitops/apps/paas-api/loadtest/`](../gitops/apps/paas-api/loadtest/) - performance testing scripts
- **Validation**: Load tests verify HPA functionality under realistic traffic patterns
- **Metrics**: Real-time scaling behavior monitoring

### 3. Update Functionality

Implement an update endpoint for modifying instance specifications:

```
PATCH /api/v1/instances/{id}
```

**Supported Updates:**
- Number of replicas
- Storage size (expansion only)
- Resource limits
- Connection parameters

---

## Key Learnings

### RESTful API Design
- Resource-based URL structure
- Proper HTTP methods (GET, POST, DELETE, PATCH)
- Status codes and error handling
- Idempotent operations

### Kubernetes Integration
- Using the Kubernetes Go client
- Creating and managing Custom Resources
- Watching resource status
- Namespace isolation

### Container Best Practices
- Multi-stage Docker builds
- Minimal base images
- Non-root user execution
- Health checks and readiness probes

### Cloud-Native Scalability
- Horizontal Pod Autoscaler configuration
- Resource requests and limits
- Load testing methodology
- Performance optimization

### GitOps Workflow
- Declarative configuration
- ArgoCD application management
- Automated sync and rollback
- Version control for infrastructure

---

## Troubleshooting

### API Not Starting
```bash
# Check pod logs
kubectl logs -n paas-api -l app=paas-api

# Common issues:
# - Kubeconfig not mounted correctly
# - RBAC permissions missing
# - Service account not configured
```

### HPA Not Scaling
```bash
# Check metrics server
kubectl top nodes
kubectl top pods -n paas-api

# Verify HPA configuration
kubectl describe hpa paas-api -n paas-api

# Ensure resource requests are set in deployment
```

### Performance Test Failing
```bash
# Check k6 job logs
kubectl logs -n paas-api -l job-name=k6-loadtest

# Verify service endpoint
kubectl get svc -n paas-api

# Test connectivity from within cluster
kubectl run -it --rm debug --image=curlimages/curl --restart=Never -- \
  curl http://paas-api.paas-api.svc.cluster.local:8080/health
```

---

## Next Steps

Week 5-6 will add:
- **Web UI**: Vue.js dashboard for visual management
- **Ingress**: External access with TLS termination
- **Monitoring**: Prometheus and Grafana integration
- **Logging**: Centralized logging with Loki
- **Alerting**: Proactive incident detection

---

## Resources

- [Go HTTP Server Tutorial](https://go.dev/doc/tutorial/web-service-gin)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Kubernetes Go Client](https://github.com/kubernetes/client-go)
- [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
- [k6 Load Testing](https://k6.io/docs/)
- [Docker Multi-Stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [GitOps with ArgoCD](https://argo-cd.readthedocs.io/)
