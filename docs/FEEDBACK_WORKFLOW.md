# Feedback Workflow

This policy extends the Harness workflow for feedback captured by a human. It
defines how feedback becomes investigated work without allowing raw feedback to
trigger direct product or code changes.

## Source Of Feedback

Humans record feedback in:

```text
docs/feedbacks/USER_FEEDBACK.md
```

Agents must treat that file as human-authored input. Do not delete, rewrite, or
normalize the original human wording. Add links or status metadata only when a
feedback item is intentionally being processed.

Agents read feedback and create processing records only when explicitly asked to
triage, plan, or implement feedback. Do not triage existing feedback, create item
folders, fix issues, or change feedback status during unrelated work.

## Statuses

Use only these statuses:

| Status | Meaning |
| --- | --- |
| `TODO` | Recorded but not yet processed. |
| `IN_PROGRESS` | Being investigated, planned, implemented, or verified. |
| `RESOLVED` | Final conclusion exists, with fixed, deferred, rejected, or duplicate outcome when applicable, and clear reason plus evidence. |

Feedback can move to `RESOLVED` only after there is verification evidence or a
clear closing decision.

## Item Folders

Each feedback item that is processed must have a dedicated folder:

```text
docs/feedbacks/items/FB-YYYYMMDD-NN-short-title/
```

The main feedback entry in `docs/feedbacks/USER_FEEDBACK.md` must link to the
corresponding item folder once processing starts.

Required structure:

```text
docs/feedbacks/
├── USER_FEEDBACK.md
└── items/
    └── FB-YYYYMMDD-NN-short-title/
        ├── README.md
        ├── root_cause.md
        ├── impact_analysis.md
        ├── implementation_plan.md
        ├── verification.md
        └── design.md
```

File requirements:

| File | Required when |
| --- | --- |
| `README.md` | Always required for every processed feedback item. |
| `root_cause.md` | Required for bugs, mismatches, missing docs, or wrong context. |
| `impact_analysis.md` | Required when feedback may affect UI, API, data, auth, or another flow. |
| `implementation_plan.md` | Required before any fix or product-doc change. |
| `verification.md` | Required before moving feedback to `RESOLVED`. |
| `design.md` | Required only for larger issues that change design, business rules, API contracts, role/auth model, data model, topology/realtime behavior, deployment, or cross-service flows. |

## Processing Rules

Do not change code or product documentation directly from raw feedback before a
root cause exists.

Root cause analysis must distinguish:

- Symptom.
- Actual cause.
- Docs, code, test, or context mismatch.
- Missing or misleading source of truth.

Implementation plans must state:

- Expected changes.
- Impacted scope.
- Required validation.
- Expected impact after the fix.

Large feedback items still follow `docs/FEATURE_INTAKE.md`. If feedback becomes
a larger work item, classify it through feature intake and use the appropriate
tiny, normal, or high-risk lane before implementation.

## Resolution Rules

Before setting `RESOLVED`, `verification.md` must contain either:

- Evidence that the fix works, such as test output, screenshots, traces, or
  command results.
- A clear closure decision for deferred, rejected, or duplicate feedback,
  including the reason and any supporting evidence.

If the issue was implemented, verification must match the validation listed in
`implementation_plan.md` or explain why a planned validation could not run.

## Harness Integration

This document is a supplemental Harness policy. Use it together with:

- `docs/HARNESS.md` for the agent task loop and durable records.
- `docs/FEATURE_INTAKE.md` for classifying feedback that becomes work.
- `scripts/harness query matrix` for proof status when feedback produces or
  changes executable behavior.
