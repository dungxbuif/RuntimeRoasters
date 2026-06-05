---
version: 1.0.0
name: The Architectural Narrator
colors:
  surface: "#f7f9fb"
  surface-container-low: "#f2f4f6"
  surface-container-lowest: "#ffffff"
  surface-container-high: "#e6e8ea"
  primary: "#004ac6"
  primary-container: "#2563eb"
  secondary: "#515f74"
  tertiary: "#006242"
  tertiary-fixed: "#6ffbbe"
  outline: "#737686"
  outline-variant: "#c3c6d7"
  on-surface: "#191c1e"
  on-surface-variant: "#434655"
  inverse-surface: "#2d3133"
  inverse-on-surface: "#ffffff"
typography:
  display: "Space Grotesk"
  body: "Inter"
  annotation: "Indie Flower"
spacing:
  sm: "0.25rem"
  md: "0.75rem"
  lg: "1rem"
  xl: "1.5rem"
---

# Design System: The Architectural Narrator

## Overview
The **Architectural Narrator** is the visual soul of Runtime Roasters. It is designed to move away from sterile, standard documentation toward an editorial, high-clarity experience. The system interprets complex technical concepts through a "human-in-the-loop" lens, blending technical precision with instructional warmth.

### Concept: The Editorial Mix
The UI is divided into two distinct mental models:
- **Control Plane (/control/*):** A BFF-backed monitoring environment for Admins and Developers. It emphasizes service health, API exploration, and observability metrics.
- **Business UI (/app/*):** A direct browser-to-gateway interface for Operators. It uses metaphors like "Green Beans" (Farm events) and "Roasted Beans" (Order events) to visualize the supply chain.

---

## Colors
The palette is rooted in a crisp `surface` (#f7f9fb) that mimics high-quality bleached paper. Color is functional, never decorative.

- **Surface Layers:** Hierarchy is achieved by stacking tonal layers (`lowest` -> `low` -> `base` -> `high`).
- **Functional Accent:** `primary` (#004ac6) drives core interactions.
- **Status & Logic:** `tertiary` (#006242) denotes "Live" or "Success" states.
- **Highlighter:** `tertiary-fixed` (#6ffbbe) is used for highlighter-style annotations in technical diagrams.

---

## Typography
We use a high-contrast pairing to balance engineered precision with human storytelling.

- **Display & Headlines:** *Space Grotesk*. Geometric and tech-forward.
- **Body & Titles:** *Inter*. Optimized for technical legibility.
- **Annotations:** *Indie Flower / Gaegu*. Handwritten marker aesthetic for whiteboard-style notes and diagrams.

---

## Layout
We utilize intentional asymmetry and expansive white space.

- **The "No-Line" Rule:** Major content blocks are defined by background color shifts, not 1px solid borders.
- **The Split-Screen Pattern:** A persistent Sidebar navigation for Supply Chain domains (Farms, Batches, Logistics, Warehouse, Retail).
- **Diagram Spacing:** Technical illustrations must have at least `2rem` of internal padding to ensure they "breathe."

---

## Elevation & Depth
Depth is created through **Tonal Layering** and **Ambient Light** rather than heavy drop shadows.

- **Ambient Shadow:** `box-shadow: 0 12px 32px -4px rgba(25, 28, 30, 0.06)`. Diffused and subtle.
- **The Ghost Border:** For card containment, use `outline-variant` at 15% opacity to suggest a boundary.
- **Glassmorphism:** Modals and tooltips use 70% opacity with a `20px` backdrop-blur for a premium feel.

---

## Shapes
Corners and containers follow a strict roundedness scale.

- **xl (1.5rem):** Primary cards and technical illustration containers.
- **md (0.75rem):** Buttons and interactive elements.
- **sm (0.25rem):** Tooltips, inputs, and minor annotations.
- **Minimalist Fields:** Form inputs have no background fill, utilizing only a bottom "Ghost Border" that transforms to a 2px `primary` line on focus.

---

## Components

### Technical Illustration Cards
The heart of the system.
- **The Stroke:** Lines use an SVG displacement map to create a hand-drawn, marker feel.
- **The Logic:** `primary` is used only for the component being explained; `outline-variant` for everything else.

### Mapping to Routes
| Design Folder | Description | Target Route |
| :--- | :--- | :--- |
| `control_plane_visualization_system_diagram` | Visual System Diagram | `/control/diagram` |
| `chaos_control_resiliency_logic`| Chaos Engineering Dashboard | `/control/chaos` |
| `farm_origin_traceability` | Supply Chain Farm Overview | `/app/farms` |
| `logistics_real_time_transit` | Real-time Logistics Tracking | `/app/logistics` |
| `traceability_journey_map` | Batch/Harvest Journey Map | `/app/batches` |
| `retail_saga_orchestrator` | Retail & Order Flow | `/app/retail` |
| `warehouse_stock_logic` | Inventory & Warehouse Management | `/app/warehouse` |

---

## Do's and Don'ts

### Do
- **Do** overlap elements. Marker-style arrows "breaking" card boundaries adds to the hand-crafted feel.
- **Do** use `on_surface` (#191c1e) for text. Never use pure black.
- **Do** lean into white space to reduce cognitive load on complex pages.

### Don't
- **Don't** use standard 90-degree corners. Always use the roundedness scale.
- **Don't** use "Default" system icons. Use thin-stroke (1.5px) icons matching the system's weight.
- **Don't** use secondary colors as text on dark backgrounds.
