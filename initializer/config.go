package initializer

import (
	"goTechReady/model"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DbInstance *gorm.DB

// LoadConfig loads environment variables from a .env file
func LoadConfig() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Ensure the required environment variables are set
	// if os.Getenv("DB_PROJECT_ID") == "" {
	// 	log.Fatal("Missing DB_PROJECT_ID environment variable")
	// }
	// if os.Getenv("DB_INSTANCE_ID") == "" {
	// 	log.Fatal("Missing DB_INSTANCE_ID environment variable")
	// }
	if os.Getenv("DB_NAME") == "" {
		log.Fatal("Missing DB_NAME environment variable")
	}

	// Ensure GOOGLE_APPLICATION_CREDENTIALS is set if running locally
	// if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
	// 	log.Fatal("Missing GOOGLE_APPLICATION_CREDENTIALS environment variable")
	// }
}

// GetEnv is a helper function to fetch environment variables
func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}

// ConnectCloudSQL establishes a GORM connection to a Cloud SQL instance using environment variables
func ConnectCloudSQL() *gorm.DB {
	dbName := GetEnv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	cloudSQLConn := os.Getenv("CLOUDSQL_CONNECTION_NAME")
	cloudSQLHost := os.Getenv("CLOUDSQL_HOST") // e.g., 127.0.0.1 or public IP
	cloudSQLPort := os.Getenv("CLOUDSQL_PORT") // e.g., 3306
	if cloudSQLHost == "" {
		cloudSQLHost = "127.0.0.1" // default to localhost
	}
	if cloudSQLPort == "" {
		cloudSQLPort = "3306" // default MySQL port
	}

	var dsn string
	if os.PathSeparator == '\\' { // Windows uses backslash
		// Use TCP connection on Windows
		dsn = user + ":" + password + "@tcp(" + cloudSQLHost + ":" + cloudSQLPort + ")/" + dbName + "?parseTime=true"
	} else {
		// Use Unix socket on Linux/Mac
		if cloudSQLConn == "" {
			log.Fatal("Missing CLOUDSQL_CONNECTION_NAME environment variable")
		}
		dsn = user + ":" + password + "@unix(/cloudsql/" + cloudSQLConn + ")/" + dbName + "?parseTime=true"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to Cloud SQL with GORM: %v", err)
	}

	// Conditionally automigrate only if the table does not exist
	if !db.Migrator().HasTable(&model.Order{}) {
		db.AutoMigrate(&model.Order{})
	}
	DbInstance = db
	return db
}

// GetDB returns the initialized database instance
func GetDB() *gorm.DB {
	if DbInstance == nil {
		log.Fatal("Database not initialized. Call ConnectSpannerDB first.")
		ConnectCloudSQL()
	}
	return DbInstance
}
