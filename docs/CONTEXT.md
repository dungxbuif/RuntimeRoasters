# Session Context & Development State

Last Updated: 2026-05-30

## Active Task
- Added a new agent constraint in `AGENTS.md` to export session/task context to `docs/CONTEXT.md`.
- Initialized `docs/CONTEXT.md` with current system development state and context rules.

## Current Context & System Shape
- **Mission**: Farm-to-Cup supply-chain system built as a Go microservices monorepo with Next.js control-plane UI.
- **Client App Port**: `http://localhost:3000`
- **KrakenD Gateway**: `http://localhost:8081`

## Agent Constraints (Updated)
- Always export the current session/task context, active changes, and outstanding tasks to `docs/CONTEXT.md` before concluding.
- Follow Clean Architecture patterns and central authorization designs.
- Strictly adhere to specified user/manager role behaviors.
