# Research

## Scope

This research pass reviewed the documentation surfaces that currently shape reviewer and agent understanding:

- `README.md`
- `docs/README.md`
- `docs/HARNESS.md`
- `docs/FEATURE_INTAKE.md`
- `docs/ARCHITECTURE.md`
- `docs/product/SPEC.md`
- `docs/product/TECH.md`
- `docs/product/GUIDE.md`
- `docs/product/DEMO_SETUP_RUNBOOK.md`
- `docs/product/domain/README.md`
- `docs/product/standards/BOOTSTRAP.md`
- `docs/product/ui-ux/MISSING_UI.md`
- `docs/feedbacks/USER_FEEDBACK.md`
- `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/*`

## Evidence

### Top-Level Documentation

- `README.md` presents a clean entrypoint and links to `TECH.md`, `GUIDE.md`, `DEMO_SETUP_RUNBOOK.md`, `domain/README.md`, `SPEC.md`, UI/UX notes, roadmap, and ADRs.
- `docs/README.md` still lists `TEST_MATRIX.md` and `HARNESS_BACKLOG.md` as main files, while also saying current proof/backlog state lives in the Harness CLI. Those legacy files do not appear in the current docs tree.
- `docs/HARNESS.md` defines a useful source hierarchy: product docs are current product contract, feedback workflow is supplemental policy, stories are work packets/history, Harness CLI stores proof state, and ADRs explain why contracts changed.

### Product Docs

- `docs/product/domain/README.md` is already described as the BA-facing source of truth from `SPEC.md` and `GUIDE.md`.
- `docs/product/SPEC.md` starts with a useful note pointing readers to `domain/README.md` and FB-08 for demo truth, then continues as a large mixed document: business specification, farm operations, batch lifecycle, traceability, role responsibility, and historical-style appendix material in one file.
- `docs/product/TECH.md` is the strongest technical source, but its headings restart numbering and it contains some older target/current wording. Early sections still name `Monitor Service`, `SSE`, and `OSRM` while later sections describe Socket/WebSocket and current demo contracts.
- `docs/product/GUIDE.md` begins as a developer guide, then includes operational reference, webhook notes, socket topology notes, Kafka topics, decision links, an engineering log, milestones, and master demo data. It is useful but overloaded.
- `docs/product/DEMO_SETUP_RUNBOOK.md` is a clearer operational runbook. It explicitly calls out the current seeded-account mismatch and says some personas must be created through Admin/User flow if absent.
- `docs/product/standards/BOOTSTRAP.md` declares deterministic identities such as `mgr.caudat@runtimeroasters.com`, `mgr.songthan@runtimeroasters.com`, and `driver.songthan01@runtimeroasters.com`, while `DEMO_SETUP_RUNBOOK.md` currently lists only `admin@runtimeroasters.com` as guaranteed by executable seeding and uses different documented persona emails.
- `docs/product/ui-ux/MISSING_UI.md` contains both missing UI analysis and implemented-screen walkthrough. The title says missing UI, but some sections mark existing/mostly complete surfaces, so readers need extra context to know whether it is current contract or gap backlog.

### Feedback Items

- `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/` contains detailed planning and current demo truth labels, including demo flow contract, logistics realtime, trace CQRS, multi-account demo, UI audit, and order approval/fulfillment.
- Product docs already link to FB-08 for current demo/runtime boundaries. That makes FB-08 semi-canonical today, but there is no explicit rule for when feedback-item findings are promoted into `docs/product/*`.
- Several feedback folders contain investigation documents that are valuable during planning, but their status and relationship to product truth are not standardized beyond `USER_FEEDBACK.md`.

## Main Research Conclusion

The repo does not need less documentation by default. It needs a clearer lifecycle:

1. Product docs should state current accepted truth.
2. Standards should hold durable engineering rules.
3. Runbooks should hold executable operating steps.
4. Feedback items should hold investigation/planning records until promoted.
5. Story history should remain archive.

