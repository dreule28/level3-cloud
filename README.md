# Building a Cloud-Native Platform: From IaaS to PaaS

**Meta-Track**: Infrastructure → Platform → Product

Building a complete cloud platform stack from scratch - starting with OpenStack IaaS, provisioning Kubernetes, and developing a PaaS product with APIs, UI, and observability.


## Table of Contents

- [Architecture](#architecture)
- [Progress](#progress)
- [Technologies](#technologies)


## Architecture

Layered approach from infrastructure to application:

```
┌─────────────────────────────────────────────────────────────┐
│          Week 5: Web UI + Ingress + TLS                     │
│  (Vue.js SPA, nginx Ingress, cert-manager, GitOps deploy)  │
├─────────────────────────────────────────────────────────────┤
│             Week 4: RESTful API Layer                       │
│  (Go + Echo, JWT Auth, OpenAPI, HPA, k6, GitOps deploy)    │
├─────────────────────────────────────────────────────────────┤
│           Week 3: PaaS Product on SKE                       │
│      (Operators, CRDs, Managed PostgreSQL Service)          │
├─────────────────────────────────────────────────────────────┤
│          Week 2: Kubernetes on OpenStack                    │
│           (IaC with Terraform + Ansible)                    │
├─────────────────────────────────────────────────────────────┤
│           Week 1: IaaS with OpenStack                       │
│                (DevStack Installation)                      │
└─────────────────────────────────────────────────────────────┘
```

## Progress

### Week 1-2: IaaS Foundations & Kubernetes Provisioning ✓

OpenStack setup with DevStack, Infrastructure as Code with Terraform, automated Kubernetes installation with Ansible.

→ [Documentation](Week_1+2/README.md)

### Week 3: Developing a PaaS Product on SKE ✓

Provisioning SKE cluster with STACKIT Terraform Provider, deploying Kubernetes Operators, implementing managed database service using Custom Resources.

→ [Documentation](Week_3/README.md)

### Week 4: Provisioning and Interaction via RESTful API ✓

Go + Echo RESTful API with JWT authentication (HS256, role-based), OpenAPI 3.0 spec, unit tests with mock service, Docker image pushed to STACKIT Container Registry, GitOps deployment via ArgoCD, HPA (1–10 replicas), and k6 performance tests.

→ [Documentation](Week_4+5/README.md)

### Week 5: Web UI and Secure Public Access ✓

Vue.js SPA with Pinia state management, JWT-authenticated Axios client, Vue Router auth guards, GSAP animations. Deployed on SKE behind nginx Ingress with TLS (cert-manager / Let's Encrypt) on public STACKIT subdomains. Automated via ArgoCD GitOps.

→ [Documentation](Week_4+5/README.md)

---

## Technologies

**Infrastructure**: OpenStack (Nova, Neutron, Cinder, Glance), Terraform, Ansible

**Platform**: Kubernetes, STACKIT Kubernetes Engine (SKE), Kubernetes Operators, Custom Resource Definitions, CloudNativePG

**API**: Go, Echo, JWT (HS256), OpenAPI 3.0, Docker, STACKIT Container Registry

**Scalability**: Horizontal Pod Autoscaler, k6 Performance Testing

**UI**: Vue.js, Vite, Pinia, Vue Router, Axios, GSAP

**Networking**: nginx Ingress Controller, cert-manager, Let's Encrypt TLS

**GitOps**: ArgoCD (automated sync, self-heal, pruning)

## What You'll Learn

**Infrastructure Management**: OpenStack architecture, VM provisioning, network configuration, storage setup

**Infrastructure as Code**: Terraform workflows, state management, provider orchestration

**Kubernetes**: Cluster architecture, CRDs, Operator pattern, reconciliation loops

**Platform Engineering**: PaaS design, managed services, API-first development, multi-tenancy

**Operations**: CI/CD, GitOps, monitoring, logging, audit trails, performance testing

---
