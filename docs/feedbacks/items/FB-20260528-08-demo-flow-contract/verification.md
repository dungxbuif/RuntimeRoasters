# Verification

## Docs-Only Checklist

- [x] All new files exist in `docs/feedbacks/items/FB-20260528-08-demo-flow-contract/`.
- [x] `docs/feedbacks/USER_FEEDBACK.md` links to the ticket folder.
- [x] `DEMO_FLOW_CONTRACT.md` states the current truth for real, backend-simulated, browser-simulated, UI mock, and planned flows.
- [x] No code changes are included.
- [x] Multi-account demo instructions mention browser profile, Kratos cookie, and `localStorage` constraints.
- [x] `LOGISTICS_REALTIME.md` includes current implementation state, simulated pieces, omitted pieces, and recommended demo script.
- [x] `TRACE_CQRS.md` includes current implementation state, simulated pieces, omitted pieces, and recommended demo script.
- [x] `ORDER_APPROVAL_AND_FULFILLMENT.md` includes current implementation state, simulated pieces, omitted pieces, and recommended demo script.

## Evidence

- `find docs/feedbacks/items/FB-20260528-08-demo-flow-contract -maxdepth 1 -type f | sort` listed all required and recommended docs.
- `rg` confirmed the feedback link, truth labels, multi-account session constraints, `logistics.gps.updated`, and `coffee_traceability` references.
- `git status --short` showed only documentation changes plus the new feedback docs folder.

## Suggested Commands

```bash
find docs/feedbacks/items/FB-20260528-08-demo-flow-contract -maxdepth 1 -type f | sort
```

```bash
git diff --name-only
```

Expected: changed files are documentation and Harness local metadata only.
