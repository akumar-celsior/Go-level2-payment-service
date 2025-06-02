package routes

import (
	"goTechReady/controller"

	"github.com/kataras/iris/v12"
)

func RegisterPaymentRoutes(app *iris.Application) {
	orderAPI := app.Party("/payments")

	orderAPI.Get("/", controller.GetPaymentHandler)
}
