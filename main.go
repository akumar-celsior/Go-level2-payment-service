package main

import (
	"fmt"
	"log"
	"os"

	"goTechReady/migrations"
	"goTechReady/routes"

	"github.com/kataras/iris/v12"
)

func main() {
	// Create an Iris application instance
	app := iris.New()
	routes.ProductRoutes(app)
	fmt.Println("Routes initialized successfully!")
	log.Println("Routes registered successfully.")
	migrations.Migrate()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on port %s...\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
