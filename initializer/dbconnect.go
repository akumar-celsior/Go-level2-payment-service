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

// ConnectSpannerDB initializes GORM with Google Cloud Spanner.
func ConnectSpannerDB() (*gorm.DB, error) {
	// Get the database connection details from the environment
	projectID := GetEnv("DB_PROJECT_ID")
	instanceID := GetEnv("DB_INSTANCE_ID")
	dbName := GetEnv("DB_NAME")

	// Build the Spanner connection string
	dsn := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, dbName)

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
		log.Fatalf("Failed to open database: %v", err)
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
