package routes

import (
	"fmt"
	"goPocDemo/controller"
	"goPocDemo/services"

	"github.com/kataras/iris/v12"
)

func RegisterAuthRoutes(app *iris.Application, svc *services.UserService) {
	fmt.Println("Registering auth routes")
	app.Post("/signup", func(ctx iris.Context) {
		fmt.Println("Registering signup route came here!")
		controller.SignupHandler(svc, ctx)
	})

	app.Post("/login", func(ctx iris.Context) {
		controller.LoginHandler(svc, ctx)
	})
}
