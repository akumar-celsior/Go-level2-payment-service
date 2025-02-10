package middleware

import (
	"strings"

	"github.com/kataras/iris/v12"
)

const staticToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6InByYWthc2giLCJpYXQiOjE1MTYyMzkwMjJ9.0fwxvOHydN6uWYfbLazT0QcVeIMMDx6acveQhaY4eho"

func AuthMiddleware(ctx iris.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{
			"error": "Missing or invalid Authorization header",
		})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	if token != staticToken {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{
			"error": "Invalid token",
		})
		return
	}

	ctx.Next()
}
