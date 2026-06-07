# Module 4: The Retail Dance

### Teaching Arc
- **Metaphor:** A "Choreographed Performance". When you buy a coffee, five different things must happen in perfect order: Pay -> Reserve Beans -> Update Stock -> Confirm. If any step trips, the whole dance stops and resets.
- **Opening hook:** You click "Buy". In that split second, the system starts a complex "Saga" to make sure your coffee is guaranteed.
- **Key insight:** The "Saga Pattern" ensures that complex actions involving many services either all succeed together or fail safely without leaving a mess.
- **"Why should I care?":** This is the secret to building high-quality, professional software. You'll know how to tell AI to "Handle this as a Saga" so you never lose a customer's money.

### Code Snippets (pre-extracted)

File: src/apps/retail-service/internal/domain/order.go (lines 40-55)
```go
const (
	OrderStatusPending          = "PENDING"
	OrderStatusPaymentCompleted = "PAYMENT_COMPLETED"
	OrderStatusStockReserved    = "STOCK_RESERVED"
	OrderStatusInTransit        = "IN_TRANSIT"
	OrderStatusCompleted        = "COMPLETED"
	OrderStatusRejected         = "REJECTED"
)
```

File: src/apps/retail-service/internal/usecase/service.go (lines 450-465 - Saga trigger)
```go
func (s *Service) handlePaymentCompleted(ctx context.Context, msg events.PaymentCompleted) error {
	// 1. Mark payment done
	// 2. Request stock reservation from Warehouse
	// 3. Update Order status to PAYMENT_COMPLETED
    // This is the first step of the Saga
}
```

### Interactive Elements

- [x] **Status Pipeline** — A horizontal timeline of the Saga steps: [Payment ➔ Reserved ➔ Shipping ➔ Delivered].
- [x] **Code↔English translation** — Order statuses and the Saga trigger function.
- [x] **"The Undo Button" visualization** — An animation showing what happens when a step fails: the "baton" goes backwards and undoes the previous steps (Compensating Transactions).
- [x] **Quiz** — 4 questions, style: architecture. "Why don't we just update everything at once? Why use steps?" (Reliability/Consistency).

### Reference Files to Read

- `references/interactive-elements.md` → "Status Pipeline", "Code ↔ English Translations"
- `references/content-philosophy.md`
- `references/gotchas.md`

### Connections

- **Previous module:** The Logistics Pulse (Moving data)
- **Next module:** Transparent Proof (The result of the Saga)
- **Tone/style notes:** Sophisticated but clear. Focus on "Trust" and "Guarantees".
