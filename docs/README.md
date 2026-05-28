# Documentation Map

This directory holds the project harness and the product contract.

## Main Files

- `HARNESS.md`: how humans and agents collaborate.
- `FEEDBACK_WORKFLOW.md`: how human feedback becomes investigated work.
- `FEATURE_INTAKE.md`: how prompts become tiny, normal, or high-risk work.
- `ARCHITECTURE.md`: high-level architecture rules and stack overview.
- `TEST_MATRIX.md`: legacy proof map; current proof status is queried with
  `scripts/bin/harness-cli query matrix`.
- `HARNESS_BACKLOG.md`: legacy improvement list; current improvement records
  are stored with `scripts/bin/harness-cli backlog`.
- `GLOSSARY.md`: shared terms.

## Folders

- `product/`: current product truth.
    - `SPEC.md`: The broader technical feature specs and history.
    - `TECH.md`: System architecture, technical contracts, and infrastructure.
    - `GUIDE.md`: Developer guides, engineering logs, and operational setup.
    - `domain/README.md`: **Master Business Specification** (Source of Truth).
- `stories/`: feature packets and history of previous work.
- `decisions/`: architecture decision records (ADRs).
- `templates/`: reusable spec-intake, story, plan, decision, and validation
  formats.

## Current State

Harness v0 is active. The product documentation has been migrated to the new
structure, representing the current state of the **Runtime Roasters**
microservices ecosystem.
