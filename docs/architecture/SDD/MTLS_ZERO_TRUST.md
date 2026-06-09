---
artifact_type: sdd
id: MTLS_ZERO_TRUST
status: proposed
owner: tech-lead
---

# SDD: mTLS and Zero Trust Security Mesh

## Context

- **Problem**: In the current architecture, internal service-to-service communication occurs over unencrypted HTTP/gRPC. If an internal node is compromised, sensitive supply chain and retail data can be intercepted or forged (Lateral Movement).
- **Business goal**: Establish 100% cryptographic certainty for every internal interaction, ensuring that only authorized services can talk to each other.
- **Constraints**: 
  - Must not significantly increase latency.
  - Must support automatic certificate rotation (no manual PKI management).
  - Must be compatible with the existing OIDC/Ory stack.

## System Design

- **Workload Identity**: Adopt **SPIFFE** (Secure Production Identity Framework for Everyone). Each service is assigned a "SVID" (SPIFFE Verifiable Identity Document) in the form of an X.509 certificate.
- **Service Mesh Layer**: Use **Linkerd** or **Istio** as the control plane.
  - **Sidecar Proxy**: Every service container gets a lightweight proxy (e.g., Linkerd-proxy) that intercepts traffic.
  - **Transparent mTLS**: Proxies perform the mTLS handshake automatically without application code changes.
- **Data flow**:
  - `Retail Service` -> `Warehouse Service` (Handshake verified by proxies).
  - `Farm Service` -> `Kafka` (Identity-based authentication for topics).

## Security And Privacy

- **Authn/authz impact**: 
  - **Authentication**: Moves from "trusted IP" to "cryptographic identity".
  - **Authorization**: Service-level RBAC (e.g., `Audit Service` can ONLY consume from Kafka, cannot call `Payment Service`).
- **Public/private data split**: Public traffic remains on standard TLS; internal traffic is strictly mTLS.
- **Audit/logging impact**: Every mTLS handshake failure is logged as a security event in SigNoz/Audit Service.

## Operational Impact

- **Deployment**: Requires a Kubernetes-like environment or Docker Compose with a sidecar pattern.
- **Observability**: Proxies provide "Gold Signals" (Latency, Error, Traffic) automatically for every link in the mesh.
- **Migration/backfill**: Can be rolled out service-by-service using "mTLS optional" mode before enforcing "mTLS required".

## Verification

- **Integration**: Verify that a non-meshed service (e.g., a generic `curl` container) is blocked from calling internal APIs.
- **E2E**: System Journey must pass through the encrypted mesh without performance degradation.

## Reconciliation

- **Architecture**: Complements the "Digital Birth Certificate" by ensuring the "Delivery Pipe" is equally secure.
- **ADRs**: Proposed ADR-0012: Adopt Linkerd for Workload Identity.

---
*Prepared by Gemini CLI - Zero Trust Initiative.*
