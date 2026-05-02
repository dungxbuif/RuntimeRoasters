```markdown
# Design System Specification: The Architectural Narrator

## 1. Overview & Creative North Star
The Creative North Star for this design system is **"The Technical Storyteller."** 

We are moving away from the cold, sterile nature of standard documentation and moving toward an editorial, high-clarity experience. This system interprets complex technical concepts through a "human-in-the-loop" lens. By blending the precision of high-end typography with the "organic imperfection" of marker-style illustrations, we create an environment that feels both authoritative and approachable.

The design breaks the traditional "SaaS dashboard" mold by utilizing intentional asymmetry, expansive white space, and a unique layering principle that replaces rigid lines with tonal depth. We aren't just building a UI; we are crafting a digital textbook for the modern era.

---

## 2. Colors & Surface Logic
The palette is rooted in a crisp `background` (#f7f9fb), designed to mimic high-quality bleached paper. Color is never decorative; it is functional, used to denote flow, status, and logic.

### The "No-Line" Rule
Standard 1px solid borders for sectioning are strictly prohibited. Boundaries between major content blocks must be defined solely through background color shifts.
- Use `surface-container-low` to set apart a sidebar or a secondary content area.
- Use `surface-container-lowest` (#ffffff) for primary content cards to make them "pop" against the `surface` background.

### Surface Hierarchy & Nesting
Treat the UI as a series of physical layers. Hierarchy is achieved by "stacking" container tiers:
1.  **Base Layer:** `surface` (#f7f9fb)
2.  **Section Layer:** `surface-container-low` (#f2f4f6)
3.  **Interactive Layer (Cards):** `surface-container-lowest` (#ffffff)
4.  **Highlight/Pop-out:** `surface-container-high` (#e6e8ea)

### The "Glass & Gradient" Rule
To prevent the UI from feeling "flat" or "cheap," floating elements (like modals or hovering tooltips) should utilize Glassmorphism. Use `surface-container-lowest` with a 70% opacity and a `20px` backdrop-blur. 

For primary calls to action or hero illustrations, use a subtle linear gradient from `primary` (#004ac6) to `primary_container` (#2563eb) at a 135-degree angle. This adds a "digital soul" that flat hex codes cannot provide.

---

## 3. Typography: The Editorial Mix
We use a high-contrast typographic pairing to balance technical precision with instructional warmth.

*   **Display & Headlines (Space Grotesk):** This font carries the brand’s "tech-forward" personality. Its geometric nature feels engineered. Use `display-lg` through `headline-sm` for all major headings.
*   **Body & Titles (Inter):** Inter provides world-class legibility for dense technical explanations. Use `body-md` for general documentation and `title-sm` for card headers.
*   **Annotations (Handwritten/Marker Aesthetic):** For technical illustrations and "side-notes," use a handwritten-style font (e.g., *Indie Flower* or *Gaegu*) in `secondary` (#515f74) to mimic a mentor drawing on a whiteboard.

---

## 4. Elevation & Depth
We eschew traditional "Drop Shadows" in favor of **Tonal Layering** and **Ambient Light.**

*   **The Layering Principle:** Depth is created by placing a `surface-container-lowest` card on a `surface-container-low` background. The contrast in value creates a natural lift without the "dirtiness" of a shadow.
*   **Ambient Shadows:** If an element must float (e.g., a floating action button), use an extra-diffused shadow: `box-shadow: 0 12px 32px -4px rgba(25, 28, 30, 0.06)`. Note the low opacity; it should feel like ambient light, not a dark glow.
*   **The "Ghost Border":** For card containment, avoid high-contrast outlines. Use the `outline-variant` (#c3c6d7) at **15% opacity**. This creates a "suggestion" of a boundary that keeps the layout feeling airy.

---

## 5. Components

### Buttons & Interaction
*   **Primary Button:** Uses the `primary` to `primary_container` gradient. Roundedness is `md` (0.75rem).
*   **Secondary/Tertiary:** No background. Use `on_surface` text with the "Ghost Border" logic on hover.
*   **The "Marker" Interaction:** When hovering over interactive diagram elements, use a `tertiary_fixed` (#6ffbbe) highlight—mimicking a highlighter pen stroke.

### Technical Illustration Cards
*   **Container:** `surface-container-lowest` with an `xl` (1.5rem) corner radius.
*   **Internal Lines:** All diagram lines should use a "marker" style—slightly varying thicknesses and non-perfectly straight lines. Use `outline` (#737686) for neutral flows and `primary` for the "Golden Path."
*   **Spacing:** Use generous padding (at least `2rem`) to ensure technical diagrams have room to breathe.

### Inputs & Fields
*   **Field Style:** Minimalist. No background fill; only a bottom "Ghost Border" that transitions to a 2px `primary` line on focus. 
*   **Labels:** Use `label-md` in `on_surface_variant` (#434655).

### Tooltips & Annotations
*   **Visual Style:** Rounded `sm` (0.25rem). Background is `inverse_surface` (#2d3133) with `inverse_on_surface` text.
*   **Hand-drawn Pointers:** Use marker-style arrows to point from text to specific parts of an illustration.

---

## 6. Do’s and Don’ts

### Do
*   **Do** use `tertiary` (#006242) for "Success" or "Live" states to keep the vibe friendly and instructional.
*   **Do** lean into white space. If a page feels "busy," increase the vertical spacing between sections using the `xl` roundedness scale as a spacing guide.
*   **Do** overlap elements. Having a marker-style arrow "break" the boundary of a card and point to a heading creates an editorial, hand-crafted feel.

### Don’t
*   **Don't** use pure black (#000000). Always use `on_surface` (#191c1e) for text to maintain a premium, ink-on-paper look.
*   **Don't** use standard 90-degree corners. Everything in this system—from cards to selection states—must utilize the Roundedness Scale (minimum `sm`).
*   **Don't** use "Default" system icons. Use custom, thin-stroke (1.5px) icons that match the `outline` weight of the technical illustrations.

---

## 7. Signature Technical Illustration Style
The diagrams are the heart of this system. 
1. **The Stroke:** Use a SVG filter to apply a "rough" displacement map to lines, giving them a hand-drawn feel.
2. **The Fill:** Use `primary_fixed` (#dbe1ff) or `secondary_fixed` (#d5e3fc) for block fills within diagrams, always with a slightly "messy" edge that doesn't perfectly align with the border.
3. **The Logic:** Use the `primary` color only for the specific component being explained; all supporting infrastructure should be in `outline_variant`.