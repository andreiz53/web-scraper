package db

import (
	"github.com/andreiz53/web-scraper/types"
	"gorm.io/gorm"
)

func createMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&types.User{})
	if err != nil {
		return err
	}
	return nil
}
