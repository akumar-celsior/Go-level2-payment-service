package initializer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test LoadConfig function
func TestLoadConfig(t *testing.T) {
	t.Skip()
	// Backup existing environment variables
	// originalEnv := map[string]string{
	// 	"DB_PROJECT_ID":  os.Getenv("DB_PROJECT_ID"),
	// 	"DB_INSTANCE_ID": os.Getenv("DB_INSTANCE_ID"),
	// 	"DB_NAME":        os.Getenv("DB_NAME"),
	// 	//"GOOGLE_APPLICATION_CREDENTIALS": os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
	// }

	// // Restore environment variables after test
	// defer func() {
	// 	for key, value := range originalEnv {
	// 		_ = os.Setenv(key, value)
	// 	}
	// }()

	// Set required environment variables
	_ = os.Setenv("DB_PROJECT_ID", "iconic-star-447805-v6")
	_ = os.Setenv("DB_INSTANCE_ID", "go-instance-poc")
	_ = os.Setenv("DB_NAME", "go_spanner_db")
	//_ = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "test_credentials.json")

	// Call LoadConfig
	LoadConfig()

	// Validate that environment variables are correctly set
	assert.Equal(t, "iconic-star-447805-v6", os.Getenv("DB_PROJECT_ID"))
	assert.Equal(t, "go-instance-poc", os.Getenv("DB_INSTANCE_ID"))
	assert.Equal(t, "go_spanner_db", os.Getenv("DB_NAME"))
	//assert.Equal(t, "test_credentials.json", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
}

// Test LoadConfig when .env file is missing
// func TestLoadConfig_MissingEnvFile(t *testing.T) {
// 	// Override log.Fatal to avoid test exit
// 	//logFatal := log.Fatal
// 	// defer func() { log.Fatal = logFatal }() // Restore original log.Fatal

// 	// var fatalMessage string
// 	// log.Fatal = func(v ...interface{}) {
// 	// 	fatalMessage = v[0].(string)
// 	// }

// 	// Temporarily rename the .env file to simulate missing file scenario
// 	_ = os.Rename(".env", ".env.bak")
// 	defer os.Rename(".env.bak", ".env") // Restore after test

// 	// LoadConfig should fail due to missing .env file
// 	LoadConfig()

// 	// Validate log output
// 	assert.Contains(t, "fatal message", "Error loading .env file")
// }
