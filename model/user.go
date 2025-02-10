package model

import "time"

type User struct {
	UserID       string    `gorm:"primaryKey;column:UserID"`
	Email        string    `gorm:"unique;column:Email"`
	PasswordHash string    `gorm:"column:PasswordHash"`
	FirstName    string    `gorm:"column:FirstName"`
	LastName     string    `gorm:"column:LastName"`
	IsVerified   bool      `gorm:"column:IsVerified"`
	CreatedAt    time.Time `gorm:"column:CreatedAt"`
	UpdatedAt    time.Time `gorm:"column:UpdatedAt"`
}

// TableName explicitly sets the table name to "Users" (case-sensitive)
func (User) TableName() string {
	return "AuthorizedUsers" // Use the exact case-sensitive name you want
}
