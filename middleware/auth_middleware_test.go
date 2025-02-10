package middleware

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func TestAuthMiddleware(t *testing.T) {
	app := iris.New()
	app.Use(AuthMiddleware)
	app.Get("/test", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "success"})
	})

	e := httptest.New(t, app)

	t.Run("Missing Authorization Header", func(t *testing.T) {
		e.GET("/test").
			Expect().
			Status(iris.StatusUnauthorized).
			JSON().Object().ValueEqual("error", "Missing or invalid Authorization header")
	})

	t.Run("Invalid Authorization Header", func(t *testing.T) {
		e.GET("/test").
			WithHeader("Authorization", "InvalidToken").
			Expect().
			Status(iris.StatusUnauthorized).
			JSON().Object().ValueEqual("error", "Missing or invalid Authorization header")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		e.GET("/test").
			WithHeader("Authorization", "Bearer InvalidToken").
			Expect().
			Status(iris.StatusUnauthorized).
			JSON().Object().ValueEqual("error", "Invalid token")
	})

	t.Run("Valid Token", func(t *testing.T) {
		e.GET("/test").
			WithHeader("Authorization", "Bearer "+staticToken).
			Expect().
			Status(iris.StatusOK).
			JSON().Object().ValueEqual("message", "success")
	})
}
