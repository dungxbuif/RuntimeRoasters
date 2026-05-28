# Source Of Truth Map

This map proposes where readers should look after the documentation streamlining work is approved. It does not move files in this ticket.

## Reader Paths

| Reader | Start Here | Then Read | Purpose |
| --- | --- | --- | --- |
| Project reviewer | `README.md` | `docs/product/LEARNING_PATH.md` once created, then `docs/product/DEMO_CONTRACT.md` once promoted | Understand the demo narrative and what is implemented. |
| Business/BA reviewer | `docs/product/domain/README.md` | `docs/product/SPEC.md` as reference | Understand role responsibilities, supply-chain flow, and business rules. |
| Backend engineer | `docs/product/TECH.md` | `docs/product/standards/ENGINEERING_RULES.md`, ADRs, relevant service code | Understand current architecture, service contracts, persistence, messaging, and auth. |
| Frontend engineer | `docs/product/ui-ux/DESIGN.md` | `docs/product/ui-ux/MISSING_UI.md` or successor UI gap backlog, `docs/product/DEMO_CONTRACT.md` | Understand dashboard language, UI status, and simulation boundaries. |
| Demo operator | `docs/product/DEMO_SETUP_RUNBOOK.md` | `docs/product/standards/BOOTSTRAP.md`, `docs/product/DEMO_CONTRACT.md` | Start/reset/seed/run the demo with known account expectations. |
| Agent handling feedback | `docs/FEEDBACK_WORKFLOW.md` | `docs/feedbacks/USER_FEEDBACK.md`, selected feedback item, affected product docs | Investigate or implement a bounded feedback item without promoting unreviewed truth accidentally. |

## Canonical Ownership

| Topic | Canonical Doc | Supporting Docs | Notes |
| --- | --- | --- | --- |
| Business narrative and roles | `docs/product/domain/README.md` | `docs/product/SPEC.md` | `SPEC.md` should become reference/background where it overlaps. |
| Current technical architecture | `docs/product/TECH.md` | `docs/ARCHITECTURE.md`, ADRs | `TECH.md` should describe current contracts and explicitly label future targets. |
| Developer operations and learning | `docs/product/GUIDE.md` | `docs/product/DEMO_SETUP_RUNBOOK.md`, standards | Keep detailed run commands in the runbook and link out from the guide. |
| Local start/reset/demo operation | `docs/product/DEMO_SETUP_RUNBOOK.md` | `docs/product/standards/BOOTSTRAP.md` | This should be the executable local operator path. |
| Seed and first-run rules | `docs/product/standards/BOOTSTRAP.md` | executable seed files, `DEMO_SETUP_RUNBOOK.md` | Must distinguish guaranteed executable seed from target personas. |
| Demo runtime/simulation boundary | Proposed `docs/product/DEMO_CONTRACT.md` | FB-08 until promoted | FB-08 is useful, but product truth should not permanently live only inside a feedback item. |
| Engineering/testing rules | `docs/product/standards/*` | ADRs and Harness docs | Durable rules, not per-ticket plans. |
| Feedback investigations | `docs/feedbacks/items/*` | `docs/feedbacks/USER_FEEDBACK.md` | Planning and evidence records until accepted and promoted. |
| Historical execution | `docs/stories/history/*` | `docs/stories/ROADMAP.md` | Archive/reference only. Do not treat as current product contract unless linked from current docs. |
| Architecture decisions | `docs/decisions/*` | `docs/product/TECH.md` | Explain why a durable decision exists. |

## Status Labels To Use

| Label | Meaning |
| --- | --- |
| `CANONICAL_CURRENT` | Accepted current truth. Implementation and docs should align with it. |
| `REFERENCE` | Useful background or expanded context; not the primary source if it conflicts with canonical docs. |
| `PLANNING_ONLY` | Proposed direction or investigation output; not accepted product truth yet. |
| `ARCHIVE` | Historical record. Preserve for context, but do not implement from it without checking current docs. |
| `DEPRECATED` | Known stale or replaced material. Keep only when historical context is needed and label clearly. |

