package domain

import "time"

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusSucceeded PaymentStatus = "SUCCEEDED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
)

type Payment struct {
	ID          string        `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID     string        `gorm:"type:uuid;uniqueIndex;not null" json:"order_id"`
	StoreID     string        `gorm:"type:uuid;not null;index" json:"store_id"`
	Provider    string        `gorm:"size:40;not null;index" json:"provider"`
	ProviderRef string        `gorm:"size:160;uniqueIndex;not null" json:"provider_ref"`
	RefundRef   string        `gorm:"size:160;index" json:"refund_ref"`
	Items       string        `gorm:"type:jsonb;not null" json:"items"`
	Amount      float64       `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency    string        `gorm:"size:8;not null" json:"currency"`
	Status      PaymentStatus `gorm:"size:24;not null;index" json:"status"`
	CheckoutURL string        `gorm:"size:512" json:"checkout_url"`
	Simulated   bool          `gorm:"not null;default:true" json:"simulated"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type InboxEvent struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	MessageID   string    `gorm:"uniqueIndex;not null"`
	EventType   string    `gorm:"size:160;not null;index"`
	ProcessedAt time.Time `gorm:"not null"`
}

type WebhookEvent struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	Provider    string    `gorm:"size:40;not null;uniqueIndex:idx_provider_event"`
	EventID     string    `gorm:"size:160;not null;uniqueIndex:idx_provider_event"`
	ProviderRef string    `gorm:"size:160;not null;index"`
	Status      string    `gorm:"size:40;not null"`
	Reason      string    `gorm:"size:512"`
	ProcessedAt time.Time `gorm:"not null"`
}
