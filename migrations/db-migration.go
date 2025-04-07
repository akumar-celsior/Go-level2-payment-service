package migrations

import (
	"goTechReady/initializer"
	"goTechReady/model"
	"log"
)

func Migrate() {
	var err error
	db := initializer.GetDB()

	if !db.Migrator().HasTable(&model.Product{}) {
		err = db.AutoMigrate(&model.Product{})
		if err != nil {
			log.Fatalf("Failed to migrate product table: %v", err)
		}
	}
}
