# Root Cause

## Symptom

Demo documentation and UI labels do not consistently identify which behaviors are real runtime flows, which are simulated for demo speed, and which are only visual placeholders. Reviewers and agents can interpret a UI mock or browser route replay as production-equivalent backend behavior.

## Actual Causes

- Historical sprint docs remain useful, but they are no longer the current source of truth.
- Product docs describe both target architecture and currently implemented demo behavior, sometimes in the same document.
- The UI contains real API-backed pages, browser-simulated controls, and mock/visualizer sections.
- Logistics uses browser route replay plus backend GPS publish, which can be mistaken for real vehicle telematics.
- Trace CQRS is a real read-model projection, but empty seed state makes it look broken unless a known entity ID exists.
- "Order approval" is currently an event-driven fulfillment gate, not a complex human approval workflow.
- Multi-account demos conflict with browser-scoped Kratos cookies and app tokens in `localStorage`.

## Source-of-Truth Gap

There is no single document that says:

- what the demo storyline is,
- what each flow's truth label is,
- what is intentionally omitted,
- what reviewers should see during a demo,
- what agents should not implement from stale story history.

## Current Ticket Scope

This ticket creates a documentation contract and planning packet only. It does not change backend behavior, frontend behavior, seeding logic, or authorization rules.
