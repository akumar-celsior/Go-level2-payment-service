package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"goPocDemo/initializer"
	"goPocDemo/routes"
	"goPocDemo/services"

	"runtime/pprof"

	"github.com/kataras/iris/v12"
)

var userService *services.UserService

func init() {
	// Load configuration from .env file
	initializer.LoadConfig()

	// Initialize Spanner client using GORM
	_, err := initializer.ConnectSpannerDB()
	if err != nil {
		log.Fatalf("Failed to initialize GORM Spanner client: %v", err)
	}
}

func main() {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close() // Ensure logs are written when the program exits
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// defer logFile.Close()
	log.Println("Logging started...")
	logFile.Sync()

	// create cpu profile
	cpuProfile, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer cpuProfile.Close()
	err = pprof.StartCPUProfile(cpuProfile)
	if err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	// create memory profile
	memProfile, err := os.Create("mem.prof")
	if err != nil {
		panic(err)
	}
	defer memProfile.Close()

	// Create an Iris application instance
	app := iris.New()
	log.Println("logging started")
	// Register routes
	log.Println("Registering auth routes...")
	routes.RegisterAuthRoutes(app, userService)
	fmt.Println("Registering transaction routes...")
	wg.Add(1)
	go func() {
		mutex.Lock()
		defer mutex.Unlock()
		routes.ProcessTransactionRoutes(app)
		wg.Done()
	}()
	wg.Wait()
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

	// executing memory profile
	err = pprof.WriteHeapProfile(memProfile)
	if err != nil {
		panic(err)
	}
}
