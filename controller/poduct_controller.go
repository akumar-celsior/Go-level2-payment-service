package controller

import (
	"fmt"
	"goTechReady/services"

	"github.com/kataras/iris/v12"
)

func CreateProductHandler(ctx iris.Context) {
	fmt.Println("CreateProductHandler called")
	services.CreateProduct(ctx)
}
func DeleteProductHandler(ctx iris.Context) {
	fmt.Println("DeleteProductHandler called")
	services.DeleteProduct(ctx)
}
func GetProductsHandler(ctx iris.Context) {
	fmt.Println("GetProductsHandler called")
	services.GetAllProducts(ctx)
}
func GetProductByIDHandler(ctx iris.Context) {
	fmt.Println("GetProductByIDHandler called")
	services.GetProductByID(ctx)
}
