package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func (u *Product) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	return nil
}
