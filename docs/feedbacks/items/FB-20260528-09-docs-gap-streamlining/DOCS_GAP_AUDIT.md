# Documentation Gap Audit

Severity scale:

- `HIGH`: likely to mislead implementation, demo, reviewer expectation, or future agents.
- `MEDIUM`: causes extra reading or repeated clarification, but is unlikely to break behavior by itself.
- `LOW`: cleanup or polish issue.

| ID | Severity | Area | Gap | Related Paths | Recommendation |
| --- | --- | --- | --- | --- | --- |
| GAP-01 | HIGH | Source ownership | Product truth is spread across product docs and feedback items. FB-08 currently carries important demo truth, but it is not yet clearly marked as planning-only, temporary canonical, or promoted product truth. | `docs/product/SPEC.md`, `docs/product/TECH.md`, `docs/product/GUIDE.md`, `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/*` | Add a docs governance rule and promote stable demo contract content into a product-level `DEMO_CONTRACT.md` or clearly mark FB-08 as temporary planning truth. |
| GAP-02 | HIGH | Product spec shape | `SPEC.md` is a mega-doc: product/system spec, business rules, logistics flow, farm operations, batch lifecycle, traceability reporting, and appendix-like sections live together. | `docs/product/SPEC.md` | Keep `SPEC.md` as reference until refactor; move durable business narrative to `domain/README.md` and technical current-state items to `TECH.md` or a current-state matrix. |
| GAP-03 | HIGH | Demo personas | Account/persona references diverge. `BOOTSTRAP.md` lists deterministic identities, while `DEMO_SETUP_RUNBOOK.md` says only admin is guaranteed by executable seeding and lists different persona emails. | `docs/product/standards/BOOTSTRAP.md`, `docs/product/DEMO_SETUP_RUNBOOK.md`, `docs/product/GUIDE.md`, seed files | Define one canonical demo-account table that distinguishes guaranteed executable seed, documented persona target, and manual/admin-created fallback. |
| GAP-04 | HIGH | Current vs planned runtime | Technical docs contain both current runtime details and aspirational/older architecture wording. Examples include `Monitor Service`, `SSE`, `OSRM`, and older order flow descriptions alongside Socket/WebSocket and browser-simulated route replay docs. | `docs/product/TECH.md`, `docs/product/GUIDE.md`, `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/LOGISTICS_REALTIME.md` | Create an implemented-vs-planned matrix and update `TECH.md` headings/labels so readers can tell current service names and protocols from target architecture. |
| GAP-05 | MEDIUM | Developer guide scope | `GUIDE.md` mixes developer guide, operational runbook, engineering log, demo dataset, persona table, and demo script. | `docs/product/GUIDE.md`, `docs/product/DEMO_SETUP_RUNBOOK.md` | Narrow `GUIDE.md` to developer learning and operations pointers. Move changelog/history to story archive or a release/history doc, and keep executable runbook steps in `DEMO_SETUP_RUNBOOK.md`. |
| GAP-06 | MEDIUM | UI gap status | `MISSING_UI.md` includes missing screens, existing screens, simulation requirements, and implemented walkthrough notes in one file. | `docs/product/ui-ux/MISSING_UI.md`, `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/DEMO_UI_AUDIT.md` | Split or relabel into current UI inventory, UI gap backlog, and demo simulation contract. |
| GAP-07 | MEDIUM | Docs index accuracy | `docs/README.md` lists legacy `TEST_MATRIX.md` and `HARNESS_BACKLOG.md` files that are not present in the current tree. | `docs/README.md`, Harness CLI durable layer | Update the docs map in a follow-up so absent legacy files are not presented as live main files. |
| GAP-08 | MEDIUM | Learning path | There is no single reading order for new contributors/reviewers who want to understand the technical demo, business demo, local run, auth, saga, trace, realtime, and UI boundaries. | `README.md`, `docs/README.md`, `docs/product/*`, `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/*` | Add `docs/product/LEARNING_PATH.md` with role-based reading paths: reviewer, backend engineer, frontend engineer, demo operator, and agent. |
| GAP-09 | MEDIUM | Docs governance | The repo lacks a concise rule for when to update `SPEC`, `TECH`, `GUIDE`, runbook, standards, feedback items, story history, or ADRs. | `docs/HARNESS.md`, `docs/FEEDBACK_WORKFLOW.md`, `docs/product/*` | Add `docs/product/DOCS_GOVERNANCE.md` or a short standards doc that defines ownership, promotion, archive, and stale-label rules. |
| GAP-10 | LOW | Heading and numbering consistency | Some long product docs restart numbering and use mixed title conventions. | `docs/product/TECH.md`, `docs/product/GUIDE.md`, `docs/product/SPEC.md` | Normalize headings during the refactor phase after source ownership is approved. |

## Missing Documentation Themes

- New-contributor learning sequence.
- Current-state matrix for implemented, simulated, planned, and deprecated behavior.
- Promotion process from feedback investigation into product docs.
- Canonical demo-account and role/persona table.
- Stable demo contract outside an individual feedback item.

