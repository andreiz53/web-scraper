package db

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// type Storage interface {
// 	CreateUser(*types.User) error
// 	GetUserByEmail(int) (*types.User, error)
// 	UpdateUser(*types.User) error
// 	DeleteUser(int) error
// }

type SqliteStorage struct {
	Db *gorm.DB
}

func NewSqliteStorage() (*SqliteStorage, error) {
	db, err := openDatabase()
	if err != nil {
		return &SqliteStorage{}, err
	}

	err = createMigrations(db)
	if err != nil {
		return &SqliteStorage{}, err
	}

	ok := seedUsers(db)
	if !ok {
		fmt.Println("Users already seeded")
	}

	return &SqliteStorage{Db: db}, nil
}

func openDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DB_NAME") + fmt.Sprintf("?_pragma_key=%s&_pragma_cipher_page_size=4096", os.Getenv("DB_PASSWORD"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
