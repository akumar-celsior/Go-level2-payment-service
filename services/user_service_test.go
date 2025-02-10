package services

import (
	"context"
	"testing"
	"time"

	"goPocDemo/initializer"
	"goPocDemo/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/clause"
)

func TestGetUserByEmail(t *testing.T) {
	// Initialize the database and mock data
	db := initializer.GetDB()
	initializer.DbInstance = db
	//db.AutoMigrate(&model.User{})

	// Create a user for testing
	testUser := &model.User{
		UserID:       "123456",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FirstName:    "Test",
		LastName:     "User",
		IsVerified:   true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(testUser)

	// Create an instance of UserService
	userService := &UserService{}

	t.Run("User found", func(t *testing.T) {
		user, err := userService.GetUserByEmail(context.Background(), "test@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("User not found", func(t *testing.T) {
		user, err := userService.GetUserByEmail(context.Background(), "nonexistent@example.com")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "invalid credentials", err.Error())
	})
}
func TestCreateUser(t *testing.T) {
	// Initialize the database and mock data
	db := initializer.GetDB()
	initializer.DbInstance = db

	// Create an instance of UserService
	userService := &UserService{}

	t.Run("Create user successfully", func(t *testing.T) {
		email := "newuser@example.com"
		password := "password123"
		firstName := "New"
		lastName := "User"

		user, err := userService.CreateUser(context.Background(), email, password, firstName, lastName)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, firstName, user.FirstName)
		assert.Equal(t, lastName, user.LastName)
		assert.False(t, user.IsVerified)
	})

	// t.Run("Create user with invalid password", func(t *testing.T) {
	// 	email := "invalidpassword@example.com"
	// 	password := "" // Invalid password
	// 	firstName := "Invalid"
	// 	lastName := "Password"

	// 	user, err := userService.CreateUser(context.Background(), email, password, firstName, lastName)
	// 	assert.Error(t, err)
	// 	assert.Nil(t, user)
	// })
}
