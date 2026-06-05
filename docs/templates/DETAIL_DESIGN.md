---
artifact_type: detail_design
id: DESIGN-000
status: draft
owner: ai
approval: pending
human_fields:
  - approval
  - constraints
  - scope_decisions
ai_fields:
  - problem
  - context_loaded
  - brownfield_scope
  - proposed_approach
  - alternatives_considered
  - impacted_areas
  - test_plan
  - reconciliation_plan
shared_fields:
  - status
  - trace
  - small_task_exemption
trace:
  backlog_item: TBD
  requirement: TBD
  phase: TBD
  ticket_or_bug: TBD
  test_verification: TBD
  validation_matrix: TBD
  docs_review: TBD
  adrs: []
  master_docs_touched: []
---

# Detail Design Template

Use for non-small tickets and all meaningful bug fixes.

Before filling this template, read `docs/standards/QUALITY_BAR.md#good-detail-design`.

## Copy And Fill Rules

- AI drafts this document; Human approves or blocks it.
- Replace `DESIGN-000` with a stable design ID.
- Link the source ticket or bug before writing the approach.
- State assumptions explicitly. Do not hide product, API, data, security, or runtime assumptions.
- Include at least one realistic alternative for non-trivial work.
- If this design changes architecture, API, data, security, deployment, dependency, or standards, create or link an ADR.
- Implementation must not begin until approval is `approved`, unless `Small Task Exemption` is valid.

## Field Ownership

- Human owns approval, constraints, and scope decisions.
- AI owns analysis, proposed approach, alternatives, impact, test plan, and reconciliation plan.
- Shared fields include status, trace links, and small-task exemption.

## Status

- ID: DESIGN-000
- Status: draft | ready | in_progress | blocked | in_review | verified | done
- Ticket/Bug: TBD
- Approval: not required | pending | approved
- Author: TBD
- Updated: TBD

## Trace Links

- Backlog item: TBD
- Requirement: TBD
- Phase: TBD
- Ticket/Bug: TBD
- Test verification: TBD
- Validation matrix: TBD
- Docs review: TBD
- ADRs: TBD
- Master docs touched: TBD

## Problem

AI fill:

- Problem statement:
- Why now:
- Source ticket/bug:
- Success outcome:

## Context Loaded

- `docs/CONTEXT.md`
- TBD

## Brownfield Scope

- Touched modules/files: TBD
- Direct dependencies inspected: TBD
- Contracts affected: TBD
- Known unknowns: TBD
- Scope expansion reason: TBD

## Small Task Exemption

- Small task exemption: no
- Reason: TBD
- Impact checked: API=TBD, DB=TBD, Security=TBD, Runtime=TBD, Standards=TBD

## Proposed Approach

AI fill implementation-ready approach:

1. Change:
2. Data/control flow:
3. Boundaries:
4. Failure handling:
5. Compatibility:

## Alternatives Considered

| Alternative | Pros | Cons | Decision |
| --- | --- | --- | --- |
| TBD | TBD | TBD | rejected |

## Impacted Areas

- Code/modules: TBD
- Product behavior: TBD
- API/contracts: TBD
- Data/schema: TBD
- Security/auth: TBD
- Deployment/runtime: TBD
- Docs: TBD

## Bug Analysis

Use this section for bug work.

- Reproduction: TBD
- Root cause: TBD
- Impact scope: TBD
- Fix strategy: TBD
- Regression risk: TBD

## Test Plan

- Unit:
- Integration:
- E2E:
- UAT:
- Manual/platform:
- Docs review:

## Validation Matrix Impact

- Update required: yes | no
- Row(s): TBD
- Reason if no update required: TBD

## Reconciliation Plan

- Requirements docs: update | no change
- Architecture docs: update | no change
- API docs: update | no change
- ERD docs: update | no change
- ADR: create | no change
- Context: update

## Docs Review

- Code changed but docs unchanged reason: TBD
- Requirements updated or not needed reason: TBD
- Architecture updated or not needed reason: TBD
- API updated or not needed reason: TBD
- ERD/data updated or not needed reason: TBD
- ADR created or not needed reason: TBD

## Verification Results

Fill after execution.

- Command: TBD
- Result: TBD
- Notes: TBD
