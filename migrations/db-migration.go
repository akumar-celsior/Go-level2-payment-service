package migrations

import (
	"fmt"
	"goPocDemo/model"
	"log"

	"gorm.io/gorm"
)

// MigrateDown is responsible for rolling back migrations (dropping tables).
func MigrateDown(db *gorm.DB) error {
	// Drop tables using GORM's DropTable function
	if err := db.Migrator().DropTable(&model.User{}); err != nil {
		return fmt.Errorf("failed to rollback migration: %v", err)
	}
	log.Println("Tables rolled back successfully")
	return nil
}

// MigrateUp is responsible for applying migrations (creating tables).
func MigrateUp(db *gorm.DB) error {
	// Create tables using GORM's AutoMigrate function
	if err := db.AutoMigrate(&model.AuditLog{}); err != nil {
		return fmt.Errorf("failed to migrate tables: %v", err)
	}
	log.Println("Tables migrated successfully")
	return nil
}
