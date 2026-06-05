# Debugging Standard

This file defines how agents handle bugs, test failures, build failures, production issues, and unexpected behavior.

## Principle

No fixes without root cause investigation.

Agents MUST NOT patch symptoms, stack multiple guesses, or claim a bug is fixed without reproduction and verification evidence.

## Required Debugging Flow

```text
Reproduce -> Investigate Root Cause -> Compare Patterns -> Form Hypothesis -> Create Regression Proof -> Fix -> Verify -> Reconcile Docs
```

## Phase 1: Reproduce

Before proposing a fix:

- [ ] Record exact steps, command, input, or environment that triggers the issue.
- [ ] Record expected behavior and actual behavior separately.
- [ ] Determine whether the issue is consistent, intermittent, or environment-specific.
- [ ] Capture relevant error messages, stack traces, logs, screenshots, or terminal output summaries.

If the bug cannot be reproduced, gather more evidence or mark the bug blocked. Do not guess.

## Phase 2: Root Cause Investigation

Before changing implementation:

- [ ] Read error messages and stack traces completely.
- [ ] Check recent changes, dependencies, config, and environment.
- [ ] Trace data/control flow backward from symptom to source.
- [ ] For multi-component systems, inspect each boundary where data/config/state crosses.
- [ ] Identify the earliest source of incorrect behavior.

Root cause must be written in the bug or detail design.

## Phase 3: Pattern Comparison

Before choosing a fix:

- [ ] Find similar working code or previous decisions in the repo.
- [ ] Compare broken and working paths.
- [ ] Identify differences that could explain the failure.
- [ ] Check whether the current architecture or contract is being violated.

## Phase 4: Hypothesis

Write a single hypothesis:

```text
I think <root cause> causes <symptom> because <evidence>.
```

Then test one variable at a time. Do not bundle multiple speculative changes.

## Phase 5: Regression Proof

Before or alongside the fix:

- [ ] Create or identify a failing test/reproduction that proves the bug.
- [ ] Define the regression check that must pass after the fix.
- [ ] Update `docs/work/VALIDATION_MATRIX.md` when behavior proof changes.

## Phase 6: Fix And Verify

After root cause and regression proof:

- [ ] Implement the smallest fix that addresses the root cause.
- [ ] Run regression proof.
- [ ] Run relevant surrounding tests.
- [ ] Record command/manual evidence.
- [ ] Reconcile requirements, architecture, API, ERD/data, ADR, and context as needed.

## Stop Conditions

Stop and ask for human/design review when:

- The same failure path still fails after 3 fix attempts.
- The task reaches 5 total fix/test cycles.
- Each fix reveals a new shared-state, coupling, or architecture problem.
- The issue cannot be reproduced but impact is high.
- The fix requires broad refactoring not covered by the active ticket/bug.
- The root cause implies an architecture, API, data, security, or standards decision.

When stopping, mark the work item `blocked` and record:

- Attempt count
- Commands run
- Changes made in each attempt
- Failure output summaries
- Current root-cause hypothesis
- Human/design decision needed

## Red Flags

If an agent writes or thinks any of these, return to root cause investigation:

- "Try this and see."
- "Quick fix for now."
- "It is probably X."
- "I will skip the test."
- "I do not fully understand, but..."
- "One more fix attempt."
- "Docs can be updated later."
