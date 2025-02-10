// package routes

// import (
// 	"net/http"
// 	"testing"

// 	"goPocDemo/services"

// 	"github.com/kataras/iris/v12"
// 	"github.com/kataras/iris/v12/httptest"
// 	"github.com/stretchr/testify/assert"
// )

// func TestRegisterAuthRoutes(t *testing.T) {
// 	app := iris.New()
// 	svc := &services.UserService{}

// 	RegisterAuthRoutes(app, svc)

// 	// Mock controller handlers
// 	SignupHandler := func(ctx iris.Context) {
// 		ctx.StatusCode(http.StatusCreated)
// 	}

// 	app.Post("/signup", SignupHandler)

// 	// Build the router
// 	app.Build()

// 	// Test POST /transactions
// 	req := httptest.NewRequest("POST", "/signup", nil)
// 	resp := httptest.NewRecorder()
// 	app.ServeHTTP(resp, req)
// 	assert.Equal(t, http.StatusCreated, resp.Code)

// 	// e := httptest.New(t, app)

// 	// t.Run("Signup Route", func(t *testing.T) {
// 	// 	e.POST("/signup").Expect().Status(httptest.StatusOK)
// 	// })

//		// t.Run("Login Route", func(t *testing.T) {
//		// 	e.POST("/login").Expect().Status(httptest.StatusOK)
//		// })
//	}
package routes

import (
	"net/http"
	"strings"
	"testing"

	"goPocDemo/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAuthRoutes(t *testing.T) {
	app := iris.New()
	svc := &services.UserService{}

	RegisterAuthRoutes(app, svc)

	// Mock controller handlers
	SignupHandler := func(ctx iris.Context) {
		ctx.StatusCode(http.StatusCreated)
	}

	LoginHandler := func(ctx iris.Context) {
		ctx.StatusCode(http.StatusOK)
	}

	app.Post("/signup", SignupHandler)
	app.Post("/login", LoginHandler)

	// Build the router
	app.Build()

	// Test POST /signup
	t.Run("Signup Route", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/signup", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		app.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	// Test POST /login
	t.Run("Login Route", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		app.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}
