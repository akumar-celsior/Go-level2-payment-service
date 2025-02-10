package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
)

func TestProcessTransactionRoutes(t *testing.T) {
	// t.Skip()
	app := iris.New()
	ProcessTransactionRoutes(app)
	// Mock middleware
	authMiddleware := func(ctx iris.Context) {
		ctx.Next()
	}
	app.Use(authMiddleware)

	// Mock controller handlers
	StartTransactionHandler := func(ctx iris.Context) {
		ctx.StatusCode(http.StatusCreated)
	}

	app.Post("/process-transaction", StartTransactionHandler)

	// Build the router
	app.Build()

	// Test POST /transactions
	req := httptest.NewRequest("POST", "/process-transaction", nil)
	resp := httptest.NewRecorder()
	app.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)
}
