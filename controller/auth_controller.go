package controller

import (
	"fmt"
	"goPocDemo/services"
	"goPocDemo/utils"

	"github.com/kataras/iris/v12"
)

// LoginHandler handles user login requests.
func LoginHandler(svc *services.UserService, ctx iris.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]string{"error": "Failed to parse request"})
		return
	}

	// Call the service to get the user by email (only passing the email)
	user, err := svc.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Verify the password
	if err := utils.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "Invalid credentials"})
		return
	}

	// Respond with success
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(map[string]string{"message": "Login successful", "userID": user.UserID})
}

// SignupHandler handles user signup requests.
func SignupHandler(svc *services.UserService, ctx iris.Context) {
	fmt.Print("SignupHandler ctx", ctx)
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]string{"error": "Failed to parse request"})
		return
	}
	// Call the service to create the user
	user, err := svc.CreateUser(ctx, req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": "Failed to create user: " + err.Error()})
		return
	}
	fmt.Print("SignupHandler user", user)

	// Respond with success
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(map[string]string{"message": "User created successfully", "userID": user.UserID})
}
