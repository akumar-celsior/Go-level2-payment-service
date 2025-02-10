package services

import (
	"context"
	"errors"
	"time"

	"goPocDemo/initializer"
	"goPocDemo/model"
	"goPocDemo/utils"

	"gorm.io/gorm"
)

// // Define UserServiceInterface instead of a struct
// type UserServiceInterface interface {
// 	CreateUser(ctx context.Context, email, password, firstName, lastName string) (*model.User, error)
// 	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
// }

// var _ UserServiceInterface = &UserService{}

// UserService provides methods for user-related operations.
type UserService struct {
	// DB *gorm.DB
}

// NewUserService creates a new instance of UserService.
// func NewUserService(db *gorm.DB) *UserService {
// 	return &UserService{DB: db}
// }

// CreateUser creates a new user in the database.
func (svc *UserService) CreateUser(ctx context.Context, email, password, firstName, lastName string) (*model.User, error) {
	dbInstance := initializer.GetDB()
	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	// Create new user model
	user := &model.User{
		UserID:       time.Now().Format("20060102150405"), // Unique ID based on timestamp
		Email:        email,
		PasswordHash: hashedPassword,
		FirstName:    firstName,
		LastName:     lastName,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Insert user into the database using GORM
	if err := dbInstance.Create(user).Error; err != nil {
		return nil, err
	}

	// Return user and nil error
	return user, nil
}

// GetUserByEmail retrieves a user by email from the database (for login).
func (svc *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	dbInstance := initializer.GetDB()
	var user model.User

	// Query the database using GORM to fetch the user by email
	if err := dbInstance.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Return user if found
	return &user, nil
}
