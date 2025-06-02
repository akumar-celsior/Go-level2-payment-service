package controller

import (
	"fmt"
	"goTechReady/initializer"
	"goTechReady/model"

	"github.com/kataras/iris/v12"
)

func GetPaymentHandler(ctx iris.Context) {
	fmt.Println("GetPaymentHandler called")
	var payments []model.Payment
	db := initializer.GetDB()
	if err := db.Find(&payments).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to retrieve payments", "details": err.Error()})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(payments)
}
