package types

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"unique"`
	Password string
}

func NewUser(email, password string) *User {
	return &User{
		Email:    email,
		Password: password,
		Name:     "",
	}
}

const (
	UserContextKey = "user"
)

type UserContext struct {
	IsLoggedIn bool
}

func NewUserContext() *UserContext {
	return &UserContext{
		IsLoggedIn: false,
	}
}
