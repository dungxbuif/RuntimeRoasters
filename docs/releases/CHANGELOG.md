# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]
### Added
- New Harness v1 (Markdown-First) framework.
- Core architecture docs (API, ERD, INTEGRATIONS).
- Unified standard docs.
- Reviewed Harness v1 runtime docs: feedback log, traceability matrix, requirements index, user stories, engineering setup/local/troubleshooting guides, work-area indexes/templates, and RR-URG-07A through RR-URG-07E split tickets.

### Changed
- Reconciled migrated docs with Harness v1 folder structure and replaced stale pre-Harness product/story/feedback references with current Harness paths.
- Updated RR-URG-07 scope to public sold-cup QR trace using UI-issued `product_id` values.
- Consolidated reviewed legacy non-template docs into Harness target files and removed duplicate legacy docs.
- Replaced legacy seed notes with a canonical master-data/enum contract covering migration flow, Admin seed flow, reproducible demo variation, and RR-URG-07A retail data.
- Removed the duplicate farm administration role; `FARM_MANAGER` is now the only farm-domain operator role across identity, authorization, backend, frontend, seed data, and documentation.
