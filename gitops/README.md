# GitOps Directory

This directory contains Kubernetes manifests and ArgoCD Application definitions for managing infrastructure and applications using GitOps principles.

---

## Overview

GitOps is a declarative way to manage infrastructure and applications where Git serves as the single source of truth. Changes are made through Git commits, and tools like ArgoCD automatically sync the desired state to the cluster.

---

## Structure

```
gitops/
├── argo/                      # ArgoCD Application definitions
│   ├── app-cnpg.yml          # Manages CloudNativePG operator
│   └── app-postgres.yml      # Manages PostgreSQL instances
└── apps/                      # Application manifests
    └── postgres/             # PostgreSQL PaaS resources
        ├── namespace.yml     # Namespace definition
        └── clusterr.yml      # PostgreSQL cluster definition
```

---

## Components

### ArgoCD Applications

Located in `argo/`:

- **app-cnpg.yml**: Deploys the CloudNativePG operator from the official Helm chart
  - Installs CRDs automatically
  - Creates `cnpg-system` namespace
  - Enables automated sync and self-healing

- **app-postgres.yml**: Manages PostgreSQL cluster instances
  - References manifests in `apps/postgres/`
  - Creates `paas-postgres` namespace
  - Syncs PostgreSQL Custom Resources

### Application Manifests

Located in `apps/postgres/`:

- **namespace.yml**: Defines the `paas-postgres` namespace
- **clusterr.yml**: PostgreSQL Cluster custom resource managed by CloudNativePG operator

---

## Usage

### Prerequisites

- ArgoCD installed in your cluster
- Access to the repository (if using private repo)

### Deploy with ArgoCD

1. **Install ArgoCD** (if not already installed):
```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

2. **Apply the ArgoCD Applications**:
```bash
# Deploy CloudNativePG operator
kubectl apply -f gitops/argo/app-cnpg.yml

# Deploy PostgreSQL instances
kubectl apply -f gitops/argo/app-postgres.yml
```

3. **Monitor sync status**:
```bash
kubectl get applications -n argocd
```

### Access ArgoCD UI

```bash
# Port-forward to ArgoCD server
# Protip -> run this in a seperate terminal
kubectl port-forward svc/argocd-server -n argocd 8080:443

# Get admin password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d && echo

# Access UI at https://localhost:8080
# Username: admin
# Password: (from command above)
```

---

## Continuous Integration (CI)

This GitOps directory is protected by automated CI checks that run on every pull request and push to main. The CI pipeline validates all manifests before they reach the cluster.

### CI Pipeline Steps

Located in [`.github/workflows/ci-gitops.yml`](../.github/workflows/ci-gitops.yml), the pipeline performs three validation stages:

1. **YAML Linting** (`yamllint`)
   - Validates YAML syntax and formatting
   - Enforces consistent style (document start markers, line endings, etc.)
   - Catches common YAML mistakes early

2. **Kubernetes Schema Validation** (`kubeconform`)
   - Validates manifests against Kubernetes API schemas
   - Ensures resource definitions are structurally correct
   - Runs in strict mode with summary output
   - Ignores missing schemas for CRDs

3. **ArgoCD Application Sanity Checks**
   - Verifies ArgoCD Application manifests exist
   - Ensures referenced local paths exist in the repository
   - Validates Application kind and basic structure
   - Confirms in-cluster destination server configuration

### Benefits of CI

- **Early error detection**: Catches syntax and schema errors before deployment
- **Consistency**: Enforces formatting and structural standards
- **Safety**: Prevents invalid manifests from reaching the cluster
- **Documentation**: CI checks serve as living documentation of requirements
- **Confidence**: Pull requests show validation status before merge

### Running Checks Locally

```bash
# YAML lint
yamllint gitops/

# Kubernetes validation
find gitops/apps -type f \( -name '*.yml' -o -name '*.yaml' \) -exec kubeconform -strict {} \;

# Application checks
test -f gitops/argo/app-operator-cnpg.yml && test -f gitops/argo/app-product-postgres.yml
test -d gitops/apps/postgres-cluster
```

---

## GitOps Workflow

1. **Make changes**: Edit manifests in this directory
2. **Commit & push**: Push changes to Git repository
3. **ArgoCD syncs**: ArgoCD detects changes and applies them to the cluster
4. **Automated healing**: If someone manually changes resources, ArgoCD reverts to Git state

---

## Benefits

- **Declarative**: Infrastructure defined as code in Git
- **Version controlled**: Full history of all changes
- **Auditable**: Know who changed what and when
- **Automated**: No manual kubectl apply commands needed
- **Self-healing**: Cluster state automatically matches Git state
- **Rollback**: Easy rollback by reverting Git commits

---

## Future Extensions

This GitOps structure will be expanded in later weeks to include:
- Additional PaaS products
- Multi-environment management (dev/staging/prod)
- Application deployments
- Monitoring and observability stack
- Security policies and network configurations

---

## References

- [ArgoCD Documentation](https://argo-cd.readthedocs.io/)
- [GitOps Principles](https://opengitops.dev/)
- [CloudNativePG Operator](https://cloudnative-pg.io/)
