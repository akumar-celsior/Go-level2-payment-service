package initializer

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

func TestConnectSpannerDB(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("DB_PROJECT_ID", "iconic-star-447805-v6")
	os.Setenv("DB_INSTANCE_ID", "go-instance-poc")
	os.Setenv("DB_NAME", "go_spanner_db") // Set DB_NAME here
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, "projects/iconic-star-447805-v6/instances/go-instance-poc/databases/go_spanner_db", option.WithoutAuthentication())
	assert.NoError(t, err, "Expected no error from spanner.NewClient, got %v", err)
	db := client

	// Call the function to test
	gormDB, err := ConnectSpannerDB()
	assert.NotNil(t, gormDB, "Expected gormDB instance, got nil")

	// Assert no error occurred
	assert.NoError(t, err, "Expected no error, got %v", err)

	// Assert the returned DB instance is not nil
	assert.NotNil(t, db, "Expected db instance, got nil")

	// Clean up environment variables
	os.Unsetenv("DB_PROJECT_ID")
	os.Unsetenv("DB_INSTANCE_ID")
	os.Unsetenv("DB_NAME")
}
func TestGetDB(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("DB_PROJECT_ID", "iconic-star-447805-v6")
	os.Setenv("DB_INSTANCE_ID", "go-instance-poc")
	os.Setenv("DB_NAME", "go_spanner_db")

	// Ensure the database is connected
	_, err := ConnectSpannerDB()
	assert.NoError(t, err, "Expected no error from ConnectSpannerDB, got %v", err)

	// Call the function to test
	db := GetDB()

	// Assert the returned DB instance is not nil
	assert.NotNil(t, db, "Expected db instance, got nil")

	// Clean up environment variables
	os.Unsetenv("DB_PROJECT_ID")
	os.Unsetenv("DB_INSTANCE_ID")
	os.Unsetenv("DB_NAME")
}
