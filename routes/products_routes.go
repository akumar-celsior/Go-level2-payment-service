package routes

import (
	"fmt"
	"goTechReady/controller"

	"github.com/kataras/iris/v12"
)

func ProductRoutes(app *iris.Application) {
	fmt.Println("Registering products routes")
	productParty := app.Party("/products")
	productParty.Get("/", controller.GetProductsHandler)
	productParty.Get("/{id}", controller.GetProductByIDHandler)
	productParty.Post("/", controller.CreateProductHandler)
	productParty.Delete("/{id}", controller.DeleteProductHandler)
}
