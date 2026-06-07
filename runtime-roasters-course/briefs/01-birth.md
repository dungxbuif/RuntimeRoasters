# Module 1: The Bean's Digital Birth

### Teaching Arc
- **Metaphor:** A "Digital Birth Certificate" for every bean. Just like a newborn gets an ID and a record of where they were born, every coffee harvest gets a permanent digital identity.
- **Opening hook:** Imagine you're standing in a Dalat coffee farm, holding a handful of fresh cherries. How does that physical moment become a line of data that follows the bean halfway across the world?
- **Key insight:** The journey begins with a "Declaration". Data isn't just "there" — it has to be born through a deliberate act of recording.
- **"Why should I care?":** Understanding how data is born helps you steer AI to build better input systems and catch errors at the source. If the birth record is wrong, the whole life story is a lie.

### Code Snippets (pre-extracted)

File: src/apps/farm-service/internal/domain/harvest.go (lines 10-22)
```go
type Harvest struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FarmID      uint      `json:"farm_id"`
	CoffeeType  string    `json:"coffee_type"`
	Quantity    float64   `json:"quantity"`
	HarvestDate time.Time `json:"harvest_date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
```

File: src/apps/farm-service/internal/usecase/service.go (lines 120-135)
```go
func (s *Service) CreateHarvest(ctx context.Context, input domain.CreateHarvestInput) (*domain.Harvest, error) {
	harvest := &domain.Harvest{
		FarmID:      input.FarmID,
		CoffeeType:  input.CoffeeType,
		Quantity:    input.Quantity,
		HarvestDate: input.HarvestDate,
		Status:      "NEW",
	}

	if err := s.db.Create(harvest).Error; err != nil {
		return nil, err
	}
    // ...
```

### Interactive Elements

- [x] **Code↔English translation** — `Harvest` struct and `CreateHarvest` function.
- [x] **Quiz** — 3 questions, style: scenario. Focus on: What happens if a harvest is recorded twice? (Idempotency), Why do we need a status?
- [x] **Data flow animation** — actors: [Farm Manager, Farm Service, Database]. Steps: [1. Record click, 2. API Request, 3. DB Write, 4. Success UI].
- [x] **Hero Visual** — A "Digital Birth Certificate" card showing the fields of a Harvest object.

### Reference Files to Read

- `references/interactive-elements.md` → "Data Flow Animation", "Code ↔ English Translations"
- `references/content-philosophy.md`
- `references/gotchas.md`

### Connections

- **Previous module:** None (Intro)
- **Next module:** The Cast of Services (Zooming out to the whole system)
- **Tone/style notes:** Vermillion accent. Focus on Dalat/Vietnam coffee vibe.
