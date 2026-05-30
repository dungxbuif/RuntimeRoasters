# Missing Docs Backlog

| ID | Priority | Proposed Path | Purpose | Source Inputs |
| --- | --- | --- | --- | --- |
| DOC-BL-01 | HIGH | `docs/product/LEARNING_PATH.md` | Give new readers a short route through the repo by persona: reviewer, backend, frontend, demo operator, and agent. | `README.md`, `docs/README.md`, `GUIDE.md`, `TECH.md`, `DEMO_SETUP_RUNBOOK.md` |
| DOC-BL-02 | HIGH | `docs/product/DEMO_CONTRACT.md` | Promote stable current demo truth out of FB-08 so product docs have a durable canonical demo boundary. | FB-08 demo contract docs |
| DOC-BL-03 | HIGH | `docs/product/DOCS_GOVERNANCE.md` | Define when to update `SPEC`, `TECH`, `GUIDE`, runbook, standards, feedback items, story history, and ADRs. | `docs/HARNESS.md`, `docs/FEEDBACK_WORKFLOW.md`, this audit |
| DOC-BL-04 | HIGH | `docs/product/CURRENT_STATE.md` | Track implemented, backend-simulated, browser-simulated, UI mock-only, planned, and deprecated behavior. | `TECH.md`, `GUIDE.md`, FB-08, `MISSING_UI.md` |
| DOC-BL-05 | HIGH | `docs/product/DEMO_ACCOUNTS.md` or a dedicated section in `BOOTSTRAP.md` | Reconcile guaranteed executable seed accounts, target deterministic personas, and manual fallback account creation. | `BOOTSTRAP.md`, `DEMO_SETUP_RUNBOOK.md`, executable seed files |
| DOC-BL-06 | MEDIUM | `docs/product/ui-ux/UI_CURRENT_STATE.md` | Separate implemented UI inventory from missing UI backlog and mock/demo visualization notes. | `MISSING_UI.md`, FB-08 UI audit, current client routes |
| DOC-BL-07 | MEDIUM | `docs/product/TRACE_AND_REALTIME.md` or section in `DEMO_CONTRACT.md` | Explain Trace CQRS, Socket/WebSocket, browser route replay, and polling fallback in one current-state place. | FB-08 trace/logistics docs, `TECH.md`, `VERIFICATION_PROTOCOL.md` |
| DOC-BL-08 | MEDIUM | `docs/product/CHANGELOG.md` or archive note | Move durable engineering log/history out of `GUIDE.md` if the team wants a product-level history. | `GUIDE.md`, `docs/stories/history/*` |
| DOC-BL-09 | LOW | `docs/product/README.md` | Add a compact product-doc index with status labels and ownership. | Source-of-truth map |

## Suggested First Follow-Up Ticket

Start with `DOC-BL-01`, `DOC-BL-03`, and a small `docs/README.md` cleanup. Those are low-risk and make later refactors safer because they establish navigation and governance before content moves.

