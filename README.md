# golden-path-platform
Self-service API delivery platform on Kubernetes using Knative (Golden Path): tenant isolation (RBAC/Quota/NetworkPolicy) + CLI-based deployments.

## What this repo provides
- Platform manifests (cluster admin): `deploy/platform`
- Tenant isolation templates (RBAC/Quota/NetworkPolicy/LimitRange): `internal/templates/tenant`
- Developer CLI `gpp`:
  - `gpp tenant create <tenant>`: create namespace + isolation defaults
  - `gpp deploy -f gpp.yaml`: deploy Knative Service into tenant namespace

## Quick start (local with kind)
1) Create cluster
2) Install ingress + Knative
3) Install platform base
4) Create tenant
5) Deploy example

(See docs/quickstart.md)
