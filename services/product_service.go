package services

import (
	"goTechReady/initializer"
	"goTechReady/model"

	"github.com/kataras/iris/v12"
)

var db = initializer.GetDB()

func CreateProduct(ctx iris.Context) {
	var product model.Product

	if err := ctx.ReadJSON(&product); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid input", "details": err.Error()})
		return
	}

	result := db.Create(&product)
	if result.Error != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to create product", "details": result.Error.Error()})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(product)
}
func DeleteProduct(ctx iris.Context) {
	id := ctx.Params().GetString("id")
	if id == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Product ID is required"})
		return
	}

	result := db.Delete(&model.Product{}, "id = ?", id)
	if result.Error != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to delete product", "details": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Product not found"})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Product deleted successfully"})
}
func GetAllProducts(ctx iris.Context) {
	var products []model.Product
	if err := db.Find(&products).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to retrieve products"})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(products)
}
func GetProductByID(ctx iris.Context) {
	id := ctx.Params().GetString("id")
	if id == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid ID"})
		return
	}

	var product model.Product
	result := db.Where("id = ?", id).First(&product)
	if result.Error != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Product not found"})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(product)
}
func SeedData() {
	var count int64
	db.Model(&model.Product{}).Count(&count)
	if count == 0 {
		db.Create(&model.Product{Name: "Laptop", Price: 1000})
		db.Create(&model.Product{Name: "Phone", Price: 500})
		db.Create(&model.Product{Name: "Tablet", Price: 750})
		db.Create(&model.Product{Name: "Monitor", Price: 300})
		db.Create(&model.Product{Name: "Keyboard", Price: 50})
		db.Create(&model.Product{Name: "Mouse", Price: 25})
		db.Create(&model.Product{Name: "Headphones", Price: 150})
	}
}
