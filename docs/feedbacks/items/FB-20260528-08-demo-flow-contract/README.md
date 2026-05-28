# FB-20260528-08 Demo Flow Contract & Simulation Boundaries

Status: `IN_PROGRESS`

## Problem Statement

Runtime Roasters is a production-grade demo, but the current explanation of demo behavior is spread across product specs, technical docs, sprint history, feedback items, and UI pages. This makes it too easy for reviewers or agents to confuse real runtime behavior, backend-seeded simulation, browser route replay, and UI-only visual fallback.

## Goals

- Define one canonical demo flow contract.
- List intentionally omitted flows and why they are omitted.
- Explain simulation boundaries for backend simulation, browser simulation, and UI mock sections.
- Document multi-account demo setup with browser profile/session constraints.
- Document logistics realtime get/publish, trace CQRS, and order fulfillment gates.
- Provide a UI audit that classifies real API-backed screens, browser-simulated behavior, and mock/visualizer surfaces.

## Non-Goals

- No code changes.
- No UI implementation.
- No backend seeding implementation.
- No deletion of historical story docs.
- No attempt to convert planned flows into runtime behavior in this ticket.

## Canonical Docs

- [DEMO_FLOW_CONTRACT.md](DEMO_FLOW_CONTRACT.md)
- [LOGISTICS_REALTIME.md](LOGISTICS_REALTIME.md)
- [TRACE_CQRS.md](TRACE_CQRS.md)
- [ORDER_APPROVAL_AND_FULFILLMENT.md](ORDER_APPROVAL_AND_FULFILLMENT.md)
- [MULTI_ACCOUNT_DEMO.md](MULTI_ACCOUNT_DEMO.md)
- [DEMO_UI_AUDIT.md](DEMO_UI_AUDIT.md)

## Processing Docs

- [root_cause.md](root_cause.md)
- [impact_analysis.md](impact_analysis.md)
- [design.md](design.md)
- [implementation_plan.md](implementation_plan.md)
- [verification.md](verification.md)
- [RESEARCH.md](RESEARCH.md)
