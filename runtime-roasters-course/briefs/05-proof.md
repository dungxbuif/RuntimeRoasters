# Module 5: Transparent Proof

### Teaching Arc
- **Metaphor:** The "Unbreakable Chain". Imagine a physical chain where each link is welded shut. You can follow the chain from the end back to the very first link, knowing nothing was removed or faked.
- **Opening hook:** You're holding a cup with a QR code. When you scan it, how does the system "remember" which specific farm the beans came from?
- **Key insight:** "Traceability" is about linking unique IDs. A Cup points to a Sale, which points to a Roast Batch, which points to a Harvest.
- **"Why should I care?":** Understanding data relationships (Lineage) helps you build "Search" and "Analytics" features. You'll know how to ask AI for "Foreign Keys" or "Join logic" to connect the dots.

### Code Snippets (pre-extracted)

File: src/apps/trace-service/internal/domain/trace.go (lines 15-30)
```go
type TraceDocument struct {
	ProductID    string             `json:"product_id"`
	TraceCode    string             `json:"trace_code"`
	Product      ProductInfo        `json:"product"`
	Store        StoreInfo          `json:"store"`
	HarvestID    string             `json:"harvest_id"`
	OriginBatch  string             `json:"origin_batch_id"`
	Milestones   []TraceMilestone   `json:"milestones"`
}
```

File: src/apps/trace-service/internal/usecase/public_trace.go (lines 50-65)
```go
func (u *Usecase) GetPublicTrace(ctx context.Context, code string) (*domain.TraceDocument, error) {
	// 1. Find the Sold Cup by Code
	// 2. Fetch Harvest details from Farm Service
	// 3. Fetch Batch details from Warehouse Service
	// 4. Assemble the complete story
}
```

### Interactive Elements

- [x] **Lineage Tree** — A visual tree showing: [Cup ➔ Sale ➔ Batch ➔ Harvest]. Clicking each node shows the unique ID that connects them.
- [x] **Code↔English translation** — `TraceDocument` structure and the assembly logic.
- [x] **QR Scanner Simulation** — A "Scan Me" button that triggers a transition to a "Verified" state.
- [x] **Quiz** — 3 questions, style: tracing. "If we lost the Batch ID, what part of the journey would disappear?"

### Reference Files to Read

- `references/interactive-elements.md` → "Lineage Tree", "Code ↔ English Translations"
- `references/content-philosophy.md`
- `references/gotchas.md`

### Connections

- **Previous module:** The Retail Dance (Generating the sale)
- **Next module:** When Nodes Whisper (Seeing the whole mesh)
- **Tone/style notes:** Scientific, precise, "Certified".
