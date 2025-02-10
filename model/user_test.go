package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserTableName(t *testing.T) {
	user := User{}
	expectedTableName := "AuthorizedUsers"
	assert.Equal(t, expectedTableName, user.TableName(), "Table name should be 'AuthorizedUsers'")
}

func TestUserFields(t *testing.T) {
	now := time.Now()
	user := User{
		UserID:       "123",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FirstName:    "John",
		LastName:     "Doe",
		IsVerified:   true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	assert.Equal(t, "123", user.UserID, "UserID should be '123'")
	assert.Equal(t, "test@example.com", user.Email, "Email should be 'test@example.com'")
	assert.Equal(t, "hashedpassword", user.PasswordHash, "PasswordHash should be 'hashedpassword'")
	assert.Equal(t, "John", user.FirstName, "FirstName should be 'John'")
	assert.Equal(t, "Doe", user.LastName, "LastName should be 'Doe'")
	assert.True(t, user.IsVerified, "IsVerified should be true")
	assert.Equal(t, now, user.CreatedAt, "CreatedAt should be the current time")
	assert.Equal(t, now, user.UpdatedAt, "UpdatedAt should be the current time")
}
