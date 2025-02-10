package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goPocDemo/controller"
	"goPocDemo/model"
	"goPocDemo/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

// "gorm.io/driver/sqlite"

func setupTestDB() *gorm.DB {
	// db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	db.AutoMigrate(&model.User{})
	return db
}

func TestSignupHandler(t *testing.T) {
	app := iris.New()
	db := setupTestDB()
	if db == nil {
		t.Fatal("Database initialization failed")
	}
	svc := new(services.UserService)
	if svc == nil {
		t.Fatal("UserService initialization failed")
	}

	//routes.RegisterAuthRoutes(app, svc)

	app.Post("/signup", func(ctx iris.Context) {
		controller.SignupHandler(svc, ctx)
	})

	t.Run("Successful Signup", func(t *testing.T) {
		t.Skip() // userData := map[string]string{
		// 	"email":      "test@example.com",
		// 	"password":   "password123",
		// 	"first_name": "John",
		// 	"last_name":  "Doe",
		// }
		userData := map[string]string{"UserID": "12", "email": "newuser@example.com", "passwordHash": "password123",
			"first_name": "Jane", "last_name": "Doe"}
		body, _ := json.Marshal(userData)
		fmt.Println(body)
		// request := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
		request := httptest.NewRequest(http.MethodPost, "/signup", nil)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		fmt.Println("Debug: Sending request to /signup")
		fmt.Println("respinse", response)

		//app.ServeHTTP(response, request)

		// assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		t.Skip()
		request := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader([]byte(`invalid json`)))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		// app.ServeHTTP(response, request)
		// assert.Equal(t, http.StatusBadRequest, response.Code)

		fmt.Println("respinse", response)
	})
}
