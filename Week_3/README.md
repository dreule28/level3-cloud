# Week 3 – Developing a Platform-as-a-Service (PaaS) Product on SKE

Provisioning a managed Kubernetes cluster and deploying a PaaS product using Operators and Custom Resources.

---

## Overview

This week focuses on:
- **SKE Cluster Creation**: Using the STACKIT Terraform Provider to provision a managed Kubernetes cluster
- **PaaS Product Implementation**: Deploying a managed PostgreSQL database service
- **Operator Deployment**: Installing and configuring the CloudNativePG Operator
- **Custom Resources**: Managing database instances through Kubernetes Custom Resources (CRs)
- **Connectivity**: Demonstrating how to connect to and use the managed database service

The implementation showcases the **Operator pattern** and deepens understanding of **Custom Resource Definitions (CRDs)** and **reconciliation loops**.

---

## Architecture

```
┌─────────────────────────────────────────────────┐
│              User Application                   │
│         (Connects to PostgreSQL)                │
├─────────────────────────────────────────────────┤
│         PostgreSQL Cluster Instance             │
│            (Managed by Operator)                │
├─────────────────────────────────────────────────┤
│          CloudNativePG Operator                 │
│       (Watches CRs, Reconciles State)           │
├─────────────────────────────────────────────────┤
│            STACKIT Kubernetes Engine            │
│              (Managed K8s Cluster)              │
└─────────────────────────────────────────────────┘
```

---

## Prerequisites

- STACKIT account with valid credentials
- `terraform` CLI (>= 1.6.0)
- `kubectl` CLI

---

## Infrastructure Setup

### 1. Configure STACKIT Credentials

Ensure your STACKIT credentials are configured:

```bash
export STACKIT_SERVICE_ACCOUNT_TOKEN="your-token"
export STACKIT_SERVICE_ACCOUNT_EMAIL="your-email@example.com"
```

### 2. Provision SKE Cluster with Terraform

```bash
cd Week_3
terraform init
terraform apply
```

**What this creates:**
- SKE cluster named `week3-paas` in project `your Project-ID`
- Single node pool with 1 worker node
- Machine type: `g1a.2d` (Ubuntu 22.04)
- 100GB premium storage per node

### 3. Configure kubectl Access

After the cluster is provisioned, download the kubeconfig:

```bash
# Download kubeconfig from STACKIT Portal
# Navigate to: Kubernetes > Your Cluster > Download Kubeconfig
# Save the file and move it to your kubectl config directory

mv ~/Downloads/kubeconfig-week3-paas.yaml ~/.kube/week3-config

# Set as current context
export KUBECONFIG=~/.kube/week3-config

# Verify access
kubectl get nodes
```

---

## PaaS Product Implementation

### PostgreSQL as a Service Using CloudNativePG

#### 1. Install the CloudNativePG Operator

```bash
# Add the CloudNativePG Helm repository
helm repo add cnpg https://cloudnative-pg.github.io/charts
helm repo update

# Install the operator using Helm
helm install cnpg \
  --namespace cnpg-system \
  --create-namespace \
  cnpg/cloudnative-pg
```

**What this does:**
- Creates the `cnpg-system` namespace automatically
- Installs all required CRDs (Cluster, Backup, ScheduledBackup, Pooler)
- Deploys the CloudNativePG operator controller
- Sets up RBAC permissions

#### 2. Verify Operator Installation

```bash
# Check operator pods
kubectl get pods -n cnpg-system

# Check CRDs
kubectl get crds | grep postgresql
```

Expected output:
```
backups.postgresql.cnpg.io
clusters.postgresql.cnpg.io
poolers.postgresql.cnpg.io
scheduledbackups.postgresql.cnpg.io
```

#### 3. Create Namespace for Database Instances

```bash
kubectl create namespace paas-postgres
```

#### 4. Deploy PostgreSQL Instance

Apply the Custom Resource to provision a PostgreSQL cluster:

```bash
kubectl apply -f pg-demo.yml
```

**What this does:**
- Creates a PostgreSQL cluster named `pg-demo`
- Provisions 1 instance (pod) running PostgreSQL
- Allocates 10Gi persistent storage
- The operator automatically handles database initialization, configuration, and management

---

## Operator Pattern & Reconciliation

### Understanding Custom Resources

The `pg-demo.yml` file defines a **Custom Resource** of kind `Cluster`:

```yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: pg-demo
  namespace: paas-postgres
spec:
  instances: 1
  storage:
    size: 10Gi
```

This is not a standard Kubernetes resource - it's a custom definition created by the CloudNativePG operator.

### How the Operator Works

1. **Watch**: The operator continuously watches for `Cluster` resources
2. **Compare**: When detected, it compares the desired state (spec) with actual state
3. **Reconcile**: If there's a difference, the operator takes action to align reality with the spec
4. **Repeat**: This loop runs continuously, ensuring the cluster stays in the desired state

### What Gets Created

When you apply `pg-demo.yml`, the operator automatically creates:
- StatefulSet for PostgreSQL pods
- Persistent Volume Claims (PVCs)
- Services for database connectivity
- ConfigMaps and Secrets for configuration
- Monitoring resources (if enabled)

---

## Connectivity & Usage

### 1. Check Cluster Status

```bash
kubectl get cluster -n paas-postgres
kubectl describe cluster pg-demo -n paas-postgres
```

### 2. Get Connection Information

```bash
# List all services
kubectl get svc -n paas-postgres

# List all secrets created by the operator
kubectl get secrets -n paas-postgres
```

**Get the app user password:**

```bash
kubectl get secret pg-demo-app -n paas-postgres -o jsonpath='{.data.password}' | base64 -d && echo
```

### 3. Connect to PostgreSQL

Run an ephemeral psql client pod that auto-deletes when you exit:

```bash
kubectl run psql-client --rm -it --image=postgres:16 --namespace paas-postgres -- bash
```

Inside the pod, connect to PostgreSQL:

```bash
# Connect as app user
psql -h pg-demo-rw -U app -d app
```

**Alternative: Connect directly to the database pod:**

```bash
kubectl exec -it pg-demo-1 -n paas-postgres -- psql -U postgres
```

### 4. Test the Database

```sql
-- You're already in the app database, create a test table
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  email VARCHAR(100)
);

-- Insert test data
INSERT INTO users (name, email) VALUES
  ('Alice', 'alice@example.com'),
  ('Bob', 'bob@example.com');

-- Query
SELECT * FROM users;
```

---

## Verification & Monitoring

### Check Pod Status

```bash
kubectl get pods -n paas-postgres
```

Expected output:
```
NAME        READY   STATUS    RESTARTS   AGE
pg-demo-1   1/1     Running   0          5m
```

### View Logs

```bash
kubectl logs -n paas-postgres pg-demo-1
```

### Inspect Storage

```bash
kubectl get pvc -n paas-postgres
```

### Check Operator Logs

```bash
kubectl logs -n cnpg-system -l app.kubernetes.io/name=cloudnative-pg
```

---

## Scaling the Database

To scale the PostgreSQL cluster, simply update the spec:

```yaml
spec:
  instances: 3  # Changed from 1
  storage:
    size: 10Gi
```

Then apply:

```bash
kubectl apply -f pg-demo.yml
```

The operator will automatically:
- Create additional replicas
- Configure replication
- Update services
- Balance the load

### Verify Replication Across All Instances

After scaling, verify all 3 pods are running:

```bash
kubectl get pods -n paas-postgres
```

Expected output:
```
NAME        READY   STATUS    RESTARTS   AGE
pg-demo-1   1/1     Running   0          15m
pg-demo-2   1/1     Running   0          2m
pg-demo-3   1/1     Running   0          2m
```

**Prove that all instances have the same data:**

Connect to each pod and query the users table (using postgres user for direct pod access):

```bash
# Query from pod 1 (primary)
kubectl exec -it pg-demo-1 -n paas-postgres -- psql -U postgres -d app -c "SELECT * FROM users;"

# Query from pod 2 (replica)
kubectl exec -it pg-demo-2 -n paas-postgres -- psql -U postgres -d app -c "SELECT * FROM users;"

# Query from pod 3 (replica)
kubectl exec -it pg-demo-3 -n paas-postgres -- psql -U postgres -d app -c "SELECT * FROM users;"
```

All three queries should return identical results:
```
 id | name  |       email
----+-------+-------------------
  1 | Alice | alice@example.com
  2 | Bob   | bob@example.com
(2 rows)
```

**Test write to primary and read from replicas:**

```bash
# Insert new data
kubectl exec -it pg-demo-1 -n paas-postgres -- psql -U postgres -d app -c \
  "INSERT INTO users (name, email) VALUES ('Charlie', 'charlie@example.com');"

# Immediately read from a replica
kubectl exec -it pg-demo-2 -n paas-postgres -- psql -U postgres -d app -c "SELECT * FROM users;"
```

You should see the new row replicated across all instances, demonstrating PostgreSQL streaming replication in action!

---

## Cleanup

### Delete PostgreSQL Instance

```bash
kubectl delete -f pg-demo.yml
```

### Uninstall Operator

```bash
# Uninstall the Helm release
helm uninstall cnpg -n cnpg-system

# Clean up namespaces
kubectl delete namespace paas-postgres
kubectl delete namespace cnpg-system
```

---

## Bonus: GitOps with ArgoCD

The `gitops/` directory in the project root contains manifests for managing this PaaS product using GitOps principles with ArgoCD. This demonstrates declarative, Git-based infrastructure and application management.

**What's included:**
- ArgoCD Application manifests for automated deployment
- Kubernetes manifests organized for GitOps workflows
- Will be expanded in later weeks for full GitOps implementation

See [gitops/README.md](../gitops/README.md) for more details.

---

## References

- [STACKIT Terraform Provider](https://registry.terraform.io/providers/stackitcloud/stackit/latest/docs)
- [CloudNativePG Documentation](https://cloudnative-pg.io/)
- [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

---

### Recommended Alias Commands

```bash
alias k='kubectl'
alias kgp='kubectl get pods'
alias kgs='kubectl get svc'
alias kgc='kubectl get cluster'
alias kdesc='kubectl describe'
alias klogs='kubectl logs'

# Add to ~/.bashrc or ~/.zshrc
source ~/.bashrc  # or source ~/.zshrc
```

---
