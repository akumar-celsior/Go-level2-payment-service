package routes

import (
	"goPocDemo/controller"
	"goPocDemo/middleware"

	"github.com/kataras/iris/v12"
)

func ProcessTransactionRoutes(app *iris.Application) {
	auth := app.Party("/process-transaction", middleware.AuthMiddleware)
	{
		auth.Post("/", func(ctx iris.Context) {
			controller.StartTransactionHandler(ctx)
		})
	}
}
