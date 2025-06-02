package main

import (
	"fmt"
	"goTechReady/initializer"
	"goTechReady/pubsub"
	"goTechReady/routes"
	"log"
	"os"

	"github.com/kataras/iris/v12"
)

func main() {
	// Create an Iris application instance
	app := iris.New()
	initializer.LoadConfig()
	initializer.ConnectCloudSQL() // Only call this once to initialize the global DB
	fmt.Println("db:", initializer.GetDB())
	routes.RegisterPaymentRoutes(app) // Register payment routes
	go pubsub.ListenToOrders()        // Start listening to Pub/Sub messages in a goroutine

	fmt.Println("Routes initialized successfully!")
	log.Println("Routes registered successfully.")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on port %s...\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
