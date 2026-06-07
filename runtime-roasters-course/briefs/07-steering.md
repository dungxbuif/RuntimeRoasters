# Module 7: Steering the Roaster

### Teaching Arc
- **Metaphor:** The "Captain's Wheel". You don't need to know how the engine works to sail the ship, but you *do* need to know which way to turn the wheel. You are now the Captain of this codebase.
- **Opening hook:** You've seen the bean's journey from dirt to data. Now, how do you use this knowledge to build the next version?
- **Key insight:** Technical fluency is a superpower. You can now describe *what* you want in the language the system *understands*.
- **"Why should I care?":** This is your ROI. You'll work 10x faster with AI because you'll stop saying "the button thing" and start saying "add a field to the Harvest domain model."

### Code Snippets (pre-extracted)

File: src/apps/farm-service/internal/domain/harvest.go (lines 1-5 - generic example)
```go
// Add a new quality score field!
type Harvest struct {
    // ...
    QualityScore float64 `json:"quality_score"` // 0.0 to 10.0
}
```

### Interactive Elements

- [x] **"AI Command" Translator** — A tool where you type a "Vibe" (e.g., "Add a rating to coffee") and it translates it into a "Technical Instruction" (e.g., "Add a QualityScore field to the Harvest domain and update the UI form").
- [x] **The Future Map** — A list of "What to build next" cards: [Multi-currency, Geo-fencing, AI Roasting Predictions].
- [x] **Interactive Checklist** — "You are ready to... [ ] Design new services, [ ] Debug Saga failures, [ ] Steer AI agents."
- [x] **Quiz** — 3 questions, style: application. "You want to add a 10% discount for 'Loyal Farmers'. Where do you start?"

### Reference Files to Read

- `references/content-philosophy.md`
- `references/gotchas.md`

### Connections

- **Previous module:** When Nodes Whisper (The big picture)
- **Next module:** None (Completion)
- **Tone/style notes:** Empowering, inspiring, "Graduation" vibe.
