package domain

import (
	"time"
)

type Warehouse struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string    `gorm:"size:120;not null" json:"name"`
	Code         string    `gorm:"size:20;not null;uniqueIndex" json:"code"`
	Location     string    `gorm:"size:255;not null" json:"location"`
	ManagerID    string    `gorm:"type:varchar(64);index" json:"manager_id"`
	ManagerEmail string    `gorm:"size:160;index" json:"manager_email"`
	Capacity     string    `gorm:"size:50" json:"capacity"`
	Status       string    `gorm:"size:20;not null;default:'ACTIVE'" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
