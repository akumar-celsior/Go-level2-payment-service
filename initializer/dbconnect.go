package initializer

import (
	"fmt"
	"log"

	spannergorm "github.com/googleapis/go-gorm-spanner"
	_ "github.com/googleapis/go-sql-spanner" // Import the Spanner SQL driver
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DbInstance *gorm.DB

// init initializes the database connection when the package is imported.
func init() {
	// Load configuration from .env file
	LoadConfig()
	log.Println("Initializing database connection...")
	var err error
	DbInstance, err = ConnectSpannerDB()
	if err != nil {
		log.Fatalf("Failed to initialize database in init(): %v", err)
	}
	log.Println("Database connection initialized successfully.")
}

// ConnectSpannerDB initializes GORM with Google Cloud Spanner.
func ConnectSpannerDB() (*gorm.DB, error) {
	// Get the database connection details from the environment
	projectID := GetEnv("DB_PROJECT_ID")
	instanceID := GetEnv("DB_INSTANCE_ID")
	dbName := GetEnv("DB_NAME")

	// Build the Spanner connection string
	log.Printf("DB_PROJECT_ID: %s, DB_INSTANCE_ID: %s, DB_NAME: %s", projectID, instanceID, dbName)
	dsn := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, dbName)
	log.Printf("Connecting to Spanner with DSN: %s", dsn)
	// Connect to Spanner using the GORM driver for Spanner
	db, err := gorm.Open(spannergorm.New(spannergorm.Config{
		DriverName: "spanner", // Spanner driver name
		DSN:        dsn,       // The connection string (Data Source Name)
	}), &gorm.Config{
		PrepareStmt:                      true,
		IgnoreRelationshipsWhenMigrating: true,
		Logger:                           logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Printf("[error] failed to initialize database, DSN: %s, error: %v", dsn, err)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	DbInstance = db

	// Log success
	log.Printf("Successfully connected to Spanner database: %s", dsn)
	return db, nil
}

// GetDB returns the initialized database instance
func GetDB() *gorm.DB {
	if DbInstance == nil {
		log.Fatal("Database not initialized. Call ConnectSpannerDB first.")
		ConnectSpannerDB()
	}
	return DbInstance
}
