# Streamlining Plan

This plan is intentionally staged. The current ticket creates the audit only; later tickets can apply the refactor after review.

## Phase 0: Approve Ownership

Deliverables:

- Review [DOCS_GAP_AUDIT.md](DOCS_GAP_AUDIT.md).
- Review [SOURCE_OF_TRUTH_MAP.md](SOURCE_OF_TRUTH_MAP.md).
- Confirm which documents become canonical for business, technical, runbook, demo contract, UI status, and feedback planning.

Exit criteria:

- Reviewers agree that feedback items are planning records unless explicitly promoted.
- Reviewers agree that story history is archive only.

## Phase 1: Add Navigation And Governance

Create or update:

- `docs/product/LEARNING_PATH.md`
- `docs/product/DOCS_GOVERNANCE.md`
- `docs/product/CURRENT_STATE.md` or an implemented-vs-planned matrix inside `TECH.md`
- `docs/README.md` to remove or relabel absent legacy files.

Rules:

- Do not duplicate large sections.
- Link to canonical docs instead of copying tables.
- Add status labels where a file contains mixed current/reference/planning content.

Exit criteria:

- A new reader can identify the correct first five docs to read.
- A future agent can decide where to place a docs update without guessing.

## Phase 2: Promote Stable Demo Contract

Create:

- `docs/product/DEMO_CONTRACT.md`

Promote stable content from:

- `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/DEMO_FLOW_CONTRACT.md`
- `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/LOGISTICS_REALTIME.md`
- `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/TRACE_CQRS.md`
- `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/ORDER_APPROVAL_AND_FULFILLMENT.md`
- `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/MULTI_ACCOUNT_DEMO.md`
- `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/DEMO_UI_AUDIT.md`

Exit criteria:

- Product docs link to `DEMO_CONTRACT.md` for current demo boundaries.
- FB-08 remains as investigation/history, not the permanent canonical doc.

## Phase 3: Split Or Label Mega-Docs

Recommended target:

- `docs/product/domain/README.md`: business narrative, role responsibility, and durable business rules.
- `docs/product/TECH.md`: current architecture, contracts, storage, messaging, auth, realtime, observability.
- `docs/product/GUIDE.md`: developer guide and learning pointers only.
- `docs/product/DEMO_SETUP_RUNBOOK.md`: local install/start/reset/seed/demo operation only.
- `docs/product/SPEC.md`: broader reference/background or reduced product overview after canonical content is promoted.

Exit criteria:

- `SPEC.md`, `TECH.md`, and `GUIDE.md` each have one primary job.
- Old target/current wording is either removed or explicitly labeled.

## Phase 4: Normalize UI And Account Docs

Create or update:

- A UI current inventory and UI gap backlog from `docs/product/ui-ux/MISSING_UI.md`.
- A canonical account/persona table that reconciles `BOOTSTRAP.md`, `DEMO_SETUP_RUNBOOK.md`, `GUIDE.md`, and executable seed files.

Exit criteria:

- Demo operators know which accounts are guaranteed immediately after seed.
- Reviewers can distinguish implemented UI, partial UI, planned UI, and mock/visualizer UI.

## Phase 5: Maintenance Sweep

Actions:

- Normalize heading numbering.
- Remove duplicated tables where links are enough.
- Add backlinks from promoted docs to original feedback tickets.
- Update Harness trace/proof records after each docs refactor ticket.

Exit criteria:

- `rg` for legacy names such as absent docs or old service/protocol names returns either no matches or explicitly labeled archive/planned references.

