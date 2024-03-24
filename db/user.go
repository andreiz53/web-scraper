package db

import (
	"github.com/andreiz53/web-scraper/types"
)

func (s *SqliteStorage) CreateUser(u *types.User) error {
	result := s.Db.Create(u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SqliteStorage) GetUserByEmail(email string) (types.User, error) {
	var user types.User
	result := s.Db.Where(&types.User{Email: email}).First(&user)
	if result.Error != nil {
		return types.User{}, result.Error
	}
	return user, nil
}

func (s *SqliteStorage) UpdateUser(u *types.User) error {
	result := s.Db.Save(u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SqliteStorage) DeleteUser(id int) error {
	result := s.Db.Delete(&types.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
