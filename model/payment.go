package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID     string
	Amount float64
}

type Payment struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	OrderID string
	Amount  float64
	Status  string
}

// BeforeCreate will set a UUID rather than numeric ID.
func (o *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}
