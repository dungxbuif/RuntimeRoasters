# Root Cause Analysis

**Symptom**: Architecture Diagram on the home page does not have a fixed layout, and nodes are floating freely without a clear logical grouping.
**Actual Cause**: The `ArchitectureDiagramCanvas.tsx` uses React Flow without explicit grouping nodes (Sub-flows) or a rigid layout engine (like Dagre) to enforce a left-to-right (Frontend -> Backend -> DB) layout.
**Missing/Misleading Info**: No explicit layout engine configured for the topology diagram.
