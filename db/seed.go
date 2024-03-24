package db

import (
	"github.com/andreiz53/web-scraper/types"
	"gorm.io/gorm"
)

func seedUsers(db *gorm.DB) bool {
	var user types.User
	db.First(&user, 1)
	if user.ID != 0 {
		return false
	}
	db.Create(&types.User{
		Email:    "andrei@gmail.com",
		Password: "asd",
	})
	db.Create(&types.User{
		Email:    "andrei2@gmail.com",
		Password: "asd2",
	})
	db.Create(&types.User{
		Email:    "andrei3@gmail.com",
		Password: "asd3",
	})
	return true

}
