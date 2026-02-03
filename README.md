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
┌─────────────────────────────────────────────────┐
│         Week 5-6: UI + Observability            │
│    (Web UI, Ingress, Monitoring, Logging)       │
├─────────────────────────────────────────────────┤
│           Week 4: RESTful API Layer             │
│         (Product Management API)                │
├─────────────────────────────────────────────────┤
│        Week 3: PaaS Product on SKE              │
│    (Operators, CRDs, Managed Services)          │
├─────────────────────────────────────────────────┤
│       Week 2: Kubernetes on OpenStack           │
│      (IaC with Terraform + Ansible)             │
├─────────────────────────────────────────────────┤
│        Week 1: IaaS with OpenStack              │
│           (DevStack Installation)               │
└─────────────────────────────────────────────────┘
```

## Progress

### Week 1-2: IaaS Foundations & Kubernetes Provisioning ✓

OpenStack setup with DevStack, Infrastructure as Code with Terraform, automated Kubernetes installation with Ansible.

→ [Detailed documentation](Week_1+2/README.md)

### Week 3: Developing a PaaS Product on SKE (In Progress)

Provisioning SKE cluster with STACKIT Terraform Provider, deploying Kubernetes Operators, implementing managed database service using Custom Resources.

Documentation: `Week_3/README.md` (to be created)

---

## Technologies

**Infrastructure**: OpenStack (Nova, Neutron, Cinder, Glance), Terraform, Ansible

**Platform**: Kubernetes, STACKIT Kubernetes Engine (SKE), Kubernetes Operators, Custom Resource Definitions

**Planned**: Go APIs, Vue.js UI, PostgreSQL operator, Prometheus, Grafana, Loki

## What You'll Learn

**Infrastructure Management**: OpenStack architecture, VM provisioning, network configuration, storage setup

**Infrastructure as Code**: Terraform workflows, state management, provider orchestration

**Kubernetes**: Cluster architecture, CRDs, Operator pattern, reconciliation loops

**Platform Engineering**: PaaS design, managed services, API-first development, multi-tenancy

**Operations**: CI/CD, GitOps, monitoring, logging, audit trails, performance testing

---
