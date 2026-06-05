# 📋 Task Management Conventions (Jira-in-Markdown)

This document specifies the structure and process for managing tasks (tickets) within the RuntimeRoasters project using Git/Markdown, simulating a Jira experience to ensure a clear distinction between Business and Technical aspects.

---

## 🏗️ 1. Folder Structure

Every task (Ticket) must reside within its own dedicated directory.

### A. For Single Tickets
```text
docs/stories/history/sprintX/RR-x/
├── ticket.md             # BA Documentation (What/Why)
└── technical_design.md    # Technical Documentation (How)
```

### B. For Large Tickets (Epic/Complex Ticket)
If a ticket is too large and needs to be subdivided, use a `subtickets` directory:
```text
docs/stories/history/sprintX/RR-x/
├── ticket.md             # BA Epic level
├── technical_design.md    # Tech Lead/Architect Strategy
└── subtickets/
    └── RR-x.y/           # Sub-directory for each sub-ticket
        ├── ticket.md     # BA Sub-ticket level
        └── technical_design.md # Dev Implementation Plan
```

---

## 🎭 2. Role & Content Differentiation

### 📝 ticket.md (Role: BA / Product Owner)
- **Language:** Business-oriented, non-technical.
- **Content:**
    - **User Story:** "As a... I want... so that..."
    - **Business Value:** Why is this feature important?
    - **Acceptance Criteria (AC):** Business testing scenarios (Scenario-based).
    - **Priority & Type:** Priority level and classification (Feature, Bug, Infra).

### 🛠️ technical_design.md (Role: Developer / Tech Lead)
- **Language:** Deeply technical.
- **Content:**
    - **Strategy:** Approach (Which library to choose? Why?).
    - **Implementation Steps:** Specific coding steps (Copy boilerplate, setup DI, SQL Schema).
    - **Verification:** Technical verification methods (Unit tests, Curl commands, DB checks).
    - **Trade-offs:** Assessments of performance or technical limitations.

---

## 🔄 3. Ticket Lifecycle Process

1.  **Refinement Phase:** The BA writes `ticket.md` to describe the requirements.
2.  **Technical Planning Phase:** The Developer reads `ticket.md`, analyzes it, and writes `technical_design.md`. 
    - *Purpose:* Finalize the scope and solution before coding.
3.  **Implementation Phase:** The Developer proceeds with coding based on `technical_design.md`.
4.  **Closing Phase:** 
    - (Optional) Update `technical_design.md` with actual changes if they deviate from the original plan.
    - Change status in the Sprint's Kanban board (`main.md`).

---

## 📏 4. General Rules

1.  **Language:** Use professional English for all documentation.
2.  **Linking:** Use relative paths to ensure links remain intact when moving directories.
3.  **Ticket Size:** 
    - If the ticket is small enough: Write `technical_design.md` directly under the task folder.
    - If the ticket is large: A `subtickets` folder must be created for management.
4.  **Consistency:** Do not create isolated Markdown files outside task directories (except for index files like `main.md` or `sprint-planning.md`).
